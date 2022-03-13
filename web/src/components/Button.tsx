import styled from 'styled-components';

export const Button = styled.button`
  font-size: 1rem;
  font-family: 'Roboto', sans-serif;
  color: ${(p) => p.theme.text};
  background: ${(p) => p.theme.blurple};
  border: none;
  padding: 0.8em 1em;
  border-radius: 3px;
  display: flex;
  align-items: center;
  cursor: pointer;
  transition: transform 0.2s ease;

  > svg {
    margin-right: 0.8em;
  }

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  &:enabled:hover {
    transform: translateY(-3px);
  }
`;
