package util

// These variables are set on compilation by using
// the -X ldflag.
// See this for more information:
// https://golang.org/cmd/link

var (
	AppVersion = "TESTING_BUILD"
	AppCommit  = "TESTING_BUILD"
	AppDate    = "0"
	Release    = "FALSE"
)
