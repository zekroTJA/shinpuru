package permissions

import "testing"

func TestPermissionsMatchDNs(t *testing.T) {
	check := func(dn, perm string, exp int) {
		if res := permissionMatchDNs(dn, perm); res != exp {
			t.Errorf("%s -> %s : result was %d (expexted %d)",
				perm, dn, res, exp)
		}
	}

	check("test.sub.subsub", "test.sub.subsub", maxPermIndex)
	check("test.sub.subsub", "test.sub.*", 1)
	check("test.sub.subsub", "test.*", 0)

	check("test.sub.subsub", "", -1)
	check("", "test.sub", -1)

	check("test.sub.subsub", "test.sub.notsubsub", -1)
	check("test.sub.subsub", "test.notsub", -1)
	check("test.sub.subsub", "nottest", -1)
	check("test.sub.subsub", "test.notsub.subsub", -1)
	check("test.sub.subsub", "test.subsub.*", -1)
	check("test.sub.subsub", "nottest.*", -1)

	check("!test.sub.subsub", "test.sub.subsub", maxPermIndex)
	check("!test.sub.subsub", "test.sub.*", -1)
	check("!test.sub.subsub", "test.subsub", -1)
}

func TestPermissionCheckDNs(t *testing.T) {
	check := func(dn, perm string, exp1 int, exp2 bool) {
		if res1, res2 := permissionCheckDNs(dn, perm); res1 != exp1 || res2 != exp2 {
			t.Errorf("%s -> %s : result was (%d, %v) (expexted (%d, %v))",
				perm, dn, res1, res2, exp1, exp2)
		}
	}

	check("test.sub.subsub", "+test.sub.subsub", maxPermIndex, true)
	check("test.sub.subsub", "-test.sub.subsub", maxPermIndex, false)

	check("test.sub.subsub", "+test.sub.notsubsub", -1, false)
	check("test.sub.subsub", "-test.sub.notsubsub", -1, false)

	check("test.sub.subsub", "+test.sub.*", 1, true)
	check("test.sub.subsub", "-test.sub.*", 1, false)

	check("test.sub.subsub", "test.sub.subsub", -1, false)

	check("test.sub.subsub", "", -1, false)
	check("", "test.sub", -1, false)
}
