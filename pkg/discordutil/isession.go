package discordutil

import (
	. "github.com/bwmarrin/discordgo"
	"image"
	"io"
	"time"
)

type ISession interface {
	AddHandler(handler interface{}) func()
	AddHandlerOnce(handler interface{}) func()
	Application(appID string) (st *Application, err error)
	ApplicationAssets(appID string) (ass []*Asset, err error)
	ApplicationBotCreate(appID string) (st *User, err error)
	ApplicationCommand(appID, guildID, cmdID string, options ...RequestOption) (cmd *ApplicationCommand, err error)
	ApplicationCommandBulkOverwrite(appID string, guildID string, commands []*ApplicationCommand, options ...RequestOption) (createdCommands []*ApplicationCommand, err error)
	ApplicationCommandCreate(appID string, guildID string, cmd *ApplicationCommand, options ...RequestOption) (ccmd *ApplicationCommand, err error)
	ApplicationCommandDelete(appID, guildID, cmdID string, options ...RequestOption) error
	ApplicationCommandEdit(appID, guildID, cmdID string, cmd *ApplicationCommand, options ...RequestOption) (updated *ApplicationCommand, err error)
	ApplicationCommandPermissions(appID, guildID, cmdID string, options ...RequestOption) (permissions *GuildApplicationCommandPermissions, err error)
	ApplicationCommandPermissionsBatchEdit(appID, guildID string, permissions []*GuildApplicationCommandPermissions, options ...RequestOption) (err error)
	ApplicationCommandPermissionsEdit(appID, guildID, cmdID string, permissions *ApplicationCommandPermissionsList, options ...RequestOption) (err error)
	ApplicationCommands(appID, guildID string, options ...RequestOption) (cmd []*ApplicationCommand, err error)
	ApplicationCreate(ap *Application) (st *Application, err error)
	ApplicationDelete(appID string) (err error)
	ApplicationRoleConnectionMetadata(appID string) (st []*ApplicationRoleConnectionMetadata, err error)
	ApplicationRoleConnectionMetadataUpdate(appID string, metadata []*ApplicationRoleConnectionMetadata) (st []*ApplicationRoleConnectionMetadata, err error)
	ApplicationUpdate(appID string, ap *Application) (st *Application, err error)
	Applications() (st []*Application, err error)
	AutoModerationRule(guildID, ruleID string, options ...RequestOption) (st *AutoModerationRule, err error)
	AutoModerationRuleCreate(guildID string, rule *AutoModerationRule, options ...RequestOption) (st *AutoModerationRule, err error)
	AutoModerationRuleDelete(guildID, ruleID string, options ...RequestOption) (err error)
	AutoModerationRuleEdit(guildID, ruleID string, rule *AutoModerationRule, options ...RequestOption) (st *AutoModerationRule, err error)
	AutoModerationRules(guildID string, options ...RequestOption) (st []*AutoModerationRule, err error)
	Channel(channelID string, options ...RequestOption) (st *Channel, err error)
	ChannelDelete(channelID string, options ...RequestOption) (st *Channel, err error)
	ChannelEdit(channelID string, data *ChannelEdit, options ...RequestOption) (st *Channel, err error)
	ChannelEditComplex(channelID string, data *ChannelEdit, options ...RequestOption) (st *Channel, err error)
	ChannelFileSend(channelID, name string, r io.Reader, options ...RequestOption) (*Message, error)
	ChannelFileSendWithMessage(channelID, content string, name string, r io.Reader, options ...RequestOption) (*Message, error)
	ChannelInviteCreate(channelID string, i Invite, options ...RequestOption) (st *Invite, err error)
	ChannelInvites(channelID string, options ...RequestOption) (st []*Invite, err error)
	ChannelMessage(channelID, messageID string, options ...RequestOption) (st *Message, err error)
	ChannelMessageCrosspost(channelID, messageID string, options ...RequestOption) (st *Message, err error)
	ChannelMessageDelete(channelID, messageID string, options ...RequestOption) (err error)
	ChannelMessageEdit(channelID, messageID, content string, options ...RequestOption) (*Message, error)
	ChannelMessageEditComplex(m *MessageEdit, options ...RequestOption) (st *Message, err error)
	ChannelMessageEditEmbed(channelID, messageID string, embed *MessageEmbed, options ...RequestOption) (*Message, error)
	ChannelMessageEditEmbeds(channelID, messageID string, embeds []*MessageEmbed, options ...RequestOption) (*Message, error)
	ChannelMessagePin(channelID, messageID string, options ...RequestOption) (err error)
	ChannelMessageSend(channelID string, content string, options ...RequestOption) (*Message, error)
	ChannelMessageSendComplex(channelID string, data *MessageSend, options ...RequestOption) (st *Message, err error)
	ChannelMessageSendEmbed(channelID string, embed *MessageEmbed, options ...RequestOption) (*Message, error)
	ChannelMessageSendEmbedReply(channelID string, embed *MessageEmbed, reference *MessageReference, options ...RequestOption) (*Message, error)
	ChannelMessageSendEmbeds(channelID string, embeds []*MessageEmbed, options ...RequestOption) (*Message, error)
	ChannelMessageSendEmbedsReply(channelID string, embeds []*MessageEmbed, reference *MessageReference, options ...RequestOption) (*Message, error)
	ChannelMessageSendReply(channelID string, content string, reference *MessageReference, options ...RequestOption) (*Message, error)
	ChannelMessageSendTTS(channelID string, content string, options ...RequestOption) (*Message, error)
	ChannelMessageUnpin(channelID, messageID string, options ...RequestOption) (err error)
	ChannelMessages(channelID string, limit int, beforeID, afterID, aroundID string, options ...RequestOption) (st []*Message, err error)
	ChannelMessagesBulkDelete(channelID string, messages []string, options ...RequestOption) (err error)
	ChannelMessagesPinned(channelID string, options ...RequestOption) (st []*Message, err error)
	ChannelNewsFollow(channelID, targetID string, options ...RequestOption) (st *ChannelFollow, err error)
	ChannelPermissionDelete(channelID, targetID string, options ...RequestOption) (err error)
	ChannelPermissionSet(channelID, targetID string, targetType PermissionOverwriteType, allow, deny int64, options ...RequestOption) (err error)
	ChannelTyping(channelID string, options ...RequestOption) (err error)
	ChannelVoiceJoin(gID, cID string, mute, deaf bool) (voice *VoiceConnection, err error)
	ChannelVoiceJoinManual(gID, cID string, mute, deaf bool) (err error)
	ChannelWebhooks(channelID string, options ...RequestOption) (st []*Webhook, err error)
	Close() error
	CloseWithCode(closeCode int) (err error)
	FollowupMessageCreate(interaction *Interaction, wait bool, data *WebhookParams, options ...RequestOption) (*Message, error)
	FollowupMessageDelete(interaction *Interaction, messageID string, options ...RequestOption) error
	FollowupMessageEdit(interaction *Interaction, messageID string, data *WebhookEdit, options ...RequestOption) (*Message, error)
	ForumThreadStart(channelID, name string, archiveDuration int, content string, options ...RequestOption) (th *Channel, err error)
	ForumThreadStartComplex(channelID string, threadData *ThreadStart, messageData *MessageSend, options ...RequestOption) (th *Channel, err error)
	ForumThreadStartEmbed(channelID, name string, archiveDuration int, embed *MessageEmbed, options ...RequestOption) (th *Channel, err error)
	ForumThreadStartEmbeds(channelID, name string, archiveDuration int, embeds []*MessageEmbed, options ...RequestOption) (th *Channel, err error)
	Gateway(options ...RequestOption) (gateway string, err error)
	GatewayBot(options ...RequestOption) (st *GatewayBotResponse, err error)
	Guild(guildID string, options ...RequestOption) (st *Guild, err error)
	GuildApplicationCommandsPermissions(appID, guildID string, options ...RequestOption) (permissions []*GuildApplicationCommandPermissions, err error)
	GuildAuditLog(guildID, userID, beforeID string, actionType, limit int, options ...RequestOption) (st *GuildAuditLog, err error)
	GuildBan(guildID, userID string, options ...RequestOption) (st *GuildBan, err error)
	GuildBanCreate(guildID, userID string, days int, options ...RequestOption) (err error)
	GuildBanCreateWithReason(guildID, userID, reason string, days int, options ...RequestOption) (err error)
	GuildBanDelete(guildID, userID string, options ...RequestOption) (err error)
	GuildBans(guildID string, limit int, beforeID, afterID string, options ...RequestOption) (st []*GuildBan, err error)
	GuildChannelCreate(guildID, name string, ctype ChannelType, options ...RequestOption) (st *Channel, err error)
	GuildChannelCreateComplex(guildID string, data GuildChannelCreateData, options ...RequestOption) (st *Channel, err error)
	GuildChannels(guildID string, options ...RequestOption) (st []*Channel, err error)
	GuildChannelsReorder(guildID string, channels []*Channel, options ...RequestOption) (err error)
	GuildCreate(name string, options ...RequestOption) (st *Guild, err error)
	GuildCreateWithTemplate(templateCode, name, icon string, options ...RequestOption) (st *Guild, err error)
	GuildDelete(guildID string, options ...RequestOption) (err error)
	GuildEdit(guildID string, g *GuildParams, options ...RequestOption) (st *Guild, err error)
	GuildEmbed(guildID string, options ...RequestOption) (st *GuildEmbed, err error)
	GuildEmbedEdit(guildID string, data *GuildEmbed, options ...RequestOption) (err error)
	GuildEmoji(guildID, emojiID string, options ...RequestOption) (emoji *Emoji, err error)
	GuildEmojiCreate(guildID string, data *EmojiParams, options ...RequestOption) (emoji *Emoji, err error)
	GuildEmojiDelete(guildID, emojiID string, options ...RequestOption) (err error)
	GuildEmojiEdit(guildID, emojiID string, data *EmojiParams, options ...RequestOption) (emoji *Emoji, err error)
	GuildEmojis(guildID string, options ...RequestOption) (emoji []*Emoji, err error)
	GuildIcon(guildID string, options ...RequestOption) (img image.Image, err error)
	GuildIntegrationCreate(guildID, integrationType, integrationID string, options ...RequestOption) (err error)
	GuildIntegrationDelete(guildID, integrationID string, options ...RequestOption) (err error)
	GuildIntegrationEdit(guildID, integrationID string, expireBehavior, expireGracePeriod int, enableEmoticons bool, options ...RequestOption) (err error)
	GuildIntegrations(guildID string, options ...RequestOption) (st []*Integration, err error)
	GuildInvites(guildID string, options ...RequestOption) (st []*Invite, err error)
	GuildLeave(guildID string, options ...RequestOption) (err error)
	GuildMember(guildID, userID string, options ...RequestOption) (st *Member, err error)
	GuildMemberAdd(guildID, userID string, data *GuildMemberAddParams, options ...RequestOption) (err error)
	GuildMemberDeafen(guildID string, userID string, deaf bool, options ...RequestOption) (err error)
	GuildMemberDelete(guildID, userID string, options ...RequestOption) (err error)
	GuildMemberDeleteWithReason(guildID, userID, reason string, options ...RequestOption) (err error)
	GuildMemberEdit(guildID, userID string, data *GuildMemberParams, options ...RequestOption) (st *Member, err error)
	GuildMemberEditComplex(guildID, userID string, data *GuildMemberParams, options ...RequestOption) (st *Member, err error)
	GuildMemberMove(guildID string, userID string, channelID *string, options ...RequestOption) (err error)
	GuildMemberMute(guildID string, userID string, mute bool, options ...RequestOption) (err error)
	GuildMemberNickname(guildID, userID, nickname string, options ...RequestOption) (err error)
	GuildMemberRoleAdd(guildID, userID, roleID string, options ...RequestOption) (err error)
	GuildMemberRoleRemove(guildID, userID, roleID string, options ...RequestOption) (err error)
	GuildMemberTimeout(guildID string, userID string, until *time.Time, options ...RequestOption) (err error)
	GuildMembers(guildID string, after string, limit int, options ...RequestOption) (st []*Member, err error)
	GuildMembersSearch(guildID, query string, limit int, options ...RequestOption) (st []*Member, err error)
	GuildOnboarding(guildID string, options ...RequestOption) (onboarding *GuildOnboarding, err error)
	GuildOnboardingEdit(guildID string, o *GuildOnboarding, options ...RequestOption) (onboarding *GuildOnboarding, err error)
	GuildPreview(guildID string, options ...RequestOption) (st *GuildPreview, err error)
	GuildPrune(guildID string, days uint32, options ...RequestOption) (count uint32, err error)
	GuildPruneCount(guildID string, days uint32, options ...RequestOption) (count uint32, err error)
	GuildRoleCreate(guildID string, data *RoleParams, options ...RequestOption) (st *Role, err error)
	GuildRoleDelete(guildID, roleID string, options ...RequestOption) (err error)
	GuildRoleEdit(guildID, roleID string, data *RoleParams, options ...RequestOption) (st *Role, err error)
	GuildRoleReorder(guildID string, roles []*Role, options ...RequestOption) (st []*Role, err error)
	GuildRoles(guildID string, options ...RequestOption) (st []*Role, err error)
	GuildScheduledEvent(guildID, eventID string, userCount bool, options ...RequestOption) (st *GuildScheduledEvent, err error)
	GuildScheduledEventCreate(guildID string, event *GuildScheduledEventParams, options ...RequestOption) (st *GuildScheduledEvent, err error)
	GuildScheduledEventDelete(guildID, eventID string, options ...RequestOption) (err error)
	GuildScheduledEventEdit(guildID, eventID string, event *GuildScheduledEventParams, options ...RequestOption) (st *GuildScheduledEvent, err error)
	GuildScheduledEventUsers(guildID, eventID string, limit int, withMember bool, beforeID, afterID string, options ...RequestOption) (st []*GuildScheduledEventUser, err error)
	GuildScheduledEvents(guildID string, userCount bool, options ...RequestOption) (st []*GuildScheduledEvent, err error)
	GuildSplash(guildID string, options ...RequestOption) (img image.Image, err error)
	GuildTemplate(templateCode string, options ...RequestOption) (st *GuildTemplate, err error)
	GuildTemplateCreate(guildID string, data *GuildTemplateParams, options ...RequestOption) (st *GuildTemplate)
	GuildTemplateDelete(guildID, templateCode string, options ...RequestOption) (err error)
	GuildTemplateEdit(guildID, templateCode string, data *GuildTemplateParams, options ...RequestOption) (st *GuildTemplate, err error)
	GuildTemplateSync(guildID, templateCode string, options ...RequestOption) (err error)
	GuildTemplates(guildID string, options ...RequestOption) (st []*GuildTemplate, err error)
	GuildThreadsActive(guildID string, options ...RequestOption) (threads *ThreadsList, err error)
	GuildWebhooks(guildID string, options ...RequestOption) (st []*Webhook, err error)
	GuildWithCounts(guildID string, options ...RequestOption) (st *Guild, err error)
	HeartbeatLatency() time.Duration
	InteractionRespond(interaction *Interaction, resp *InteractionResponse, options ...RequestOption) error
	InteractionResponse(interaction *Interaction, options ...RequestOption) (*Message, error)
	InteractionResponseDelete(interaction *Interaction, options ...RequestOption) error
	InteractionResponseEdit(interaction *Interaction, newresp *WebhookEdit, options ...RequestOption) (*Message, error)
	Invite(inviteID string, options ...RequestOption) (st *Invite, err error)
	InviteAccept(inviteID string, options ...RequestOption) (st *Invite, err error)
	InviteComplex(inviteID, guildScheduledEventID string, withCounts, withExpiration bool, options ...RequestOption) (st *Invite, err error)
	InviteDelete(inviteID string, options ...RequestOption) (st *Invite, err error)
	InviteWithCounts(inviteID string, options ...RequestOption) (st *Invite, err error)
	MessageReactionAdd(channelID, messageID, emojiID string, options ...RequestOption) error
	MessageReactionRemove(channelID, messageID, emojiID, userID string, options ...RequestOption) error
	MessageReactions(channelID, messageID, emojiID string, limit int, beforeID, afterID string, options ...RequestOption) (st []*User, err error)
	MessageReactionsRemoveAll(channelID, messageID string, options ...RequestOption) error
	MessageReactionsRemoveEmoji(channelID, messageID, emojiID string, options ...RequestOption) error
	MessageThreadStart(channelID, messageID string, name string, archiveDuration int, options ...RequestOption) (ch *Channel, err error)
	MessageThreadStartComplex(channelID, messageID string, data *ThreadStart, options ...RequestOption) (ch *Channel, err error)
	Open() error
	Request(method, urlStr string, data interface{}, options ...RequestOption) (response []byte, err error)
	RequestGuildMembers(guildID, query string, limit int, nonce string, presences bool) error
	RequestGuildMembersBatch(guildIDs []string, query string, limit int, nonce string, presences bool) (err error)
	RequestGuildMembersBatchList(guildIDs []string, userIDs []string, limit int, nonce string, presences bool) (err error)
	RequestGuildMembersList(guildID string, userIDs []string, limit int, nonce string, presences bool) error
	RequestWithBucketID(method, urlStr string, data interface{}, bucketID string, options ...RequestOption) (response []byte, err error)
	RequestWithLockedBucket(method, urlStr, contentType string, b []byte, bucket *Bucket, sequence int, options ...RequestOption) (response []byte, err error)
	StageInstance(channelID string, options ...RequestOption) (si *StageInstance, err error)
	StageInstanceCreate(data *StageInstanceParams, options ...RequestOption) (si *StageInstance, err error)
	StageInstanceDelete(channelID string, options ...RequestOption) (err error)
	StageInstanceEdit(channelID string, data *StageInstanceParams, options ...RequestOption) (si *StageInstance, err error)
	ThreadJoin(id string, options ...RequestOption) error
	ThreadLeave(id string, options ...RequestOption) error
	ThreadMember(threadID, memberID string, withMember bool, options ...RequestOption) (member *ThreadMember, err error)
	ThreadMemberAdd(threadID, memberID string, options ...RequestOption) error
	ThreadMemberRemove(threadID, memberID string, options ...RequestOption) error
	ThreadMembers(threadID string, limit int, withMember bool, afterID string, options ...RequestOption) (members []*ThreadMember, err error)
	ThreadStart(channelID, name string, typ ChannelType, archiveDuration int, options ...RequestOption) (ch *Channel, err error)
	ThreadStartComplex(channelID string, data *ThreadStart, options ...RequestOption) (ch *Channel, err error)
	ThreadsActive(channelID string, options ...RequestOption) (threads *ThreadsList, err error)
	ThreadsArchived(channelID string, before *time.Time, limit int, options ...RequestOption) (threads *ThreadsList, err error)
	ThreadsPrivateArchived(channelID string, before *time.Time, limit int, options ...RequestOption) (threads *ThreadsList, err error)
	ThreadsPrivateJoinedArchived(channelID string, before *time.Time, limit int, options ...RequestOption) (threads *ThreadsList, err error)
	UpdateCustomStatus(state string) (err error)
	UpdateGameStatus(idle int, name string) (err error)
	UpdateListeningStatus(name string) (err error)
	UpdateStatusComplex(usd UpdateStatusData) (err error)
	UpdateStreamingStatus(idle int, name string, url string) (err error)
	UpdateWatchStatus(idle int, name string) (err error)
	User(userID string, options ...RequestOption) (st *User, err error)
	UserApplicationRoleConnection(appID string) (st *ApplicationRoleConnection, err error)
	UserApplicationRoleConnectionUpdate(appID string, rconn *ApplicationRoleConnection) (st *ApplicationRoleConnection, err error)
	UserAvatar(userID string, options ...RequestOption) (img image.Image, err error)
	UserAvatarDecode(u *User, options ...RequestOption) (img image.Image, err error)
	UserChannelCreate(recipientID string, options ...RequestOption) (st *Channel, err error)
	UserChannelPermissions(userID, channelID string, fetchOptions ...RequestOption) (apermissions int64, err error)
	UserConnections(options ...RequestOption) (conn []*UserConnection, err error)
	UserGuildMember(guildID string, options ...RequestOption) (st *Member, err error)
	UserGuilds(limit int, beforeID, afterID string, options ...RequestOption) (st []*UserGuild, err error)
	UserUpdate(username, avatar string, options ...RequestOption) (st *User, err error)
	VoiceRegions(options ...RequestOption) (st []*VoiceRegion, err error)
	Webhook(webhookID string, options ...RequestOption) (st *Webhook, err error)
	WebhookCreate(channelID, name, avatar string, options ...RequestOption) (st *Webhook, err error)
	WebhookDelete(webhookID string, options ...RequestOption) (err error)
	WebhookDeleteWithToken(webhookID, token string, options ...RequestOption) (st *Webhook, err error)
	WebhookEdit(webhookID, name, avatar, channelID string, options ...RequestOption) (st *Role, err error)
	WebhookEditWithToken(webhookID, token, name, avatar string, options ...RequestOption) (st *Role, err error)
	WebhookExecute(webhookID, token string, wait bool, data *WebhookParams, options ...RequestOption) (st *Message, err error)
	WebhookMessage(webhookID, token, messageID string, options ...RequestOption) (message *Message, err error)
	WebhookMessageDelete(webhookID, token, messageID string, options ...RequestOption) (err error)
	WebhookMessageEdit(webhookID, token, messageID string, data *WebhookEdit, options ...RequestOption) (st *Message, err error)
	WebhookThreadExecute(webhookID, token string, wait bool, threadID string, data *WebhookParams, options ...RequestOption) (st *Message, err error)
	WebhookWithToken(webhookID, token string, options ...RequestOption) (st *Webhook, err error)
}
