package commands

type Command interface {
	GetInvokes() []string
	GetDescription() string
	GetHelp() string
	GetGroup() string
	GetPermission() int
	SetPermission(permLvl int)
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
