import { DependencyList, useEffect } from 'react';

type AsyncEffectCallback = () => Promise<void>;

export const useEffectAsync = (
  effect: AsyncEffectCallback,
  deps?: DependencyList,
  destructor?: () => void,
) => {
  useEffect(() => {
    effect();
    return destructor;
  }, deps);
};
