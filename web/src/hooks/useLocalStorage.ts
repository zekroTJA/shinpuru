import { useEffect, useState } from 'react';

import LocalStorageUtil from '../util/localstorage';

type GetterSetter<T> = [T | undefined, (v: T) => void];

export function useLocalStorage<T>(key: string, def?: T): GetterSetter<T> {
  const [get, set] = useState<T | undefined>(LocalStorageUtil.get(key, def));

  useEffect(() => {
    LocalStorageUtil.set(key, get);
  }, [get]);

  return [get, set];
}
