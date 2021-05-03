package codeexec

import (
	"strings"

	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/jdoodle"
)

var langs = []string{"java", "c", "cpp", "c99", "cpp14", "php", "perl", "python3", "ruby", "go", "scala", "bash", "sql", "pascal", "csharp",
	"vbn", "haskell", "objc", "ell", "swift", "groovy", "fortran", "brainfuck", "lua", "tcl", "hack", "rust", "d", "ada", "r", "freebasic",
	"verilog", "cobol", "dart", "yabasic", "clojure", "nodejs", "scheme", "forth", "prolog", "octave", "coffeescript", "icon", "fsharp", "nasm",
	"gccasm", "intercal", "unlambda", "picolisp", "spidermonkey", "rhino", "bc", "clisp", "elixir", "factor", "falcon", "fantom", "pike", "smalltalk",
	"mozart", "lolcode", "racket", "kotlin"}

type JdoodleFactory struct {
	db database.Database
}

func NewJdoodleFactory(container di.Container) (e *JdoodleFactory) {
	e = &JdoodleFactory{}

	e.db = container.Get(static.DiDatabase).(database.Database)

	return
}

func (e *JdoodleFactory) Name() string {
	return "jdoodle"
}

func (e *JdoodleFactory) Languages() ([]string, error) {
	return langs, nil
}

func (e *JdoodleFactory) NewExecutor(guildID string) (exec Executor, err error) {
	jdCreds, err := e.db.GetGuildJdoodleKey(guildID)
	if err != nil || jdCreds == "" {
		return
	}

	jdCredsSplit := strings.Split(jdCreds, "#")
	if len(jdCredsSplit) < 2 {
		return
	}

	exec = &JdoodleExecutor{jdCredsSplit[0], jdCredsSplit[1]}
	return
}

type JdoodleExecutor struct {
	clientId     string
	clientSecret string
}

func (e *JdoodleExecutor) Exec(p Payload) (res Response, err error) {
	w := jdoodle.NewWrapper(e.clientId, e.clientSecret)
	r, err := w.ExecuteScript(p.Language, p.Code)
	if err != nil {
		return
	}

	res.StdOut = r.Output
	res.MemUsed = r.Memory + " Byte"
	res.CpuUsed = r.CPUTime + " Seconds"

	return
}
