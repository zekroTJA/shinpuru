package commands

// SubPermission wraps information about
// a sub permissions of commands.
type SubPermission struct {
	Term        string
	Explicit    bool
	Description string
}

// Command describes the functionalities of a
// command struct which can be registered
// in the CommandHandler.
type Command interface {
	// GetInvokes returns the unique strings udes to
	// call the command. The first invoke is the
	// primary command invoke and each following is
	// treated as command alias.
	GetInvokes() []string
	// GetDescription returns a brief description about
	// the functionality of the command.
	GetDescription() string
	// GetHelp returns detailed information on how to
	// use the command and their sub commands.
	GetHelp() string
	// GetGroup returns the group name of the command.
	GetGroup() string
	// GetDomainName returns the commands domain name.
	// The domain name is specified like following:
	//   sp.{group}(.{subGroup}...).{primaryInvoke}
	GetDomainName() string
	// GetSubPermissionRules returns optional sub
	// permissions of the command.
	GetSubPermissionRules() []SubPermission
	// Exec is called when the command is executed and
	// is getting passed the command CommandArgs.
	// When the command was executed successfully, it
	// should return nil. Otherwise, the error
	// encountered should be returned.
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
