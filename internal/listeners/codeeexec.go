package listeners

import (
	"fmt"
	"strings"
	"time"

	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/timedmap"
	"github.com/zekrotja/dgrs"
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

	replaces = map[string]string{
		"js":         "nodejs",
		"javascript": "nodejs",
		"c++":        "cpp",
		"c#":         "csharp",
		"python":     "python3",
		"py":         "python3",
	}
)

type ListenerCodeexec struct {
	db       database.Database
	execFact codeexec.Factory
	pmw      *middleware.PermissionsMiddleware
	st       *dgrs.State

	langs  []string
	limits *timedmap.TimedMap
	msgMap *timedmap.TimedMap
}

type execMessage struct {
	*discordgo.Message

	wrapper codeexec.Executor
	lang    string
	script  string

	embLang string
}

func NewListenerJdoodle(container di.Container) (l *ListenerCodeexec, err error) {
	l = &ListenerCodeexec{}

	l.db = container.Get(static.DiDatabase).(database.Database)
	l.pmw = container.Get(static.DiPermissionMiddleware).(*middleware.PermissionsMiddleware)
	l.execFact = container.Get(static.DiCodeExecFactory).(codeexec.Factory)
	l.st = container.Get(static.DiState).(*dgrs.State)
	l.limits = timedmap.New(limitTMCleanupInterval)
	l.msgMap = timedmap.New(removeHandlerCleanupInterval)

	l.langs, err = l.execFact.Languages()

	return
}

func (l *ListenerCodeexec) HandlerMessageCreate(s *discordgo.Session, e *discordgo.MessageCreate) {
	l.handler(s, e.Message)
}

func (l *ListenerCodeexec) HandlerMessageUpdate(s *discordgo.Session, e *discordgo.MessageUpdate) {
	l.handler(s, e.Message)
}

func (l *ListenerCodeexec) handler(s *discordgo.Session, e *discordgo.Message) {
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

	jdMsg := &execMessage{
		Message: e,
		wrapper: wrapper,
		lang:    lang,
		script:  cont,
		embLang: embLang,
	}

	self, err := l.st.SelfUser()
	if err != nil {
		return
	}

	l.msgMap.Set(e.ID, jdMsg, removeHandlerTimeout, func(v interface{}) {
		s.MessageReactionRemove(e.ChannelID, e.ID, runReactionEmoji, self.ID)
	})
}

func (l *ListenerCodeexec) HandlerReactionAdd(s *discordgo.Session, eReact *discordgo.MessageReactionAdd) {
	self, err := l.st.SelfUser()
	if err != nil {
		return
	}

	if eReact.UserID == self.ID {
		return
	}

	if eReact.Emoji.Name != runReactionEmoji {
		return
	}

	jdMsg, ok := l.msgMap.GetValue(eReact.MessageID).(*execMessage)
	if !ok || jdMsg == nil {
		return
	}

	if eReact.UserID != jdMsg.Author.ID {
		s.MessageReactionRemove(eReact.ChannelID, eReact.MessageID, eReact.Emoji.Name, eReact.UserID)
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

		if l.execFact.Name() == "ranna" {
			emb.WithDescription("*Code execution is provided by [ranna](https://github.com/ranna-go).*")
		}

		if result.StdOut != "" {
			emb.AddField("StdOut", "```\n"+result.StdOut+"\n```")
		}
		if result.StdErr != "" {
			emb.AddField("StdErr", "```\n"+result.StdErr+"\n```")
		}
		if result.CpuUsed != "" {
			emb.AddInlineField("CPU Time", result.CpuUsed)
		}
		if result.MemUsed != "" {
			emb.AddInlineField("Memory", result.MemUsed)
		}
		if result.ExecTime != 0 {
			emb.AddInlineField("Execution Time", result.ExecTime.Round(time.Millisecond).String())
		}

		s.ChannelMessageEditEmbed(resMsg.ChannelID, resMsg.ID, emb.Build())

		l.msgMap.Remove(jdMsg.ID)
	}
}

func (l *ListenerCodeexec) parseMessageContent(content string) (lang string, script string, ok bool) {
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

func (l *ListenerCodeexec) checkLimit(userID string) bool {
	limiter, ok := l.limits.GetValue(userID).(*rate.Limiter)
	if !ok || limiter == nil {
		limiter = rate.NewLimiter(rate.Limit(limitRate), limitBurst)
		l.limits.Set(userID, limiter, limitTMLifetime)
	}

	return limiter.Allow()
}

func (l *ListenerCodeexec) checkPermission(s *discordgo.Session, guildID, userID string) (bool, error) {
	allowed, _, err := l.pmw.CheckPermissions(s, guildID, userID, "sp.chat.exec.exec")
	return allowed, err
}
