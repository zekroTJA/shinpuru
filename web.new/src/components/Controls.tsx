import styled from 'styled-components';

export const Controls = styled.div`
  display: flex;
  gap: 1em;
  margin-top: 2em;

  > * {
    width: 100%;
  }

  @media (orientation: portrait) {
    flex-direction: column;
  }
`;
