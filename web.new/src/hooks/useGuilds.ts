import { useApi } from './useApi';
import { useEffect } from 'react';
import { useStore } from '../services/store';

export const useGuilds = () => {
  const fetch = useApi();
  const [guilds, setGuilds] = useStore((s) => [s.guilds, s.setGuilds]);

  useEffect(() => {
    if (!guilds) {
      setGuilds(undefined);
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
