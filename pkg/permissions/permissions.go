// Package permissions provides functionalities to
// calculate, update and merge arrays of permission
// domain rules.
//
// Read this to get more information about how
// permission domains and rules are working:
// https://github.com/zekroTJA/shinpuru/wiki/Permissions-Guide
package permissions

import (
	"fmt"
	"strings"
)

const maxPermIndex = 999

// permissionMatchDNs tries to match the passed
// domainName on the passed perm.
//
// This also respects explicit domainNames
// prefixed with '!'.
//
// The resulting match index is returned. If the
// match index is < 0, this must be interpreted as
// no match.
func permissionMatchDNs(domainName, perm string) int {
	if domainName == "" {
		return -1
	}

	var needsExplicitAllow bool

	// A domainName with the prefix '!' sets
	// needsExplicitAllow to true.
	// This means, the domainName must be
	// explicitely allowed and can not be matched
	// by wildcard.
	if domainName[0] == '!' {
		needsExplicitAllow = true
		domainName = domainName[1:]
	}

	// If the domain name equals perm, return
	// 999 match index.
	if domainName == perm {
		return maxPermIndex
	}

	// ...otherwise, if needsExplicitAllow is
	// true and it is not an exact match,
	// return negative match.
	if needsExplicitAllow {
		return -1
	}

	// Split domainName in areas seperated by '.'
	dnAreas := strings.Split(domainName, ".")
	assembled := ""
	for i, dnArea := range dnAreas {
		if assembled == "" {
			// If assembled is empty, set assembled to
			// current dnArea.
			assembled = dnArea
		} else {
			// Otherwise, add current dnArea to assembled.
			assembled = fmt.Sprintf("%s.%s", assembled, dnArea)
		}

		// If perm equals assembled area with trailing
		// wildcard selector ".*", return current index
		// as match index.
		if perm == fmt.Sprintf("%s.*", assembled) {
			return i
		}
	}

	// Otherwise, return negative match index.
	return -1
}

// permissionCheckDNs tries to match domainName on the
// passed perm and returns the match index and if
// it matched and perm is not prefixed with '-'.
func permissionCheckDNs(domainName, perm string) (int, bool) {
	if perm == "" {
		return -1, false
	}

	if perm[0] != '+' && perm[0] != '-' {
		return -1, false
	}

	match := permissionMatchDNs(domainName, perm[1:])
	if match < 0 {
		return match, false
	}

	return match, !strings.HasPrefix(perm, "-")
}
