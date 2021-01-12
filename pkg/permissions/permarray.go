package permissions

// PermissionsArray describes a set of permission
// rules.
//
// Example:
//   +sp.guild.config.*
//   +sp.*
//   +sp.guild.*
//   -sp.guild.mod.ban
//   +sp.etc.*
//   +sp.chat.*
type PermissionArray []string

// Updates "adds" the passed newPerm to the permission array
// p by merging the permissions and returns the result as
// new permission array.
//
// This means, if p looks like following
//   +sp.guild.*
//   +sp.guild.mod.ban
// and newPerm is '-sp.guild.mod.ban', the
// returned permission array will be
//   +sp.guild.*
func (p PermissionArray) Update(newPerm string, override bool) (newPermsArray PermissionArray, changed bool) {
	newPermsArray = make(PermissionArray, len(p)+1)

	i := 0
	add := true
	for _, perm := range p {
		// If the permission rule equals and overrride
		// is true, set newPerm at this point to
		// newPermsArray.
		// Otherwise, if the prefix of perm and newPerm
		// are unequal and prefix of newPerm is '-',
		// newPerm is not being added.
		// If prefix of perm and newPerm are unequal and
		// prefix of newPerm is '+', newPerm will be added
		// to newPermsArray.
		//
		// Otherwise, perm is added to newPermArray.
		if len(perm) > 0 && perm[1:] == newPerm[1:] {
			add = false

			if override {
				newPermsArray[i] = newPerm
				i++
				continue
			}

			if perm[0] != newPerm[0] {
				continue
			}
		}

		newPermsArray[i] = perm
		i++
	}

	if add {
		newPermsArray[i] = newPerm
		i++
	}

	newPermsArray = newPermsArray[:i]

	changed = !p.Equals(newPermsArray)

	return
}

// Merge updates all entries of p using Update one
// by one with all entries of newPerms. Parameter
// override is passed to the Update function.
//
// A new permissions array is returned with the
// resulting permission rule set.
func (p PermissionArray) Merge(newPerms PermissionArray, override bool) PermissionArray {
	for _, cp := range newPerms {
		p, _ = p.Update(cp, override)
	}
	return p
}

// Equals returns true when p2 has the same elements
// in the same order as p.
func (p PermissionArray) Equals(p2 PermissionArray) bool {
	if len(p) != len(p2) {
		return false
	}

	for i, v := range p {
		if v != p2[i] {
			return false
		}
	}

	return true
}

// Check returns true if the passed domainName
// matches positively on the permission array p.
func (p PermissionArray) Check(domainName string) bool {
	lvl := -1
	allow := false

	for _, perm := range p {
		m, a := permissionCheckDNs(domainName, perm)
		if m > lvl {
			allow = a
			lvl = m
		}
	}

	return allow
}
