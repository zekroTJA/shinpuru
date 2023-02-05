import { PropsWithChildren } from 'react';
import styled from 'styled-components';
import { Styled } from '../props';

type Theme = {
  disabledColor: string;
  enabledColor: string;
};

type Props = Styled &
  PropsWithChildren<{
    enabled: boolean;
    onChange?: (s: boolean) => void;
    labelBefore?: string | JSX.Element;
    labelAfter?: string | JSX.Element;
    theaming?: Partial<Theme>;
  }>;

const SwitchContainer = styled.div<{ enabled: boolean; theaming?: Partial<Theme> }>`
  width: fit-content;
  display: flex;
  align-items: center;
  gap: 1em;

  &,
  > * {
    cursor: pointer;
  }

  > div {
    min-width: 4em;
    height: 2em;
    border-radius: 1em;
    background-color: ${(p) =>
      p.enabled
        ? p.theaming?.enabledColor ?? p.theme.accent
        : p.theaming?.disabledColor ?? p.theme.background3};
    transition: all 0.3s ease;
    padding: 0.25em;

    > div {
      height: 100%;
      border-radius: 2em;
      background-color: ${(p) => p.theme.white};
      margin-left: ${(p) => (p.enabled ? '2em' : '0')};
      margin-right: ${(p) => (p.enabled ? '0' : '2em')};

      transition: ${(p) => (p.enabled ? 'margin-left' : 'margin-right')} 0.3s ease 0.15s,
        ${(p) => (p.enabled ? 'margin-right' : 'margin-left')} 0.3s ease;
    }
  }
`;

export const Switch: React.FC<Props> = ({
  enabled,
  onChange = () => {},
  labelBefore,
  labelAfter,
  children,
  ...props
}) => {
  const toLabel = (v?: string | JSX.Element) => {
    return v ? typeof v === 'string' ? <label>{v}</label> : v : <></>;
  };

  return (
    <SwitchContainer enabled={enabled} onClick={() => onChange(!enabled)} {...props}>
      {toLabel(labelBefore)}
      <div>
        <div>{children}</div>
      </div>
      {toLabel(labelAfter)}
    </SwitchContainer>
  );
};
