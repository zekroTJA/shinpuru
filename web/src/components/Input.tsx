import Color from 'color';
import styled from 'styled-components';

export const Input = styled.input`
  border-radius: 3px;
  background-color: ${(p) => p.theme.background2};
  border: none;
  font-size: 1rem;
  color: ${(p) => p.theme.text};
  padding: 0.5em;
  transition: outline 0.2s ease;
  outline: solid 2px ${(p) => new Color(p.theme.accent).fade(1).hexa()};

  &:enabled:focus {
    outline: solid 2px ${(p) => p.theme.accent};
  }

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
`;
