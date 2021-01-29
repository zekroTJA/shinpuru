package permissions

import "testing"

func TestUpdateAdd(t *testing.T) {
	p1 := PermissionArray{
		"+a.a",
		"-a.b",
	}

	pNew, changed := p1.Update("+a.c", false)
	if !changed {
		t.Error("did not changed")
	}
	if !equalsUnordered(
		pNew,
		PermissionArray{
			"+a.a",
			"-a.b",
			"+a.c",
		},
	) {
		t.Error("unexprected update result")
	}

	pNew, changed = p1.Update("-a.c", false)
	if !changed {
		t.Error("did not changed")
	}
	if !equalsUnordered(
		pNew,
		PermissionArray{
			"+a.a",
			"-a.b",
			"-a.c",
		},
	) {
		t.Error("unexprected update result")
	}

	pNew, changed = p1.Update("+a.a", false)
	if changed {
		t.Error("did changed")
	}
	if !equalsUnordered(
		pNew,
		p1,
	) {
		t.Error("unexprected update result")
	}

	pNew, changed = p1.Update("+a.b", false)
	if !changed {
		t.Error("did not changed")
	}
	if !equalsUnordered(
		pNew,
		PermissionArray{
			"+a.a",
		},
	) {
		t.Error("unexprected update result")
	}

	pNew, changed = p1.Update("-a.a", false)
	if !changed {
		t.Error("did not changed")
	}
	if !equalsUnordered(
		pNew,
		PermissionArray{
			"-a.b",
		},
	) {
		t.Error("unexprected update result")
	}

	pNew, changed = p1.Update("-a.a", false)
	if !changed {
		t.Error("did not changed")
	}
	if !equalsUnordered(
		pNew,
		PermissionArray{
			"-a.b",
		},
	) {
		t.Error("unexprected update result")
	}

	pNew, changed = p1.Update("-a.a", false)
	if !changed {
		t.Error("did not changed")
	}
	if !equalsUnordered(
		pNew,
		PermissionArray{
			"-a.b",
		},
	) {
		t.Error("unexprected update result")
	}

	pNew, changed = p1.Update("-a.a", true)
	if !changed {
		t.Error("did not changed")
	}
	if !equalsUnordered(
		pNew,
		PermissionArray{
			"-a.a",
			"-a.b",
		},
	) {
		t.Error("unexprected update result")
	}

	pNew, changed = p1.Update("+a.b", true)
	if !changed {
		t.Error("did not changed")
	}
	if !equalsUnordered(
		pNew,
		PermissionArray{
			"+a.a",
			"+a.b",
		},
	) {
		t.Error("unexprected update result")
	}

	pNew, changed = p1.Update("-a.c", true)
	if !changed {
		t.Error("did not changed")
	}
	if !equalsUnordered(
		pNew,
		PermissionArray{
			"+a.a",
			"-a.b",
			"-a.c",
		},
	) {
		t.Error("unexprected update result")
	}

	pNew, changed = p1.Update("+a.a", true)
	if changed {
		t.Error("did changed")
	}
	if !equalsUnordered(
		pNew,
		PermissionArray{
			"+a.a",
			"-a.b",
		},
	) {
		t.Error("unexprected update result")
	}
}

func TestMerge(t *testing.T) {
	p1 := PermissionArray{
		"+a.a",
		"-a.b",
	}

	if !equalsUnordered(
		p1.Merge(PermissionArray{
			"+a.b",
			"+a.c",
		}, false),
		PermissionArray{
			"+a.a",
			"+a.c",
		},
	) {
		t.Error("unexprected update result")
	}

	if !equalsUnordered(
		p1.Merge(PermissionArray{
			"+a.c",
			"+a.d",
		}, false),
		PermissionArray{
			"+a.a",
			"-a.b",
			"+a.c",
			"+a.d",
		},
	) {
		t.Error("unexprected update result")
	}

	if !equalsUnordered(
		p1.Merge(PermissionArray{
			"-a.a",
			"+a.b",
		}, false),
		PermissionArray{},
	) {
		t.Error("unexprected update result")
	}

	if !equalsUnordered(
		p1.Merge(PermissionArray{
			"-a.a",
			"+a.c",
		}, true),
		PermissionArray{
			"-a.a",
			"-a.b",
			"+a.c",
		},
	) {
		t.Error("unexprected update result")
	}

	if !equalsUnordered(
		p1.Merge(PermissionArray{
			"-a.a",
			"+a.b",
			"+a.c",
			"-a.d",
		}, true),
		PermissionArray{
			"-a.a",
			"+a.b",
			"+a.c",
			"-a.d",
		},
	) {
		t.Error("unexprected update result")
	}
}

func TestEquals(t *testing.T) {
	p1 := PermissionArray{
		"+a.a",
		"-a.b",
	}
	if !p1.Equals(p1) {
		t.Error("equal arrays have unequal res")
	}

	p1 = PermissionArray{}
	if !p1.Equals(p1) {
		t.Error("equal arrays have unequal res")
	}

	p1 = PermissionArray{
		"+a.a",
		"-a.b",
	}
	p2 := PermissionArray{
		"-a.b",
		"+a.a",
	}
	if p1.Equals(p2) {
		t.Error("unequal arrays have equal res")
	}

	p1 = PermissionArray{
		"+a.a",
		"-a.b",
	}
	p2 = PermissionArray{
		"-a.b",
	}
	if p1.Equals(p2) {
		t.Error("unequal arrays have equal res")
	}
}

func TestCheck(t *testing.T) {
	p := PermissionArray{
		"+a.a",
		"+a.b.*",
		"-a.c",
		"-a.d.*",
	}
	if !p.Check("a.a") {
		t.Error("check failed")
	}

	if !p.Check("a.b.c") {
		t.Error("check failed")
	}

	if p.Check("a.c") {
		t.Error("check failed")
	}

	if p.Check("a.d.c") {
		t.Error("check failed")
	}

	if p.Check("a.d") {
		t.Error("check failed")
	}

	if p.Check("x.y.z") {
		t.Error("check failed")
	}

	if p.Check("x") {
		t.Error("check failed")
	}

	if p.Check("") {
		t.Error("check failed")
	}
}

// --- HELPER ---------

func equalsUnordered(p1, p2 PermissionArray) bool {
	for _, v1 := range p1 {
		contains := false
		for _, v2 := range p2 {
			if v1 == v2 {
				contains = true
			}
		}
		if !contains {
			return false
		}
	}

	for _, v2 := range p2 {
		contains := false
		for _, v1 := range p1 {
			if v1 == v2 {
				contains = true
			}
		}
		if !contains {
			return false
		}
	}

	return true
}
