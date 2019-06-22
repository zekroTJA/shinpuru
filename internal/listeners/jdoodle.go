package listeners

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/zekroTJA/timedmap"
	"golang.org/x/time/rate"

	"github.com/zekroTJA/shinpuru/internal/util"

	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/internal/core"
)

const (
	removeHandlerTimeout   = 3 * time.Minute
	limitTMCleanupInterval = 30 * time.Second // 10 * time.Minute
	limitTMLifetime        = 24 * time.Hour
	limitBurst             = 2                // 5
	limitRate              = float64(1) / 900 // one token per 15 minutes
)

var (
	runReactionEmoji = "▶"

	embRx = regexp.MustCompile("```([\\w\\+\\#]+)(\\n|\\s)((.|\\n)*)```")

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
	db     core.Database
	limits *timedmap.TimedMap
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

func NewListenerJdoodle(db core.Database) *ListenerJdoodle {
	return &ListenerJdoodle{
		db:     db,
		limits: timedmap.New(limitTMCleanupInterval),
	}
}

func (l *ListenerJdoodle) Handler(s *discordgo.Session, e *discordgo.MessageCreate) {
	if e.Author.Bot || e.GuildID == "" {
		return
	}

	if !embRx.MatchString(e.Content) {
		return
	}

	_matches := embRx.FindAllStringSubmatch(e.Content, 1)
	if len(_matches) < 1 {
		return
	}
	matches := _matches[0]

	if len(matches) < 4 {
		return
	}

	lang := strings.ToLower(strings.Trim(matches[1], " \t"))
	cont := strings.Trim(matches[3], " \t")
	embLang := lang

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
		fmt.Println(err)
		return
	}

	var removeHandler func()

	removeHandler = s.AddHandler(func(_ *discordgo.Session, eReact *discordgo.MessageReactionAdd) {
		if eReact.UserID == s.State.User.ID || eReact.GuildID != e.GuildID || eReact.MessageID != e.ID {
			return
		}

		if eReact.Emoji.Name != runReactionEmoji {
			return
		}

		if !l.checkLimit(eReact.UserID) {
			s.MessageReactionRemove(eReact.ChannelID, eReact.MessageID, eReact.Emoji.Name, eReact.UserID)
			return
		}

		s.MessageReactionsRemoveAll(eReact.ChannelID, eReact.MessageID)
		removeHandler()

		resMsg, err := util.SendEmbed(s, eReact.ChannelID, "Executing...", "", util.ColorEmbedGray)
		if err != nil {
			return
		}

		requestBody := &jdoodleRequestBody{
			ClientID:     jdCredsSplit[0],
			ClientSecret: jdCredsSplit[1],
			Script:       cont,
			Language:     lang,
		}

		bodyBuffer, err := json.Marshal(requestBody)
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
				Color:       util.ColorEmbedError,
				Title:       "Execution Error",
				Description: fmt.Sprintf("API responded with following error: ```\nCode: %d\nMsg:  %s\n```", res.StatusCode, errBody.Error),
			})
			util.DeleteMessageLater(s, resMsg, 15*time.Second)

		} else {

			result := new(jdoodleResponseResult)
			err = dec.Decode(result)

			if err != nil {
				unexpectedError(s, eReact.ChannelID, err)
				return
			}

			executor, _ := s.GuildMember(eReact.GuildID, eReact.UserID)

			emb := &discordgo.MessageEmbed{
				Color: util.ColorEmbedCyan,
				Title: "Compilation Result",
				Fields: []*discordgo.MessageEmbedField{
					&discordgo.MessageEmbedField{
						Name:  "Code",
						Value: fmt.Sprintf("```%s\n%s\n```", embLang, cont),
					},
					&discordgo.MessageEmbedField{
						Name:  "Output",
						Value: "```\n" + result.Output + "\n```",
					},
					&discordgo.MessageEmbedField{
						Name:   "CPU Time",
						Value:  result.CPUTime + " Seconds",
						Inline: true,
					},
					&discordgo.MessageEmbedField{
						Name:   "Memory",
						Value:  result.Memory + " Byte",
						Inline: true,
					},
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: strings.ToUpper(lang),
				},
			}

			if executor != nil {
				emb.Footer.Text += " | Executed by " + executor.User.String()
			}

			s.ChannelMessageEditEmbed(resMsg.ChannelID, resMsg.ID, emb)
		}
	})

	time.AfterFunc(removeHandlerTimeout, func() {
		s.MessageReactionsRemoveAll(e.ChannelID, e.ID)
		removeHandler()
	})
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
