package permissions

type SubPermission struct {
	Term        string `json:"term"`
	Explicit    bool   `json:"explicit"`
	Description string `json:"description"`
}

type PermCommand interface {
	Domain() string
	SubDomains() []SubPermission
}
