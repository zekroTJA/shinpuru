import { useEffect, useState } from 'react';
import { Member } from '../lib/shinpuru-ts/src';
import { useApi } from './useApi';
import { useSelfUser } from './useSelfUser';

export const useSelfMember = (guildid?: string) => {
  const [member, setMember] = useState<Member>();
  const selfUser = useSelfUser();
  const fetch = useApi();

  useEffect(() => {
    if (!guildid || !selfUser) return;
    setMember(undefined);
    fetch((c) => c.guilds.member(guildid, selfUser.id).get())
      .then((res) => setMember(res))
      .catch();
  }, [guildid, selfUser]);

  return member;
};
