package angularservice

import "io"

var (
	// EmptyOptions is an empty options instance.
	EmptyOptions = Options{}
)

// Options contains configuration properties
// for an Angular dev server command instance.
type Options struct {
	Stdout io.Writer `json:"-"`    // Stdout writer
	Stderr io.Writer `json:"-"`    // Stderr writer
	Cd     string    `json:"cd"`   // Working direcrtory (where your Angular files are)
	Port   int       `json:"port"` // Port to expose angular server at
	Args   []string  `json:"args"` // Additional args passed to the Angular CLI
}
