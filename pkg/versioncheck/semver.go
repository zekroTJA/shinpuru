package versioncheck

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var svRx = regexp.MustCompile(`^[vV]?\.?(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:[\-+]([\w_\-\.]+))?$`)

var (
	ErrNoMatch = errors.New("no match")
)

// CompareType specifies the type of comparison used for
// different semantic versions.
type CompareType int

const (
	Major CompareType = iota
	Minor
	Patch
	Exact
)

// Semver represents the elements of a
// semantic version.
// See: https://semver.org/spec/v2.0.0.html
type Semver struct {
	Major      int    `json:"major"`
	Minor      int    `json:"minor"`
	Patch      int    `json:"patch"`
	Attachment string `json:"attachment"`
}

// ParseSemver takes a semantic version as string
// and tries to parse it into a Semver instance.
func ParseSemver(raw string) (v Semver, err error) {
	res := svRx.FindAllStringSubmatch(raw, -1)

	if len(res) == 0 || len(res[0]) != 5 {
		err = ErrNoMatch
		return
	}

	r := res[0]

	if v.Major, err = atoiSafe(r[1]); err != nil {
		return
	}
	if v.Minor, err = atoiSafe(r[2]); err != nil {
		return
	}
	if v.Patch, err = atoiSafe(r[3]); err != nil {
		return
	}
	v.Attachment = r[4]

	return
}

// String returns the string representation
// of the current Semver.
func (v Semver) String() (r string) {
	var sb strings.Builder

	sb.WriteString(strconv.Itoa(v.Major))
	sb.WriteRune('.')
	sb.WriteString(strconv.Itoa(v.Minor))
	sb.WriteRune('.')
	sb.WriteString(strconv.Itoa(v.Patch))

	if v.Attachment != "" {
		sb.WriteRune('-')
		sb.WriteString(v.Attachment)
	}

	return sb.String()
}

// Equal returns true if the given semver equals
// the current semver down to the given compare type.
func (v1 Semver) Equal(v2 Semver, et ...CompareType) (ok bool) {
	return v1.compare(v2, func(i, j int) bool {
		return i != j
	}, et)
}

// OlderThan returns true if the given semver is older than
// the current semver down to the given compare type.
func (v1 Semver) OlderThan(v2 Semver, et ...CompareType) (ok bool) {
	return !v1.Equal(v2, et...) && v1.compare(v2, func(i, j int) bool {
		return i > j
	}, et)
}

// LaterThan returns true if the given semver is later than
// the current semver down to the given compare type.
func (v1 Semver) LaterThan(v2 Semver, et ...CompareType) (ok bool) {
	return !v1.Equal(v2, et...) && v1.compare(v2, func(i, j int) bool {
		return i < j
	}, et)
}

func (v1 Semver) compare(
	v2 Semver,
	cf func(i, j int) bool,
	et []CompareType,
) (ok bool) {
	e := optEq(et)

	switch e {
	case Exact:
		if v1.Attachment != v2.Attachment {
			return
		}
		fallthrough
	case Patch:
		if cf(v1.Patch, v2.Patch) {
			return
		}
		fallthrough
	case Minor:
		if cf(v1.Minor, v2.Minor) {
			return
		}
		fallthrough
	case Major:
		if cf(v1.Major, v2.Major) {
			return
		}
	}

	ok = true
	return
}

func atoiSafe(s string) (i int, err error) {
	if s == "" {
		return
	}
	i, err = strconv.Atoi(s)
	return
}

func optEq(e []CompareType) CompareType {
	if len(e) != 0 {
		return e[0]
	}
	return Exact
}
