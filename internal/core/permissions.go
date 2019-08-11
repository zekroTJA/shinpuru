package core

import (
	"fmt"
	"strings"
)

type PermissionArray []string

func (p PermissionArray) Update(npdn string) PermissionArray {
	newp := make(PermissionArray, len(p)+1)

	i := 0
	add := true
	for _, cp := range p {
		if len(cp) > 0 && cp[1:] == npdn[1:] {
			add = false
			if cp[0] != npdn[0] {
				continue
			}
		}
		newp[i] = cp
		i++
	}

	if add {
		newp[i] = npdn
		i++
	}

	return newp[:i]
}

func (p PermissionArray) Merge(np PermissionArray) PermissionArray {
	for _, cp := range np {
		p = p.Update(cp)
	}
	return p
}

func PermissionMatchDNs(dn, p string) int {
	if dn == p {
		return 999
	}

	dnA := strings.Split(dn, ".")
	ass := ""
	for i, dnP := range dnA {
		if ass == "" {
			ass = dnP
		} else {
			ass = fmt.Sprintf("%s.%s", ass, dnP)
		}

		if p == ass+".*" {
			return i
		}
	}

	return -1
}

func PermissionCheckDNs(dn, p string) (int, bool) {
	match := PermissionMatchDNs(dn, p[1:])
	if match < 0 {
		return match, false
	}

	return match, !strings.HasPrefix(p, "-")
}

func PermissionCheck(dn string, p PermissionArray) bool {
	lvl := -1
	allow := false

	for _, perm := range p {
		m, a := PermissionCheckDNs(dn, perm)
		if m > lvl {
			allow = a
			lvl = m
		}
	}

	return allow
}
