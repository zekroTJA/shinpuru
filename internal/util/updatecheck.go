package util

import (
	"github.com/zekroTJA/shinpuru/internal/util/embedded"
	"github.com/zekroTJA/shinpuru/pkg/versioncheck"
	"github.com/zekrotja/rogu/log"
)

var (
	versionProvider versioncheck.Provider = versioncheck.NewGitHubProvider("zekroTJA", "shinpuru")
	currVersion                           = mustCurrSemver()
)

func CheckForUpdate() (isOld bool, current, latest versioncheck.Semver) {
	if currVersion == nil {
		return
	}

	latest, err := versionProvider.GetLatestVersion()
	if err != nil {
		log.Error().Tag("VersionCheck").Err(err).Msg("Failed retrieving latest version")
		return
	}

	current = *currVersion
	isOld = currVersion.OlderThan(latest, versioncheck.Patch)
	return
}

func mustCurrSemver() *versioncheck.Semver {
	curr, err := versioncheck.ParseSemver(embedded.AppVersion)
	if err != nil {
		log.Error().Tag("VersionCheck").Err(err).Field("retrieved", embedded.AppVersion).Msg("Failed parsing current version - versioncheck skipped")
		return nil
	}
	return &curr
}
