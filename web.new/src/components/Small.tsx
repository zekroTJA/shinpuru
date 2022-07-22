import styled from 'styled-components';
import { TextAlignProps } from './props';

type Props = TextAlignProps & {};

export const Small = styled.p<Props>`
  opacity: 0.75;
  font-size: 0.9em;
  text-align: ${(p) => p.textAlign};
`;
