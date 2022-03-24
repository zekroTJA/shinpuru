import styled from 'styled-components';
import { LinearGradient } from './styleParts';

type Props = {
  variant?: 'default' | 'red' | 'green' | 'blue' | 'yellow' | 'orange' | 'gray';
};

export const Button = styled.button<Props>`
  font-size: 1rem;
  font-family: 'Roboto', sans-serif;
  color: ${(p) => p.theme.text};
  border: none;
  padding: 0.8em 1em;
  border-radius: 3px;
  display: flex;
  align-items: center;
  cursor: pointer;
  transition: transform 0.2s ease;

  ${(p) => {
    switch (p.variant ?? 'default') {
      case 'red':
        return LinearGradient(p.theme.red);
      case 'green':
        return LinearGradient(p.theme.green);
      case 'blue':
        return LinearGradient(p.theme.blurple);
      case 'yellow':
        return LinearGradient(p.theme.yellow);
      case 'orange':
        return LinearGradient(p.theme.orange);
      case 'gray':
        return LinearGradient(p.theme.background3);
      default:
        return LinearGradient(p.theme.accent);
    }
  }}

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
