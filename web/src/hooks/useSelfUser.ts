import { useApi } from './useApi';
import { useEffect } from 'react';
import { useStore } from '../services/store';

export const useSelfUser = () => {
  const [selfUser, setSelfUser] = useStore((s) => [s.selfUser, s.setSelfUser]);
  const fetch = useApi();

  useEffect(() => {
    if (selfUser.value || selfUser.isFetching) return;
    setSelfUser({ isFetching: true });
    fetch((c) => c.etc.me())
      .then((me) => setSelfUser({ value: me, isFetching: false }))
      .catch();
  }, []);

  return selfUser.value;
};
