import { Clickable, LinearGradient } from './styleParts';

import { Container } from './Container';
import styled from 'styled-components';

type Props = {
  color: string;
};

export const ColorTile = styled(Container)<Props>`
  ${Clickable()}

  font-family: 'Cantarell', sans-serif;
  font-size: 1.2em;
  ${(p) => LinearGradient(p.color)};
`;
