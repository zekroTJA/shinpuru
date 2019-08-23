/** @format */

export interface ListReponse<T> {
  n: number;
  data: T[];
}

export interface User {
  id: string;
  username: string;
  avatar: string;
  locale: string;
  discriminator: string;
  verified: boolean;
  bot: boolean;
  avatar_url: string;
  created_at?: string;
  bot_owner?: boolean;
}

export interface Role {
  id: string;
  name: string;
  managed: boolean;
  mentionable: boolean;
  hoist: boolean;
  color: number;
  position: number;
  permission: number;
}

export interface Member {
  user: User;
  guild_id: string;
  joined_at: string;
  nick: string;
  avatar_url?: string;
  roles: string[];
  created_at?: string;
  dominance?: number;
}

export interface Channel {
  id: string;
  guild_id: string;
  name: string;
  topic: string;
  type: number;
  nsfw: boolean;
  position: number;
  user_limit: number;
  parent_id: string;
}

export interface Guild {
  id: string;
  name: string;
  icon: string;
  icon_url: string;
  region: string;
  owner_id: string;
  joined_at: string;
  member_count: number;

  self_member?: Member;

  roles?: Role[];
  members?: Member[];
  channels?: Channel[];
}

export interface PermissionResponse {
  permissions: number;
}

export interface Report {
  id: string;
  type: number;
  type_name?: string;
  guild_id: string;
  executor_id: string;
  victim_id: string;
  message: string;
  attachment_url: string;
  created?: string;
}

export interface GuildSettings {
  prefix: string;
  perms: Map<string, string[]>;
  autorole: string;
  modlogchannel: string;
  voicelogchannel: string;
  joinmessagechannel: string;
  joinmessagetext: string;
  leavemessagechannel: string;
  leavemessagetext: string;
}

export interface PermissionsUpdate {
  perm: string;
  role_ids: string[];
}

export interface ReasonRequest {
  reason: string;
  attachment: string;
}

export interface ReportRequest extends ReasonRequest {
  type: number;
}

export interface Presence {
  game: string;
  status: string;
}

export interface InviteSettingsRequest {
  guild_id: string;
  message: string;
  invite_code?: string;
}

export interface InviteSettingsResponse {
  guild: Guild;
  invite_url: string;
  message: string;
}

export interface Count {
  count: number;
}

export interface SystemInfo {
  version: string;
  commit_hash: string;
  build_date: Date;
  go_version: string;
  uptime: number;
  uptime_str: string;
  os: string;
  arch: string;
  cpus: number;
  go_routines: number;
  stack_use_str: string;
  heap_use_str: string;
  bot_user_id: string;
  bot_invite: string;
  guilds: number;
}
