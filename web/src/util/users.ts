import { Member } from '../lib/shinpuru-ts/src';

export function memberName(member: Member): string {
  return !!member.nick ? member.nick : member.user.username;
}
