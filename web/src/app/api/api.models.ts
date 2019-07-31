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
}

export interface Member extends User {
  guild_id: string;
}
