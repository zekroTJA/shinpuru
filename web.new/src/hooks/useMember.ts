import { useEffect, useState } from 'react';
import { Member } from '../lib/shinpuru-ts/src';
import { GuildMemberClient } from '../lib/shinpuru-ts/src/bindings';
import { useApi } from './useApi';

type MemberRequester = <T>(
  req: (c: GuildMemberClient) => Promise<T>
) => Promise<T>;

export const useMember = (
  guildid?: string,
  memberid?: string
): [Member | undefined, MemberRequester] => {
  const fetch = useApi();
  const [member, setMember] = useState<Member>();

  const memberAction: MemberRequester = (req) =>
    fetch((c) => req(c.guilds.member(guildid!, memberid!)));

  useEffect(() => {
    if (!guildid || !memberid) return;
    memberAction((c) => c.get())
      .then((res) => setMember(res))
      .catch();
  }, [guildid, memberid]);

  return [member, memberAction];
};
