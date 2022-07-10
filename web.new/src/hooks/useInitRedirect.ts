import { useEffect } from 'react';
import { useNavigate } from 'react-router';
import LocalStorageUtil from '../util/localstorage';

const LS_KEY = 'shnp.redirect-after-login';

export const useInitRedirect = () => {
  const nav = useNavigate();

  useEffect(() => {
    const ref = LocalStorageUtil.get<string>(LS_KEY);
    if (!ref) return;
    LocalStorageUtil.del(LS_KEY);
    nav(ref);
  }, []);
};

export const setInitRedirect = (ref: string) => LocalStorageUtil.set(LS_KEY, ref);
