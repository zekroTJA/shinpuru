import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { Modal } from '../Modal';

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

const ValueContainer = styled.div`
  border-radius: 8px;
  background-color: ${(p) => p.theme.background};
  width: 100%;
  padding: 0.6em;
  border: solid 1px ${(p) => p.theme.accentDarker};
`;

const OptionContainer = styled.div`
  padding: 0.6em;
  transition: background-color 0.2s ease;
  &:hover {
    background-color: ${(p) => p.theme.accentDarker};
  }
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
  box-shadow: 0 0.3em 2em 0 rgba(0 0 0 / 25%);
  background-color: ${(p) => p.theme.background3};
  transition: opacity 0.2s ease;
`;

const OptionsModal = styled(Modal)`
  > div > section {
    padding: 0px;

    > ${OptionContainer} {
      padding: 1em;
    }
  }
`;

const SelectContainer = styled.div`
  position: relative;
  cursor: pointer;

  @media (orientation: portrait) {
    ${OptionsList} {
      display: none;
    }
  }

  @media (orientation: landscape) {
    ${OptionsModal} {
      display: none;
    }
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

  const _onSelectionClick: React.MouseEventHandler<HTMLDivElement> = (e) =>
    stopPropagation(e.nativeEvent, () => setSelect(!select));

  useEffect(() => {
    window.addEventListener('click', _onWindowClick);
    return () => window.removeEventListener('click', _onWindowClick);
  }, []);

  const _options = options
    .filter((o) => o.id !== value?.id)
    .map((o) => (
      <OptionContainer key={o.id} onClick={() => _onSelect(o)}>
        {o.display}
      </OptionContainer>
    ));

  return (
    <SelectContainer {...props}>
      <ValueContainer onClick={_onSelectionClick}>{value?.display}</ValueContainer>
      <OptionsList show={select}>{_options}</OptionsList>
      <OptionsModal show={select} onClose={() => setSelect(false)}>
        {_options}
      </OptionsModal>
    </SelectContainer>
  );
};
