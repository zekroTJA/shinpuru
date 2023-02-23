import { AppTheme, DarkTheme, LightTheme, Theme } from '../theme/theme';

import Color from 'color';
import { useStore } from '../services/store';

export function useStoredTheme() {
  const [appTheme, accentColor] = useStore((s) => [s.theme, s.accentColor]);

  let theme: Theme;
  let editorTheme: string;

  switch (appTheme) {
    case AppTheme.LIGHT:
      theme = {
        ...LightTheme,
        accent: accentColor ?? LightTheme.accent,
        accentDarker: accentColor
          ? Color(accentColor).mix(Color('#ffffffa0')).hexa()
          : LightTheme.accentDarker,
      };
      editorTheme = 'light';
      break;
    case AppTheme.DARK:
    default:
      theme = {
        ...DarkTheme,
        accent: accentColor ?? DarkTheme.accent,
        accentDarker: accentColor
          ? Color(accentColor).mix(Color('#000000a0')).hexa()
          : DarkTheme.accentDarker,
      };
      editorTheme = 'vs-dark';
  }

  return { theme, editorTheme };
}
