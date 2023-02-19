import {
  APIToken,
  AccessTokenModel,
  AntiraidAction,
  AntiraidSettings,
  Channel,
  CodeExecSettings,
  CodeResponse,
  CommandInfo,
  Count,
  Guild,
  GuildBackup,
  GuildScoreboardEntry,
  GuildSettings,
  GuildSettingsApi,
  GuildStarboardEntry,
  InviteSettingsRequest,
  InviteSettingsResponse,
  JoinlogEntry,
  KarmaRule,
  KarmaSettings,
  LandingPageInfo,
  ListReponse as ListResponse,
  Member,
  MessageEmbed,
  PermissionResponse,
  PermissionsMap,
  PermissionsUpdate,
  Presence,
  PrivacyInfo,
  ReasonRequest,
  Report,
  ReportRequest,
  SearchResult,
  StarboardSortOrder,
  State,
  SystemInfo,
  UnbanRequest,
  UserSettingsOTA,
  UserSettingsPrivacy,
  VerificationSiteKey,
} from './models';
import { GuildLogEntry, GuildSettingsVerification, User } from './models';

import { Client } from './client';
import { SubClient } from './subclient';

export class EtcClient extends SubClient {
  constructor(client: Client) {
    super(client, '');
  }

  me(): Promise<User> {
    return this.req('GET', 'me');
  }

  privacyInfo(): Promise<PrivacyInfo> {
    return this.req('GET', 'privacyinfo');
  }

  sysinfo(): Promise<SystemInfo> {
    return this.req('GET', 'sysinfo');
  }
}

export class UtilClient extends SubClient {
  constructor(client: Client) {
    super(client, 'util');
  }

  color(hexCode: string, size: number = 24): string {
    return `${this.endpoint}/color/${hexCode}?size=${size}`;
  }

  commands(): Promise<ListResponse<CommandInfo>> {
    return this.req('GET', 'commands');
  }

  landingPageInfo(): Promise<LandingPageInfo> {
    return this.req('GET', 'landingpageinfo');
  }

  slashcommands(): Promise<ListResponse<CommandInfo>> {
    return this.req('GET', 'landingpageinfo');
  }
}

export class AuthClient extends SubClient {
  constructor(client: Client) {
    super(client, 'auth');
  }

  accesstoken(): Promise<AccessTokenModel> {
    return this.req('POST', 'accesstoken');
  }

  check(): Promise<CodeResponse> {
    return this.req('GET', 'check');
  }

  logout(): Promise<CodeResponse> {
    return this.req('POST', 'logout');
  }

  pushCode(code: string): Promise<CodeResponse> {
    return this.req('POST', 'pushcode', { code });
  }
}

export class SearchClient extends SubClient {
  constructor(client: Client) {
    super(client, 'search');
  }

  query(query: string, limit: number = 50): Promise<SearchResult> {
    return this.req('GET', `/?query=${query}&limit=${limit}`);
  }
}

export class TokensClient extends SubClient {
  constructor(client: Client) {
    super(client, 'token');
  }

  delete(): Promise<CodeResponse> {
    return this.req('DELETE', `/`);
  }

  info(): Promise<APIToken> {
    return this.req('GET', `/`);
  }

  generate(): Promise<APIToken> {
    return this.req('POST', `/`);
  }
}

export class GlobalSettingsClient extends SubClient {
  constructor(client: Client) {
    super(client, 'settings');
  }

  noGuildInvitesStatus(): Promise<InviteSettingsResponse> {
    return this.req('GET', 'noguildinvite');
  }

  setNoGuildInvitesStatus(ngi: InviteSettingsRequest): Promise<InviteSettingsResponse> {
    return this.req('POST', 'noguildinvite', ngi);
  }

  presence(): Promise<Presence> {
    return this.req('GET', 'presence');
  }

  setPresence(presence: Presence): Promise<Presence> {
    return this.req('POST', 'presence', presence);
  }
}

export class ReportsClient extends SubClient {
  constructor(client: Client) {
    super(client, 'reports');
  }

  get(id: string): Promise<Report> {
    return this.req('GET', id);
  }

  revoke(id: string, reason: ReasonRequest): Promise<Report> {
    return this.req('POST', `${id}/revoke`, reason);
  }
}

export class GuildMemberClient extends SubClient {
  constructor(client: Client, guildID: string, memberID: string) {
    super(client, `guilds/${guildID}/${memberID}`);
  }

  get(): Promise<Member> {
    return this.req('GET', '/');
  }

  permissions(): Promise<PermissionResponse> {
    return this.req('GET', 'permissions');
  }

  permissionsAllowed(): Promise<ListResponse<string>> {
    return this.req('GET', 'permissions/allowed');
  }

  reports(): Promise<ListResponse<Report>> {
    return this.req('GET', 'reports');
  }

  reportsCount(): Promise<Count> {
    return this.req('GET', 'reports/count');
  }

  ban(reason: ReasonRequest): Promise<Report> {
    return this.req('POST', 'ban', reason);
  }

