package versioncheck

// Provider implements a getter for version
// information from a specified source.
type Provider interface {
	GetLatestVersion() (v Semver, err error)
}
