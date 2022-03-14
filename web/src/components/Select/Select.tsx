import React, { useEffect, useState } from 'react';
import styled from 'styled-components';

export type Props<T extends unknown> = React.HTMLAttributes<HTMLDivElement> & {
  options: Element<T>[];
  value?: Element<T>;
  onElementSelect?: (v: Element<T>) => void;
};

export type Element<T> = {
  id: string;
  display: string | JSX.Element;
  value: T;
};

const SelectContainer = styled.div`
  position: relative;
  cursor: pointer;
`;

const ValueContainer = styled.div`
  border-radius: 8px;
  background-color: ${(p) => p.theme.background};
  width: 100%;
  padding: 0.6em;
  border: solid 1px ${(p) => p.theme.accent};
`;

const OptionsList = styled.div<{ show: boolean }>`
  pointer-events: ${(p) => (p.show ? 'all' : 'none')};
  opacity: ${(p) => (p.show ? 1 : 0)};

  position: absolute;
  z-index: 100;
  width: 100%;
  border-radius: 8px;
  max-height: 20em;
  overflow-y: auto;
  transform: translateY(0.3em);
  box-shadow: 0 0.5em 3em 0 rgba(0 0 0 / 50%);
  background-color: ${(p) => p.theme.background3};
  transition: opacity 0.2s ease;
`;

const OptionContainer = styled.div`
  padding: 0.6em;
  transition: background-color 0.2s ease;
  &:hover {
    background-color: ${(p) => p.theme.accent};
  }
`;

const stopPropagation = <T extends Event>(e: T, handler: (e: T) => void) => {
  e.stopPropagation();
  e.preventDefault();
  handler(e);
};

export const Select = <T extends unknown>({
  options,
  value,
  onElementSelect = () => {},
  ...props
}: Props<T>) => {
  const [select, setSelect] = useState(false);

  const _onWindowClick = () => {
    setSelect(false);
  };

  const _onSelect = (e: Element<T>) => {
    onElementSelect(e);
  };

  useEffect(() => {
    window.addEventListener('click', _onWindowClick);
    return () => window.removeEventListener('click', _onWindowClick);
  }, []);

  return (
    <SelectContainer {...props}>
      <ValueContainer
        onClick={(e) =>
          stopPropagation(e.nativeEvent, () => setSelect(!select))
        }
      >
        {value?.display}
      </ValueContainer>
      <OptionsList show={select}>
        {options
          .filter((o) => o.id !== value?.id)
          .map((o) => (
            <OptionContainer key={o.id} onClick={() => _onSelect(o)}>
              {o.display}
            </OptionContainer>
          ))}
      </OptionsList>
    </SelectContainer>
  );
};
