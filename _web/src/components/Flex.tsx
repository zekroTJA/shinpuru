import styled, { css } from 'styled-components';

type Props = {
  wrap?: boolean;
};

export const Flex = styled.div<Props>`
  display: flex;
  align-items: center;
  flex-wrap: ${(p) => (p.wrap ? 'wrap' : 'nowrap')};
`;
