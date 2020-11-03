package listeners

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/zekroTJA/timedmap"
	"golang.org/x/time/rate"

	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"

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

	langs = []string{"java", "c", "cpp", "c99", "cpp14", "php", "perl", "python3", "ruby", "go", "scala", "bash", "sql", "pascal", "csharp",
		"vbn", "haskell", "objc", "ell", "swift", "groovy", "fortran", "brainfuck", "lua", "tcl", "hack", "rust", "d", "ada", "r", "freebasic",
		"verilog", "cobol", "dart", "yabasic", "clojure", "nodejs", "scheme", "forth", "prolog", "octave", "coffeescript", "icon", "fsharp", "nasm",
		"gccasm", "intercal", "unlambda", "picolisp", "spidermonkey", "rhino", "bc", "clisp", "elixir", "factor", "falcon", "fantom", "pike", "smalltalk",
		"mozart", "lolcode", "racket", "kotlin"}

	replaces = map[string]string{
		"js":         "nodejs",
		"javascript": "nodejs",
		"c++":        "cpp",
		"c#":         "csharp",
		"python":     "python3",
		"py":         "python3",
	}

	apiURL = "https://api.jdoodle.com/v1/execute"
)

type ListenerJdoodle struct {
	db     database.Database
	limits *timedmap.TimedMap
	msgMap *timedmap.TimedMap
}

type jdoodleMessage struct {
	*discordgo.Message

	rb *jdoodleRequestBody

	embLang string
}

type jdoodleRequestBody struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Script       string `json:"script"`
	Language     string `json:"language"`
}

type jdoodleResponseError struct {
	Error string `json:"error"`
}

type jdoodleResponseResult struct {
	Output  string `json:"output"`
	Memory  string `json:"memory"`
	CPUTime string `json:"cpuTime"`
}

func NewListenerJdoodle(db database.Database) *ListenerJdoodle {
	return &ListenerJdoodle{
		db:     db,
		limits: timedmap.New(limitTMCleanupInterval),
		msgMap: timedmap.New(removeHandlerCleanupInterval),
	}
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

	fmt.Println("<" + cont + ">")

	if lang == "" || cont == "" {
		return
	}

	_repl, ok := replaces[lang]
	if ok {
		lang = _repl
	}

	var isValidLang bool
	for _, l := range langs {
		if lang == l {
			isValidLang = true
		}
	}

	if !isValidLang {
		return
	}

	jdCreds, err := l.db.GetGuildJdoodleKey(e.GuildID)
	if err != nil || jdCreds == "" {
		return
	}

	jdCredsSplit := strings.Split(jdCreds, "#")
	if len(jdCredsSplit) < 2 {
		return
	}

	err = s.MessageReactionAdd(e.ChannelID, e.ID, runReactionEmoji)
	if err != nil {
		return
	}

	jdMsg := &jdoodleMessage{
		Message: e,
		rb: &jdoodleRequestBody{
			ClientID:     jdCredsSplit[0],
			ClientSecret: jdCredsSplit[1],
			Script:       cont,
			Language:     lang,
		},
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

	if !l.checkLimit(eReact.UserID) {
		s.MessageReactionRemove(eReact.ChannelID, eReact.MessageID, eReact.Emoji.Name, eReact.UserID)
		return
	}

	s.MessageReactionsRemoveAll(eReact.ChannelID, eReact.MessageID)

	resMsg := util.SendEmbed(s, eReact.ChannelID, "Executing...", "", static.ColorEmbedGray)
	if resMsg.Error() != nil {
		return
	}

	bodyBuffer, err := json.Marshal(jdMsg.rb)
	if err != nil {
		unexpectedError(s, eReact.ChannelID, err)
		return
	}
	req, err := http.NewRequest("POST", apiURL, bytes.NewReader(bodyBuffer))
	if err != nil {
		unexpectedError(s, eReact.ChannelID, err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		unexpectedError(s, eReact.ChannelID, err)
		return
	}

	dec := json.NewDecoder(res.Body)
	if res.StatusCode != 200 {
		errBody := new(jdoodleResponseError)
		err = dec.Decode(errBody)

		if err != nil {
			unexpectedError(s, eReact.ChannelID, err)
			return
		}

		s.ChannelMessageEditEmbed(resMsg.ChannelID, resMsg.ID, &discordgo.MessageEmbed{
			Color:       static.ColorEmbedError,
			Title:       "Execution Error",
			Description: fmt.Sprintf("API responded with following error: ```\nCode: %d\nMsg:  %s\n```", res.StatusCode, errBody.Error),
		})
		discordutil.DeleteMessageLater(s, resMsg.Message, 15*time.Second)

	} else {

		result := new(jdoodleResponseResult)
		err = dec.Decode(result)

		if err != nil {
			unexpectedError(s, eReact.ChannelID, err)
			return
		}

		executor, _ := s.GuildMember(eReact.GuildID, eReact.UserID)

		emb := &discordgo.MessageEmbed{
			Color: static.ColorEmbedCyan,
			Title: "Compilation Result",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Code",
					Value: fmt.Sprintf("```%s\n%s\n```", jdMsg.embLang, jdMsg.rb.Script),
				},
				{
					Name:  "Output",
					Value: "```\n" + result.Output + "\n```",
				},
				{
					Name:   "CPU Time",
					Value:  result.CPUTime + " Seconds",
					Inline: true,
				},
				{
					Name:   "Memory",
					Value:  result.Memory + " Byte",
					Inline: true,
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: strings.ToUpper(jdMsg.rb.Language),
			},
		}

		if executor != nil {
			emb.Footer.Text += " | Executed by " + executor.User.String()
		}

		s.ChannelMessageEditEmbed(resMsg.ChannelID, resMsg.ID, emb)

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

func unexpectedError(s *discordgo.Session, chanID string, err error) {
	util.SendEmbedError(s, chanID, "An unexpected error occured. Please inform the host of this bot about that: ```\n"+err.Error()+"\n```")
}
