import styled from 'styled-components';
import { Container } from './Container';
import { Clickable, LinearGradient } from './styleParts';

type Props = {
  color: string;
};

export const ColorTile = styled(Container)<Props>`
  ${Clickable()}

  font-family: 'Cantarell', sans-serif;
  font-size: 1.2em;
  ${(p) => LinearGradient(p.color)};
`;
