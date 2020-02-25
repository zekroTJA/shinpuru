package commands

type SubPermission struct {
	Term        string
	Explicit    bool
	Description string
}

type Command interface {
	GetInvokes() []string
	GetDescription() string
	GetHelp() string
	GetGroup() string
	GetDomainName() string
	GetSubPermissionRules() []SubPermission
	Exec(args *CommandArgs) error
}

const (
	GroupGlobalAdmin = "GLOBAL ADMIN"
	GroupGuildAdmin  = "GUILD ADMIN"
	GroupModeration  = "MODERATION"
	GroupFun         = "FUN"
	GroupGame        = "GAME"
	GroupChat        = "CHAT"
	GroupEtc         = "ETC"
	GroupGeneral     = "GENERAL"
	GroupGuildConfig = "GUILD CONFIG"
)