  kick(reason: ReasonRequest): Promise<Report> {
    return this.req('POST', 'kick', reason);
  }

  mute(reason: ReasonRequest): Promise<Report> {
    return this.req('POST', 'mute', reason);
  }

  report(reason: ReportRequest): Promise<Report> {
    return this.req('POST', 'reports', reason);
  }
}

export class GuildsClient extends SubClient {
  constructor(private _client: Client) {
    super(_client, 'guilds');
  }

  list(): Promise<ListResponse<Guild>> {
    return this.req('GET', '/');
  }

  guild(id: string): Promise<Guild> {
    return this.req('GET', id);
  }

  antiraidJoinlog(id: string): Promise<ListResponse<JoinlogEntry>> {
    return this.req('GET', `${id}/antiraid/joinlog`);
  }

  deleteAntiraidJoinlog(id: string): Promise<ListResponse<JoinlogEntry>> {
    return this.req('DELETE', `${id}/antiraid/joinlog`);
  }

  setInviteBlock(id: string, enabled: boolean): Promise<ListResponse<JoinlogEntry>> {
    return this.req('POST', `${id}/inviteblock`, { enabled });
  }

  permissions(id: string): Promise<PermissionsMap> {
    return this.req('GET', `${id}/permissions`);
  }

  applyPermission(id: string, update: PermissionsUpdate): Promise<PermissionsMap> {
    return this.req('POST', `${id}/permissions`, update);
  }

  reports(id: string, limit: number = 20, offset: number = 0): Promise<ListResponse<Report>> {
    return this.req('GET', `${id}/reports?limit=${limit}&offset=${offset}`);
  }

  reportsCount(id: string): Promise<Count> {
    return this.req('GET', `${id}/reports/count`);
  }

  scoreboard(id: string, limit: number = 20): Promise<ListResponse<GuildScoreboardEntry>> {
    return this.req('GET', `${id}/scoreboard?limit=${limit}`);
  }

  starboard(
    id: string,
    sort: StarboardSortOrder = 'latest',
    limit: number = 20,
    offset: number = 0,
  ): Promise<ListResponse<GuildStarboardEntry>> {
    return this.req('GET', `${id}/starboard?limit=${limit}&offset=${offset}&sort=${sort}`);
  }

  unbanrequests(id: string): Promise<ListResponse<UnbanRequest>> {
    return this.req('GET', `${id}/unbanrequests`);
  }

  unbanrequestsCount(id: string): Promise<Count> {
    return this.req('GET', `${id}/unbanrequests/count`);
  }

  unbanrequest(id: string, requestId: string): Promise<UnbanRequest> {
    return this.req('GET', `${id}/unbanrequests/${requestId}`);
  }

  respondUnbanrequest(id: string, requestId: string, request: UnbanRequest): Promise<UnbanRequest> {
    return this.req('POST', `${id}/unbanrequests/${requestId}`, request);
  }

  settings(id: string): GuildSettingsClient {
    return new GuildSettingsClient(this._client, id);
  }

  backups(id: string): GuildBackupsClient {
    return new GuildBackupsClient(this._client, id);
  }

  members(id: string, limit = 50, after = '', filter = ''): Promise<ListResponse<Member>> {
    return this.req('GET', `${id}/members?limit=${limit}&after=${after}&filter=${filter}`);
  }

  member(id: string, memberID: string): GuildMemberClient {
    return new GuildMemberClient(this._client, id, memberID);
  }
}

export class GuildSettingsClient extends SubClient {
  constructor(client: Client, id: string) {
    super(client, `guilds/${id}/settings`);
  }

  settings(): Promise<GuildSettings> {
    return this.req('GET', '/');
  }

  setSettings(settings: GuildSettings): Promise<GuildSettings> {
    return this.req('POST', '/', settings);
  }

  antiraid(): Promise<AntiraidSettings> {
    return this.req('GET', 'antiraid');
  }

  setAntiraid(settings: AntiraidSettings): Promise<AntiraidSettings> {
    return this.req('POST', 'antiraid', settings);
  }

  addAntiraidAction(payload: AntiraidAction): Promise<CodeResponse> {
    return this.req('POST', 'antiraid/action', payload);
  }

  api(): Promise<GuildSettingsApi> {
    return this.req('GET', 'api');
  }

  setApi(settings: GuildSettingsApi): Promise<GuildSettingsApi> {
    return this.req('POST', 'api', settings);
  }

  codeexec(): Promise<CodeExecSettings> {
    return this.req('GET', 'codeexec');
  }

  setCodeexec(settings: CodeExecSettings): Promise<CodeExecSettings> {
    return this.req('POST', 'codeexec', settings);
  }

  flushData(leave_after: boolean, validation: string): Promise<CodeResponse> {
    return this.req('POST', 'flushguilddata', { leave_after, validation });
  }

  karma(): Promise<KarmaSettings> {
    return this.req('GET', 'karma');
  }

  setKarma(settings: KarmaSettings): Promise<KarmaSettings> {
    return this.req('POST', 'karma', settings);
  }

