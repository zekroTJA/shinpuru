import Color from 'color';
import { useState } from 'react';
import { uid } from 'react-uid';
import styled, { css } from 'styled-components';
import { ReactComponent as RemoveIcon } from '../../assets/close.svg';

type Props<T> = {
  selected: TagElement<T>[];
  options: TagElement<T>[];
  onChange?: (n: TagElement<T>[]) => void;
  placeholder?: string;
};

export type TagElement<T> = {
  id: string;
  display: string | JSX.Element;
  value: T;
  keywords: string[];
};

const Selectables = styled.div`
  display: flex;
  gap: 0.5em;
  flex-wrap: wrap;
  padding: 0.4em;

  > span {
    padding: 0.2em 0.4em;
    background-color: ${(p) => p.theme.accentDarker};
    border-radius: 8px;
    cursor: pointer;
  }
`;

const Selected = styled(Selectables)`
  > span {
    background-color: ${(p) => p.theme.accent};
    cursor: default;

    > svg {
      margin-left: 0.2em;
      height: 0.8em;
      cursor: pointer;
    }
  }
`;

const TagsInputContainer = styled.div<{ focussed: boolean }>`
  width: 100%;
  border-radius: 5px;
  background-color: ${(p) => p.theme.background2};
  border: none;
  padding: 0.2em;
  transition: outline 0.2s ease;
  border: solid 2px ${(p) => new Color(p.theme.accent).fade(1).hexa()};

  ${(p) =>
    p.focussed &&
    css`
      border: solid 2px ${(p) => p.theme.accent};
    `}

  transition: all 0.25s ease;

  > div {
    display: flex;

    > input {
      width: 12ch;
      background-color: transparent;
      border: none;
      padding: 0.5em;
      outline: none;
      color: ${(p) => p.theme.text};
      font-size: 1rem;
    }
  }
`;

export const TagsInput = <T extends unknown>({
  selected,
  options,
  onChange = () => {},
  placeholder,
}: Props<T>) => {
  const [val, setVal] = useState('');
  const [focussed, setFocussed] = useState(false);

  const _onSelect = (e: TagElement<T>) => {
    setVal('');
    onChange([...selected, e]);
  };

  const _onRemove = (e: TagElement<T>) => {
    onChange(selected.filter((s) => s.id !== e.id));
  };

  const valLower = val.toLowerCase();
  const _selectables = options
    .filter((o) => !selected.find((s) => s.id === o.id))
    .filter((o) => !!o.keywords.find((kw) => kw.toLowerCase().includes(valLower)))
    .map((s) => (
      <span key={uid(s)} onClick={() => _onSelect(s)}>
        {s.display}
      </span>
    ));

  const _selected = selected.map((s) => (
    <span>
      {s.display}
      <RemoveIcon onClick={() => _onRemove(s)} />
    </span>
  ));

  return (
    <TagsInputContainer focussed={focussed}>
      <div>
        {_selected?.length > 0 && <Selected>{_selected}</Selected>}
        <input
          onFocus={() => setFocussed(true)}
          onBlur={() => setFocussed(false)}
          value={val}
          onInput={(e) => setVal(e.currentTarget.value)}
          placeholder={placeholder}
        />
      </div>
      {valLower && <Selectables>{_selectables}</Selectables>}
    </TagsInputContainer>
  );
};
