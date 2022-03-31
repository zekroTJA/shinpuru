import styled from 'styled-components';

type Props = {
  wrap?: boolean;
  gap?: string;
  direction?: 'row' | 'column';
};

export const Flex = styled.div<Props>`
  display: flex;
  flex-wrap: ${(p) => (p.wrap ? 'wrap' : 'nowrap')};
  gap: ${(p) => p.gap};
  flex-direction: ${(p) => p.direction};
`;

Flex.defaultProps = {
  direction: 'row',
};
