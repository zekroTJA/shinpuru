import { LinearGradient } from './styleParts';
import styled from 'styled-components';

export type ButtonVariant =
  | 'default'
  | 'red'
  | 'green'
  | 'blue'
  | 'yellow'
  | 'orange'
  | 'gray'
  | 'pink';

export type ButtonProps = {
  variant?: ButtonVariant;
  nvp?: boolean;
  margin?: string;
};

export const Button = styled.button<ButtonProps>`
  font-size: 1rem;
  font-family: 'Roboto', sans-serif;
  color: ${(p) => p.theme.textAlt};
  border: none;
  padding: ${(p) => (p.nvp ? '0' : '0.8em')} 1em;
  border-radius: 3px;
  display: flex;
  gap: 0.8em;
  align-items: center;
  cursor: pointer;
  transition: transform 0.2s ease;
  justify-content: center;
  margin: ${(p) => p.margin};

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
      case 'pink':
        return LinearGradient(p.theme.pink);
      default:
        return LinearGradient(p.theme.accent);
    }
  }}

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  &:enabled:hover {
    transform: translateY(-3px);
  }
`;
