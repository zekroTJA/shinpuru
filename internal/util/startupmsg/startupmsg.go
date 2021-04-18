package startupmsg

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"time"

	"github.com/zekroTJA/shinpuru/internal/util"
)

//go:embed template.txt
var templateTxt string

type information struct {
	Appname   string
	Copyright string
	Version   string
	Commit    string
	Release   bool
}

func getInformation() information {
	return information{
		Appname: "shinpuru",
		Copyright: fmt.Sprintf("© %d Ringo Hoffmann (zekro Development)",
			time.Now().Year()),
		Version: util.AppVersion,
		Commit:  util.AppCommit,
		Release: util.IsRelease(),
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func Output(w io.Writer) {
	t, err := template.New("startupmsg").Parse(templateTxt)
	must(err)
	must(t.Execute(w, getInformation()))
}
