import styled, { css } from 'styled-components';

import Color from 'color';

type Props = {
  colors?: string | number;
  borderRadius?: string;
};

export const Tag = styled.span<Props>`
  padding: 0.2em 0.5em;
  font-size: 0.8rem;
  border-radius: ${(p) => p.borderRadius};

  ${(p) => {
    const clr = new Color(!!p.colors ? p.colors : p.theme.text);
    return css`
      color: ${clr.hex()};
      border: solid 1px ${clr.hex()};
      background-color: ${clr.fade(0.9).hexa()};
    `;
  }}
`;

Tag.defaultProps = {
  borderRadius: '3px',
};