  karmaBlocklist(): Promise<ListResponse<Member>> {
    return this.req('GET', 'karma/blocklist');
  }

  addKarmaBlocklist(memberId: string): Promise<Member> {
    return this.req('PUT', `karma/blocklist/${memberId}`);
  }

  removeKarmaBlocklist(memberId: string): Promise<CodeResponse> {
    return this.req('DELETE', `karma/blocklist/${memberId}`);
  }

  karmaRules(): Promise<ListResponse<KarmaRule>> {
    return this.req('GET', 'karma/rules');
  }

  addKarmaRule(rule: KarmaRule): Promise<KarmaRule> {
    return this.req('POST', 'karma/rules', rule);
  }

  removeKarmaRule(id: string): Promise<CodeResponse> {
    return this.req('DELETE', `karma/rules/${id}`);
  }

  updateKarmaRules(rule: KarmaRule): Promise<KarmaRule> {
    return this.req('POST', `karma/rules/${rule.id}`, rule);
  }

  logs(limit = 50, offset = 0, severity = -1): Promise<ListResponse<GuildLogEntry>> {
    return this.req('GET', `logs?limit=${limit}&offset=${offset}&severity=${severity}`);
  }

  logsCount(severity = -1): Promise<Count> {
    return this.req('GET', `logs/count?severity=${severity}`);
  }

  flushLogs(): Promise<CodeResponse> {
    return this.req('DELETE', 'logs');
  }

  removeLogEntry(id: string): Promise<CodeResponse> {
    return this.req('DELETE', `logs/${id}`);
  }

  logsEnabled(): Promise<State> {
    return this.req('GET', 'logs/state');
  }

  setLogsEnabled(state: boolean): Promise<State> {
    return this.req('POST', 'logs/state', { state });
  }

  verification(): Promise<GuildSettingsVerification> {
    return this.req('GET', 'verification');
  }

  setVerification(state: GuildSettingsVerification): Promise<GuildSettingsVerification> {
    return this.req('POST', 'verification', state);
  }
}

export class GuildBackupsClient extends SubClient {
  constructor(private _client: Client, id: string) {
    super(_client, `guilds/${id}/backups`);
  }

  list(): Promise<ListResponse<GuildBackup>> {
    return this.req('GET', '/');
  }

  download(id: string): Promise<AccessTokenModel> {
    return this.req('POST', `${id}/download`);
  }

  downloadUrl(id: string, otaToken: string): Promise<string> {
    return Promise.resolve(
      `${this._client.clientEndpoint}/${this.sub}/${id}/download?ota_token=${otaToken}`,
    );
  }

  toggle(enabled: boolean): Promise<CodeResponse> {
    return this.req('POST', 'toggle', { enabled });
  }
}

export class UnbanRequestsClient extends SubClient {
  constructor(client: Client) {
    super(client, 'unbanrequests');
  }

  list(): Promise<ListResponse<UnbanRequest>> {
    return this.req('GET', '/');
  }

  create(request: UnbanRequest): Promise<UnbanRequest> {
    return this.req('POST', '/', request);
  }

  guilds(): Promise<ListResponse<Guild>> {
    return this.req('GET', 'bannedguilds');
  }
}

export class ChannelsClient extends SubClient {
  constructor(client: Client) {
    super(client, 'channels');
  }

  list(guildID: string): Promise<ListResponse<Channel>> {
    return this.req('GET', `${guildID}`);
  }

  sendEmbedMessage(guildID: string, channelID: string, embed: MessageEmbed): Promise<Channel> {
    return this.req('POST', `${guildID}/${channelID}`, embed);
  }

  updateEmbedMessage(
    guildID: string,
    channelID: string,
    messageID: string,
    embed: MessageEmbed,
  ): Promise<Channel> {
    return this.req('POST', `${guildID}/${channelID}/${messageID}`, embed);
  }
}

export class VerificationsClient extends SubClient {
  constructor(client: Client) {
    super(client, 'verification');
  }

  sitekey(): Promise<VerificationSiteKey> {
    return this.req('GET', 'sitekey');
  }

  verify(token: string): Promise<VerificationSiteKey> {
    return this.req('POST', 'verify', { token });
  }
}

export class UsersClient extends SubClient {
  constructor(client: Client) {
    super(client, 'users');
  }

  get(id: string): Promise<User> {
    return this.req('GET', id);
  }
}

export class UserSettingsClient extends SubClient {
  constructor(client: Client) {
    super(client, 'usersettings');
  }

  ota(): Promise<UserSettingsOTA> {
    return this.req('GET', 'ota');
  }

  setOta(state: UserSettingsOTA): Promise<CodeResponse> {
    return this.req('POST', 'ota', state);
  }

  privacy(): Promise<UserSettingsPrivacy> {
    return this.req('GET', 'privacy');
  }

  setPrivacy(state: UserSettingsPrivacy): Promise<CodeResponse> {
    return this.req('POST', 'privacy', state);
  }

  flush(): Promise<CodeResponse> {
    return this.req('POST', 'flush');
  }
}
