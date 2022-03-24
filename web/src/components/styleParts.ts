import Color from 'color';
import { css } from 'styled-components';

export const Clickable = (scaling = 1.03) => css`
  cursor: pointer;
  transition: transform 0.2s ease;

  &:hover {
    transform: scale(${scaling});
  }
`;

export const LinearGradient = (c1: string, c2?: string) => {
  const _c2 = c2 ?? new Color(c1).darken(0.15).hex();
  return css`
    background: linear-gradient(140deg, ${c1}, ${_c2});
  `;
};
