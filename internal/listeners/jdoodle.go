package listeners

import (
	"fmt"
	"strings"
	"time"

	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/timedmap"
	"golang.org/x/time/rate"

	"github.com/zekroTJA/shinpuru/internal/middleware"
	"github.com/zekroTJA/shinpuru/internal/services/codeexec"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/embedbuilder"

	"github.com/bwmarrin/discordgo"
)

const (
	removeHandlerTimeout         = 3 * time.Minute
	removeHandlerCleanupInterval = 1 * time.Minute

	// Ratelimiter settings
	limitTMCleanupInterval = 30 * time.Second // 10 * time.Minute
	limitTMLifetime        = 24 * time.Hour
	limitBurst             = 2                // 5
	limitRate              = float64(1) / 900 // one token per 15 minutes
)

var (
	runReactionEmoji = "▶"

	// langs = []string{"java", "c", "cpp", "c99", "cpp14", "php", "perl", "python3", "ruby", "go", "scala", "bash", "sql", "pascal", "csharp",
	// 	"vbn", "haskell", "objc", "ell", "swift", "groovy", "fortran", "brainfuck", "lua", "tcl", "hack", "rust", "d", "ada", "r", "freebasic",
	// 	"verilog", "cobol", "dart", "yabasic", "clojure", "nodejs", "scheme", "forth", "prolog", "octave", "coffeescript", "icon", "fsharp", "nasm",
	// 	"gccasm", "intercal", "unlambda", "picolisp", "spidermonkey", "rhino", "bc", "clisp", "elixir", "factor", "falcon", "fantom", "pike", "smalltalk",
	// 	"mozart", "lolcode", "racket", "kotlin"}

	replaces = map[string]string{
		"js":         "nodejs",
		"javascript": "nodejs",
		"c++":        "cpp",
		"c#":         "csharp",
		"python":     "python3",
		"py":         "python3",
	}
)

type ListenerJdoodle struct {
	db       database.Database
	execFact codeexec.Factory
	pmw      *middleware.PermissionsMiddleware

	langs  []string
	limits *timedmap.TimedMap
	msgMap *timedmap.TimedMap
}

type jdoodleMessage struct {
	*discordgo.Message

	wrapper codeexec.Executor
	lang    string
	script  string

	embLang string
}

func NewListenerJdoodle(container di.Container) (l *ListenerJdoodle) {
	l = &ListenerJdoodle{}

	l.db = container.Get(static.DiDatabase).(database.Database)
	l.pmw = container.Get(static.DiPermissionMiddleware).(*middleware.PermissionsMiddleware)
	l.execFact = container.Get(static.DiCodeExecFactory).(codeexec.Factory)
	l.langs = l.execFact.Languages()
	l.limits = timedmap.New(limitTMCleanupInterval)
	l.msgMap = timedmap.New(removeHandlerCleanupInterval)

	return
}

func (l *ListenerJdoodle) HandlerMessageCreate(s *discordgo.Session, e *discordgo.MessageCreate) {
	l.handler(s, e.Message)
}

func (l *ListenerJdoodle) HandlerMessageUpdate(s *discordgo.Session, e *discordgo.MessageUpdate) {
	l.handler(s, e.Message)
}

func (l *ListenerJdoodle) handler(s *discordgo.Session, e *discordgo.Message) {
	if e.Author == nil || e.Author.Bot || e.GuildID == "" {
		return
	}

	lang, cont, ok := l.parseMessageContent(e.Content)
	if !ok {
		return
	}

	embLang := lang

	if lang == "" || cont == "" {
		return
	}

	_repl, ok := replaces[lang]
	if ok {
		lang = _repl
	}

	var isValidLang bool
	for _, lng := range l.langs {
		if lang == lng {
			isValidLang = true
		}
	}

	if !isValidLang {
		return
	}

	wrapper, err := l.execFact.NewExecutor(e.GuildID)
	if err != nil || wrapper == nil {
		return
	}

	err = s.MessageReactionAdd(e.ChannelID, e.ID, runReactionEmoji)
	if err != nil {
		return
	}

	jdMsg := &jdoodleMessage{
		Message: e,
		wrapper: wrapper,
		lang:    lang,
		script:  cont,
		embLang: embLang,
	}

	l.msgMap.Set(e.ID, jdMsg, removeHandlerTimeout, func(v interface{}) {
		s.MessageReactionRemove(e.ChannelID, e.ID, runReactionEmoji, s.State.User.ID)
	})
}

