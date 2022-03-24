import styled from 'styled-components';
import { Heading } from '../Heading';

export const ModalHeader = styled.header`
  display: flex;
  align-items: center;
  padding: 1em 1em 0 1em;

  > ${Heading} {
    margin: 0;
  }
`;
