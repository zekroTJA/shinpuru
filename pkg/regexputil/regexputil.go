// Package regexutil provides additional utility
// functions used with regular expressions.
package regexputil

import "regexp"

// FindNamedSubmatchMap tries to find named groups of a regular
// expression in the given string s and returns matches in a map
// with the name of the group as key and the match as value.
//
// If no match for the named group was found, there will be no
// entry in the map therefore.
func FindNamedSubmatchMap(re *regexp.Regexp, s string) map[string]string {
	results := make(map[string]string)
	matches := re.FindStringSubmatch(s)
	for i, name := range re.SubexpNames() {
		if i >= len(matches) {
			break
		}
		if i != 0 && i < len(matches) && name != "" && matches[i] != "" {
			results[name] = matches[i]
		}
	}
	return results
}
