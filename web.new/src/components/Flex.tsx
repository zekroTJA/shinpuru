import styled from 'styled-components';

type Props = {
  wrap?: boolean;
  gap?: string;
};

export const Flex = styled.div<Props>`
  display: flex;
  align-items: center;
  flex-wrap: ${(p) => (p.wrap ? 'wrap' : 'nowrap')};
  gap: ${(p) => p.gap};
`;
