import styled from 'styled-components';

export const SplitContainer = styled.div`
  width: 100%;
  display: flex;
  gap: 1em;

  > section {
    width: 100%;
  }

  @media (orientation: portrait) {
    flex-direction: column;
  }
`;
