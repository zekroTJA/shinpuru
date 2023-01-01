import Color from 'color';
import styled from 'styled-components';
import { LinearGradient } from './styleParts';

type Props = {
  color?: string;
};

export const Card = styled.div<Props>`
  padding: 1em;
  border: solid 1px ${(p) => p.color ?? p.theme.blurple};
  border-radius: 12px;
  width: fit-content;

  ${(p) =>
    LinearGradient(
      Color(p.color ?? p.theme.blurple)
        .alpha(0.1)
        .hexa(),
    )};
`;
