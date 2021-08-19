// Package argp is a stupid simple flag (argument) parser
// which allows to parse flags without panicing when
// non-registered flags are passed.
package argp

var (
	singleton = New()

	Scan   = singleton.Scan
	String = singleton.String
	Bool   = singleton.Bool
	Int    = singleton.Int
	Float  = singleton.Float
	Args   = singleton.Args
	Help   = singleton.Help
)
