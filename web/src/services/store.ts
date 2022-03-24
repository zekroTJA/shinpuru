import create from 'zustand';
import { Notification } from '../components/Notifications';
import { Guild, User } from '../lib/shinpuru-ts/src';
import { AppTheme } from '../theme/theme';
import LocalStorageUtil from '../util/localstorage';

export interface Store {
  theme: AppTheme;
  setTheme: (v: AppTheme) => void;

  selfUser?: User;
  setSelfUser: (selfUser: User) => void;

  guilds?: Guild[];
  setGuilds: (guilds: Guild[]) => void;

  selectedGuild?: Guild;
  setSelectedGuild: (selectedGuild: Guild) => void;

  notifications: Notification[];
  setNotifications: (notifications: Notification[]) => void;
}

export const useStore = create<Store>((set, get) => ({
  theme: LocalStorageUtil.get(
    'shinpuru.theme',
    window.matchMedia('(prefers-color-scheme: dark)').matches ? AppTheme.DARK : AppTheme.LIGHT,
  )!,
  setTheme: (theme) => {
    set({ theme });
    LocalStorageUtil.set('shinpuru.theme', theme);
  },

  selfUser: undefined,
  setSelfUser: (selfUser: User) => set({ selfUser }),

  guilds: undefined,
  setGuilds: (guilds) => set({ guilds }),

  selectedGuild: undefined,
  setSelectedGuild: (selectedGuild) => set({ selectedGuild }),

  notifications: [],
  setNotifications: (notifications) => set({ notifications }),
}));
