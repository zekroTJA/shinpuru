import styled from 'styled-components';

export const Container = styled.div`
  border-radius: 12px;
  background-color: ${(p) => p.theme.background2};
  padding: 1rem;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    transform: scale(1.03);
  }
`;
