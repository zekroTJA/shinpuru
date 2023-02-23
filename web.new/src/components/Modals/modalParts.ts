import { TextArea } from '../TextArea';
import styled from 'styled-components';

export const ModalTextArea = styled(TextArea)`
  background-color: ${(p) => p.theme.background3};
  min-width: 100%;
  max-width: 100%;
`;

export const ModalContainer = styled.div`
  width: 30em;
  max-width: 80vw;
  display: flex;
  flex-direction: column;
  gap: 1.5em;
`;
