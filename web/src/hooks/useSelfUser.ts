import { useEffect } from 'react';
import { useStore } from '../services/store';
import { useApi } from './useApi';

export const useSelfUser = () => {
  const [selfUser, setSelfUser] = useStore((s) => [s.selfUser, s.setSelfUser]);
  const fetch = useApi();

  useEffect(() => {
    if (selfUser) return;
    fetch((c) => c.etc.me())
      .then((me) => setSelfUser(me))
      .catch();
  }, []);

  return selfUser;
};
