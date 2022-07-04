package regexputil

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindNamedSubmatchMap(t *testing.T) {
	re := regexp.MustCompile(`^(?P<length>\d+)m$`)
	results := FindNamedSubmatchMap(re, "12m")
	assert.Equal(t, map[string]string{"length": "12"}, results)

	re = regexp.MustCompile(`^(?:(?P<length>\d+)m)+$`)
	results = FindNamedSubmatchMap(re, "12m15m")
	assert.Equal(t, map[string]string{"length": "15"}, results)

	re = regexp.MustCompile(`^(?:\d+)m$`)
	results = FindNamedSubmatchMap(re, "42m")
	assert.Equal(t, map[string]string{}, results)

	re = regexp.MustCompile(`^(?:(?P<h>\d+):)?(?:(?P<m>\d+):)?(?P<s>\d+)$`)
	results = FindNamedSubmatchMap(re, "12:13:14")
	assert.Equal(t, map[string]string{"h": "12", "m": "13", "s": "14"}, results)

	re = regexp.MustCompile(`^(?:(?P<h>\d+):)?(?:(?P<m>\d+):)?(?P<s>\d+)$`)
	results = FindNamedSubmatchMap(re, "13:14")
	assert.Equal(t, map[string]string{"h": "13", "s": "14"}, results)

	re = regexp.MustCompile(`^(?:(?P<h>\d+):)?(?:(?P<m>\d+):)?(?P<s>\d+)$`)
	results = FindNamedSubmatchMap(re, "asd")
	assert.Equal(t, map[string]string{}, results)
}
