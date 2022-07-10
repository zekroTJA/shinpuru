import styled from 'styled-components';
import { Container } from '../Container';

export const ModalControls = styled(Container)`
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 1em;
  background-color: ${(p) => p.theme.background3};
`;
