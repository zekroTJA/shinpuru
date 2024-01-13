package embedded

import (
	"embed"
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

	//go:embed webdist
	FrontendFiles embed.FS
)

func IsRelease() bool {
	return strings.ToLower(Release) == "true"
}
