package messagecommands

import (
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/slashcommands"
	"github.com/zekrotja/ken"
)

type Quote struct {
	slashcommands.Quote
}

var (
	_ ken.MessageCommand      = (*Quote)(nil)
	_ permissions.PermCommand = (*Quote)(nil)
)

func (c *Quote) TypeMessage() {}

func (c *Quote) Name() string {
	return "quotemessage"
}
