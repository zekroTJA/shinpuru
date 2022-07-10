import { useEffect, useRef, useState } from 'react';
import { Member } from '../lib/shinpuru-ts/src';
import { useApi } from './useApi';

export const useMembers = (
  guildid?: string,
  limit = 100,
  filter = ''
): [Member[] | undefined, () => Promise<any>] => {
  const fetch = useApi();
  const [members, setMembers] = useState<Member[]>();
  const afterRef = useRef('');

  const _load = (reset = false) => {
    if (!guildid) return Promise.resolve();
    if (reset) afterRef.current = '';
    return fetch((c) =>
      c.guilds.members(guildid, limit, afterRef.current, filter)
    )
      .then((res) => {
        setMembers([...(!members || reset ? [] : members), ...res.data]);
        if (res.data.length !== 0)
          afterRef.current = res.data[res.data.length - 1].user.id;
      })
      .catch();
  };

  useEffect(() => {
    setMembers(undefined);
  }, [guildid]);

  useEffect(() => {
    _load(true);
  }, [guildid, filter]);

  return [members, _load];
};
