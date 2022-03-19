import { useEffect, useRef } from 'react';
import { useStore } from '../services/store';
import { useApi } from './useApi';

export const useGuilds = () => {
  const fetch = useApi();
  const [guilds, setGuilds] = useStore((s) => [s.guilds, s.setGuilds]);

  useEffect(() => {
    if (!guilds) {
      setGuilds([]);
      fetch((c) => c.guilds.list())
        .then((res) => {
          setGuilds(res.data);
          console.log(res);
        })
        .catch();
    }
  }, []);

  return guilds;
};
