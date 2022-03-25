import { useEffect, useState } from 'react';
import { useApi } from './useApi';
import { useSelfUser } from './useSelfUser';

export const usePerms = (guildid?: string) => {
  const fetch = useApi();
  const selfUser = useSelfUser();
  const [allowedPerms, setAllowedPerms] = useState<string[]>();

  useEffect(() => {
    if (!!selfUser && !!guildid) {
      fetch((c) => c.guilds.member(guildid, selfUser.id).permissionsAllowed())
        .then((res) => setAllowedPerms(res.data))
        .catch();
    }
  }, [selfUser, guildid]);

  const isAllowed = (pattern: string) =>
    !!allowedPerms?.find((a) => a.includes(pattern));

  return { allowedPerms, isAllowed };
};
