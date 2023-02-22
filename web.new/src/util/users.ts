import { Member } from '../lib/shinpuru-ts/src';

export function memberName(member: Member, withDiscriminator = false): string {
  return !!member.nick
    ? member.nick
    : member.user.username + (withDiscriminator ? `#${member.user.discriminator}` : '');
}
