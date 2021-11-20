package usercommands

import (
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/slashcommands"
	"github.com/zekrotja/ken"
)

type User struct {
	p slashcommands.User
}

var (
	_ ken.UserCommand         = (*User)(nil)
	_ permissions.PermCommand = (*User)(nil)
)

func (c *User) Name() string {
	return "userinfo"
}

func (c *User) Description() string {
	return c.p.Description()
}

func (c *User) Domain() string {
	return c.p.Domain()
}

func (c *User) SubDomains() []permissions.SubPermission {
	return c.p.SubDomains()
}

func (c *User) Run(ctx *ken.Ctx) (err error) {
	return c.p.Run(ctx)
}
