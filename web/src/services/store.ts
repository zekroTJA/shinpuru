import { AppTheme } from '../theme/theme';
import LocalStorageUtil from '../util/localstorage';
import create from 'zustand';
import { User } from '../lib/shinpuru-ts/src';

export interface Store {
  theme: AppTheme;
  setTheme: (v: AppTheme) => void;

  selfUser?: User;
  setSelfUser: (selfUser: User) => void;
}

export const useStore = create<Store>((set, get) => ({
  theme: LocalStorageUtil.get(
    'shinpuru.theme',
    window.matchMedia('(prefers-color-scheme: dark)').matches
      ? AppTheme.DARK
      : AppTheme.LIGHT
  )!,
  setTheme: (theme) => {
    set({ theme });
    LocalStorageUtil.set('shinpuru.theme', theme);
  },

  selfUser: undefined,
  setSelfUser: (selfUser: User) => set({ selfUser }),
}));