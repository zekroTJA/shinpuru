package permissions

type SubPermission struct {
	Term        string
	Explicit    bool
	Description string
}

type PermCommand interface {
	Domain() string
	SubDomains() []SubPermission
}
