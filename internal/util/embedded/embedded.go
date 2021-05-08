package embedded

import (
	_ "embed"
	"strings"
)

var (
	//go:embed AppVersion.txt
	AppVersion string
	//go:embed AppCommit.txt
	AppCommit string
	//go:embed AppDate.txt
	AppDate string
	//go:embed Release.txt
	Release string
)

func IsRelease() bool {
	return strings.ToLower(Release) == "true"
}
