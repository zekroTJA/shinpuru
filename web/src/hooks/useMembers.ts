import { useEffect, useState } from 'react';
import { Member } from '../lib/shinpuru-ts/src';
import { useApi } from './useApi';

export const useMembers = (guildid?: string) => {
  const fetch = useApi();
  const [members, setMembers] = useState<Member[]>();

  useEffect(() => {
    if (!guildid) return;
    fetch((c) => c.guilds.members(guildid))
      .then((res) => setMembers(res.data))
      .catch();
  }, [guildid]);

  return members;
};
