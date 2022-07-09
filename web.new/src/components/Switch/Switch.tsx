import styled from 'styled-components';

type Props = {
  enabled: boolean;
  onChange?: (s: boolean) => void;
  labelBefore?: string | JSX.Element;
  labelAfter?: string | JSX.Element;
};

const SwitchContainer = styled.div<{ enabled: boolean }>`
  width: fit-content;
  display: flex;
  align-items: center;
  gap: 1em;

  &,
  > * {
    cursor: pointer;
  }

  > div {
    width: 4em;
    height: 2em;
    border-radius: 1em;
    background-color: ${(p) => (p.enabled ? p.theme.accent : p.theme.background3)};
    transition: all 0.2s ease;
    padding: 0.25em;

    > div {
      height: 100%;
      border-radius: 2em;
      background-color: ${(p) => p.theme.white};
      margin-left: ${(p) => (p.enabled ? '2em' : '0')};
      margin-right: ${(p) => (p.enabled ? '0' : '2em')};

      transition: ${(p) => (p.enabled ? 'margin-left' : 'margin-right')} 0.3s ease;
    }
  }
`;

export const Switch: React.FC<Props> = ({
  enabled,
  onChange = () => {},
  labelBefore,
  labelAfter,
}) => {
  const toLabel = (v?: string | JSX.Element) => {
    return v ? typeof v === 'string' ? <label>{v}</label> : v : <></>;
  };

  return (
    <SwitchContainer enabled={enabled} onClick={() => onChange(!enabled)}>
      {toLabel(labelBefore)}
      <div>
        <div></div>
      </div>
      {toLabel(labelAfter)}
    </SwitchContainer>
  );
};
