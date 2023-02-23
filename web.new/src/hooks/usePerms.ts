import { useEffect, useState } from 'react';

import { useApi } from './useApi';
import { useSelfUser } from './useSelfUser';

const allowedCache: { [key: string]: Promise<string[]> } = {};

export const usePerms = (guildid?: string) => {
  const fetch = useApi();
  const selfUser = useSelfUser();
  const [allowedPerms, setAllowedPerms] = useState<string[]>();

  useEffect(() => {
    if (!!selfUser && !!guildid) {
      if (allowedCache[guildid] === undefined) {
        allowedCache[guildid] = fetch((c) =>
          c.guilds.member(guildid, selfUser.id).permissionsAllowed(),
        )
          .then((res) => {
            allowedCache[guildid] = Promise.resolve(res.data);
            return res.data;
          })
          .catch();
      }

      allowedCache[guildid].then(setAllowedPerms);
    }
  }, [selfUser, guildid]);

  const isAllowed = (pattern: string) => {
    if (pattern.includes('*')) {
      const rx = new RegExp(pattern.replaceAll('*', '.*'));
      return !!allowedPerms?.find((a) => rx.test(a));
    }

    return !!allowedPerms?.find((a) => a.includes(pattern));
  };

  return { allowedPerms, isAllowed };
};
