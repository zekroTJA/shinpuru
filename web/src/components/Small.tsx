import { TextAlignProps } from './props';
import styled from 'styled-components';

type Props = TextAlignProps & {};

export const Small = styled.p<Props>`
  opacity: 0.75;
  font-size: 0.9em;
  text-align: ${(p) => p.textAlign};
`;
