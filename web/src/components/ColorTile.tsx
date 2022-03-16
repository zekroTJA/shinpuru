import styled from 'styled-components';
import { Container } from './Container';
import Color from 'color';

type Props = {
  color: string;
};

export const ColorTile = styled(Container)<Props>`
  font-family: 'Cantarell', sans-serif;
  font-size: 1.2em;
  background: linear-gradient(
    140deg,
    ${(p) => p.color},
    ${(p) => Color(p.color).darken(0.15).hex()}
  );
`;
