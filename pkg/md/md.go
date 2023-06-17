// Package md provides some simple
// markdown fowmatting utilities.
package md

import "fmt"

func Bold[T ~string](v T) string {
	return fmt.Sprintf("**%s**", v)
}

func Italic[T ~string](v T) string {
	return fmt.Sprintf("*%s*", v)
}

func Code[T ~string](v T) string {
	return fmt.Sprintf("`%s`", v)
}

func CodeBlock[T ~string](v T, lang ...string) string {
	var l string
	if len(lang) > 0 {
		l = lang[0]
	}
	return fmt.Sprintf("```%s\n%s\n```", l, v)
}

func Underline[T ~string](v T) string {
	return fmt.Sprintf("__%s__", v)
}

func StrikeThrough[T ~string](v T) string {
	return fmt.Sprintf("~~%s~~", v)
}

func Spoiler[T ~string](v T) string {
	return fmt.Sprintf("||%s||", v)
}
