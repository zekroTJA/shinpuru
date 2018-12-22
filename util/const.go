package util

const (
	InvitePermission = 0x1 | // CREATE INSTANT INVITE
		0x10 | // MANAGE CHANNELS
		0x20 | // MANAGE GUILD
		0x40 | // ADD REACTIONS
		0x400 | // VIEW CHANNEL
		0x800 | // SEND MESSAGES
		0x2000 | // MANAGE MESSAGES
		0x4000 | // EMBED LINKS
		0x8000 | // ATTACH FILES
		0x10000 | // READ MESSAGE HISTORY
		0x20000 | // MENTION @everyone
		0x40000 | // USE EXTERNAL EMOJIS
		0x4000000 | // CHANGE NICKNAMES
		0x8000000 | // MANAGE NICKNAMES
		0x10000000 | // MANAGE ROLES
		0x20000000 | // MANAGE WEBHOOKS
		0x40000000 // MANAGE EMOJIS

	ColorEmbedError   = 0xd32f2f
	ColorEmbedDefault = 0xffc107
	ColorEmbedUpdated = 0x8bc34a

	AutoNick = "シンプル"
	StdMotd  = "closed beta version"
)
