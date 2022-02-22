import { useStore } from 'services/store';
import { AppTheme, DarkTheme, LightTheme, Theme } from 'theme/theme';

export function useStoredTheme() {
  const appTheme = useStore((s) => s.theme);

  let theme: Theme;
  let editorTheme: string;

  switch (appTheme) {
    case AppTheme.LIGHT:
      theme = LightTheme;
      editorTheme = 'light';
      break;
    case AppTheme.DARK:
    default:
      theme = DarkTheme;
      editorTheme = 'vs-dark';
  }

  return { theme, editorTheme };
}
