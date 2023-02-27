import 'styled-components';
import { Theme } from './theme/theme';

declare module 'styled-components' {
  export interface DefaultTheme extends Theme {}
}
