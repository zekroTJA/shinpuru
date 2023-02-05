import Color from 'color';
import Fuse from 'fuse.js';
import { useEffect, useReducer, useRef } from 'react';
import { uid } from 'react-uid';
import styled, { css } from 'styled-components';
import { Input } from '../Input';

type Props = {
  value: string;
  setValue: (v: string) => void;
  selections: string[];
};

const SelectableEntry = styled.span<{ selected: boolean }>`
  display: block;
  padding: 0.4em 0.6em;
  cursor: pointer;
  transition: all ease 0.2s;
  border-radius: 12px;
  opacity: 0.8;
  ${(p) =>
    p.selected &&
    css`
      background-color: ${Color(p.theme.accent).fade(0.3).hexa()};
    `}

  &:hover {
    background-color: ${(p) => p.theme.accent};
    opacity: 1;
  }
`;

const SelectableContainer = styled.div<{ showSelectables: boolean }>`
  width: 100%;
  background-color: ${(p) => p.theme.background2};
  border-radius: 12px;
  overflow: hidden;
  margin-top: 0.6em;
  box-shadow: 0 0 20px 0 rgba(0 0 0 / 35%);
  display: ${(p) => (p.showSelectables ? 'block' : 'none')};
  position: absolute;
`;

const InputContainer = styled.div`
  width: 100%;
  > ${Input} {
    width: 100%;
  }
  position: relative;
`;

type State = {
  selectables: string[];
  showSelectables: boolean;
  selected: number;
};

const stateReducer = (
  state: State,
  [type, payload]:
    | ['set_selectables', string[]]
    | ['set_showSelectables', boolean]
    | ['move_select', number],
) => {
  switch (type) {
    case 'set_selectables':
      return { ...state, selectables: payload };
    case 'set_showSelectables':
      return {
        ...state,
        showSelectables: payload,
        selected: payload && !state.showSelectables ? 0 : state.selected,
      };
    case 'move_select':
      let selected = state.selected + payload;
      if (selected < 0) selected = 0;
      else if (selected >= state.selectables.length) selected = state.selectables.length - 1;
      return { ...state, selected };
    default:
      return state;
  }
};

export const AutocompleteInput: React.FC<Props> = ({ value, setValue, selections }) => {
  const [state, dispatchState] = useReducer(stateReducer, {
    selectables: [],
    showSelectables: false,
    selected: 0,
  });
  const fuseRef = useRef<Fuse<string>>();

  useEffect(() => {
    fuseRef.current = new Fuse(selections, { includeScore: true });
  }, [selections]);

  useEffect(() => {
    if (!fuseRef.current) return;
    const res = fuseRef.current.search(value);
    const selectables = res
      .sort((a, b) => a.score! - b.score!)
      .map((r) => r.item)
      .slice(0, 10);

    dispatchState(['set_selectables', selectables]);
  }, [value]);

  const selectEntry = (e: React.MouseEvent<HTMLSpanElement, MouseEvent>, v: string) => {
    dispatchState(['set_showSelectables', false]);
    setValue(v);
  };

  const onInputFocus: React.FocusEventHandler<HTMLInputElement> = (e) => {
    dispatchState(['set_showSelectables', true]);
  };

  const onInputBlur = () => {
    setTimeout(() => dispatchState(['set_showSelectables', false]), 100);
  };

  const onInputKeyUp: React.KeyboardEventHandler<HTMLInputElement> = (e) => {
    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault();
        dispatchState(['move_select', 1]);
        break;
      case 'ArrowUp':
        e.preventDefault();
        dispatchState(['move_select', -1]);
        break;
      case 'Enter':
        setValue(state.selectables[state.selected]);
        break;
      case 'Escape':
        dispatchState(['set_showSelectables', false]);
        break;
      default:
        dispatchState(['set_showSelectables', true]);
    }
  };

  return (
    <InputContainer>
      <Input
        value={value}
        onInput={(e) => setValue(e.currentTarget.value)}
        onFocus={onInputFocus}
        onBlur={onInputBlur}
        onKeyUp={onInputKeyUp}
      />
      <SelectableContainer showSelectables={state.showSelectables}>
        {state.selectables.map((s, i) => (
          <SelectableEntry
            selected={i === state.selected}
            key={uid(s)}
            onClick={(e) => selectEntry(e, s)}>
            {s}
          </SelectableEntry>
        ))}
      </SelectableContainer>
    </InputContainer>
  );
};
