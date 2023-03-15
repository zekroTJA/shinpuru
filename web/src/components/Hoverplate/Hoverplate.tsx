import styled, { css } from 'styled-components';

import { useRef } from 'react';

type Direction = 'top' | 'bottom';

type Props = React.HTMLAttributes<HTMLDivElement> & {
  hoverContent: JSX.Element;
  direction?: Direction;
};

const PlateContainer = styled.div<{ height?: number; direction: Direction }>`
  position: absolute;
  left: -1em;
  white-space: nowrap;
  background-color: ${(p) => p.theme.background3};
  padding: 1em 1em 1em 1em;
  z-index: -1;
  border-radius: 12px;
  box-shadow: 0 0.5em 2em 0 rgba(0 0 0 / 0.4);

  opacity: 0;
  pointer-events: none;

  transition: all 0.2s ease;

  ${(p) => {
    switch (p.direction) {
      case 'top':
        return css`
          bottom: -1em;
          padding-bottom: calc(2em + ${p.height}px);
        `;
      case 'bottom':
        return css`
          top: -1em;
          padding-top: calc(2em + ${p.height}px);
        `;
    }
  }}
`;

const ContentContainer = styled.div``;

const HoverContainer = styled.div`
  z-index: 1;
  position: relative;
  cursor: default;

  &:hover > ${PlateContainer} {
    opacity: 1;
    pointer-events: all;
  }
`;

export const Hoverplate: React.FC<Props> = ({
  children,
  hoverContent,
  direction = 'top',
  ...props
}) => {
  const contentRef = useRef<HTMLDivElement>(null);
  return (
    <HoverContainer {...props}>
      <ContentContainer ref={contentRef}>{children}</ContentContainer>
      <PlateContainer height={contentRef.current?.offsetHeight} direction={direction}>
        {hoverContent}
      </PlateContainer>
    </HoverContainer>
  );
};
