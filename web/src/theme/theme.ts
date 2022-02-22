export enum AppTheme {
  DARK = 0,
  LIGHT = 1,
}

export interface Theme {
  background: string;
  accent: string;
  accentLight: string;
  accentDark: string;
  text: string;
  gray: string;
  darkGray: string;
  textRed: string;

  info: string;
  success: string;
  warn: string;
  error: string;
}

export const DarkTheme: Theme = {
  background: '#1e1e1e',
  accent: '#00acd7',
  accentLight: '#0081d7',
  accentDark: '#001419',
  text: '#f4f4f5',
  gray: '#455a64',
  darkGray: '#263238',
  textRed: '#EF5350',

  info: '#039BE5',
  success: '#7CB342',
  warn: '#FB8C00',
  error: '#D81B60',
};

export const LightTheme: Theme = {
  background: '#fffffe',
  accent: '#00acd7',
  accentLight: '#0081d7',
  accentDark: '#bdd3d9',
  text: '#212121',
  gray: '#e5e5e5',
  darkGray: '#e5e5e5',
  textRed: '#EF5350',

  info: '#039BE5',
  success: '#7CB342',
  warn: '#FB8C00',
  error: '#D81B60',
};

export const DefaultTheme = DarkTheme;