func (l *ListenerJdoodle) HandlerReactionAdd(s *discordgo.Session, eReact *discordgo.MessageReactionAdd) {
	if eReact.UserID == s.State.User.ID {
		return
	}

	if eReact.Emoji.Name != runReactionEmoji {
		return
	}

	jdMsg, ok := l.msgMap.GetValue(eReact.MessageID).(*jdoodleMessage)
	if !ok || jdMsg == nil {
		return
	}

	allowed, err := l.checkPermission(s, eReact.GuildID, eReact.UserID)
	if !allowed || !l.checkLimit(eReact.UserID) {
		s.MessageReactionRemove(eReact.ChannelID, eReact.MessageID, eReact.Emoji.Name, eReact.UserID)
		return
	}

	s.MessageReactionsRemoveAll(eReact.ChannelID, eReact.MessageID)

	resMsg := util.SendEmbed(s, eReact.ChannelID, "Executing...", "", static.ColorEmbedGray)
	if resMsg.Error() != nil {
		return
	}

	result, err := jdMsg.wrapper.Exec(codeexec.Payload{
		Language: jdMsg.lang,
		Code:     jdMsg.script,
	})

	if err != nil {
		s.ChannelMessageEditEmbed(resMsg.ChannelID, resMsg.ID, &discordgo.MessageEmbed{
			Color:       static.ColorEmbedError,
			Title:       "Execution Error",
			Description: fmt.Sprintf("API responded with following error: ```\n%s\n```", err.Error()),
		})
		discordutil.DeleteMessageLater(s, resMsg.Message, 15*time.Second)
	} else {
		executor, _ := s.GuildMember(eReact.GuildID, eReact.UserID)

		footer := strings.ToUpper(jdMsg.lang)
		if executor != nil {
			footer += " | Executed by " + executor.User.String()
		}

		emb := embedbuilder.New().
			WithColor(static.ColorEmbedCyan).
			WithTitle("Compilation Result").
			WithFooter(footer, "", "").
			AddField("Code", fmt.Sprintf("```%s\n%s\n```", jdMsg.embLang, jdMsg.script))

		if result.StdOut != "" {
			emb.AddField("StdOut", "```\n"+result.StdOut+"\n```")
		}

		if result.StdErr != "" {
			emb.AddField("StdErr", "```\n"+result.StdErr+"\n```")
		}

		if result.CpuUsed != "" {
			emb.AddField("CPU Time", result.CpuUsed, true)
		}
		if result.MemUsed != "" {
			emb.AddField("Memory", result.MemUsed, true)
		}

		s.ChannelMessageEditEmbed(resMsg.ChannelID, resMsg.ID, emb.Build())

		l.msgMap.Remove(jdMsg.ID)
	}
}

func (l *ListenerJdoodle) parseMessageContent(content string) (lang string, script string, ok bool) {
	spl := strings.Split(content, "```")
	if len(spl) < 3 {
		return
	}

	inner := spl[1]
	iFirstLineBreak := strings.Index(inner, "\n")
	if iFirstLineBreak < 0 || len(inner)+1 <= iFirstLineBreak {
		return
	}

	lang = inner[:iFirstLineBreak]
	script = inner[iFirstLineBreak+1:]
	ok = len(lang) > 0 && len(script) > 0

	return
}

func (l *ListenerJdoodle) checkLimit(userID string) bool {
	limiter, ok := l.limits.GetValue(userID).(*rate.Limiter)
	if !ok || limiter == nil {
		limiter = rate.NewLimiter(rate.Limit(limitRate), limitBurst)
		l.limits.Set(userID, limiter, limitTMLifetime)
	}

	return limiter.Allow()
}

func (l *ListenerJdoodle) checkPermission(s *discordgo.Session, guildID, userID string) (bool, error) {
	allowed, _, err := l.pmw.CheckPermissions(s, guildID, userID, "sp.chat.exec.exec")
	return allowed, err
}
