import create from 'zustand';
import { Notification } from '../components/Notifications';
import { ModalState } from '../hooks/useModal';
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

  modal: ModalState<any>;
  setModal: (modal: ModalState<any>) => void;
}

export const useStore = create<Store>((set, get) => ({
  theme: LocalStorageUtil.get(
    'shnp.theme',
    window.matchMedia('(prefers-color-scheme: dark)').matches ? AppTheme.DARK : AppTheme.LIGHT,
  )!,
  setTheme: (theme) => {
    set({ theme });
    LocalStorageUtil.set('shnp.theme', theme);
  },

  selfUser: undefined,
  setSelfUser: (selfUser: User) => set({ selfUser }),

  guilds: undefined,
  setGuilds: (guilds) => set({ guilds }),

  selectedGuild: undefined,
  setSelectedGuild: (selectedGuild) => {
    set({ selectedGuild });
    if (!!selectedGuild) LocalStorageUtil.set('shnp.selectedguild', selectedGuild.id);
  },

  notifications: [],
  setNotifications: (notifications) => set({ notifications }),

  modal: { isOpen: false },
  setModal: (modal) => set({ modal }),
}));
