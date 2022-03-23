import Color from 'color';
import styled, { keyframes } from 'styled-components';

interface Props {
  width?: string;
  height?: string;
  borderRadius?: string;
  margin?: string;
}

const LoaderKF = keyframes`
  from {
    transform: translateX(-80%);
  }
  to {
    transform: translateX(80%);
  }
`;

export const Loader = styled.div<Props>`
  min-width: ${(p) => p.width};
  min-height: ${(p) => p.height};
  border-radius: ${(p) => p.borderRadius};
  margin: ${(p) => p.margin};
  position: relative;
  overflow: hidden;
  background-color: ${(p) => new Color(p.theme.text).fade(0.9).hexa()};

  &::after {
    content: '';
    position: absolute;
    height: 100%;
    width: 100%;
    background: linear-gradient(
      140deg,
      ${(p) => new Color(p.theme.text).fade(1).hexa()} 20%,
      ${(p) => new Color(p.theme.text).fade(0.9).hexa()} 50%,
      ${(p) => new Color(p.theme.text).fade(1).hexa()} 80%
    );
    animation: ${LoaderKF} 3s infinite;
  }
`;

Loader.defaultProps = {
  width: '100%',
  height: '3em',
  borderRadius: '12px',
  margin: '0',
};
