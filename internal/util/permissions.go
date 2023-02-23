package util

import (
	"fmt"
	"strings"

	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/ken"
	"github.com/zekrotja/sop"
)

var allPermissions sop.Enumerable[string]

func GetAllPermissions(cmdHandler *ken.Ken) sop.Enumerable[string] {
	if allPermissions != nil {
		return allPermissions
	}

	cmds := cmdHandler.GetCommandInfo()

	// Create a copy and wrap it into a set
	perms := sop.Set(append([]string{}, static.AdditionalPermissions...))

	for _, cmd := range cmds {
		rDomain := cmd.Implementations["Domain"]
		if len(rDomain) != 1 {
			continue
		}
		domain, ok := rDomain[0].(string)
		if !ok {
			continue
		}
		perms.Push(domain)

		rSubs := cmd.Implementations["SubDomains"]
		if len(rSubs) != 1 {
			continue
		}
		subs, ok := rSubs[0].([]permissions.SubPermission)
		if !ok {
			continue
		}
		for _, sub := range subs {
			var comb string
			if strings.HasPrefix(sub.Term, "/") {
				comb = sub.Term[1:]
			} else {
				comb = fmt.Sprintf("%s.%s", domain, sub.Term)
			}
			perms.Push(comb)
		}
	}

	allPermissions = perms
	return perms
}
