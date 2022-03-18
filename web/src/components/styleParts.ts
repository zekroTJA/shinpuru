import { css } from 'styled-components';

export const Clickable = (scaling = 1.03) => css`
  cursor: pointer;
  transition: transform 0.2s ease;

  &:hover {
    transform: scale(${scaling});
  }
`;
