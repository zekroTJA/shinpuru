import { Embed } from './Embed';
import { JoinlogEntry } from '../lib/shinpuru-ts/src';
import { SinceDate } from './SinceDate';
import { formatDate } from '../util/date';
import styled from 'styled-components';

const StyledTr = styled.tr<{ selected: boolean }>`
  cursor: pointer;
  td {
    text-align: start;
    padding: 0.8em;
    background-color: ${(p) => (p.selected ? p.theme.accentDarker : p.theme.background2)};
    overflow: hidden;
    &:first-child {
      border-top-left-radius: 8px;
      border-bottom-left-radius: 8px;
    }
    &:last-child {
      border-top-right-radius: 8px;
      border-bottom-right-radius: 8px;
    }
  }
`;

type Props = {
  entry: JoinlogEntry;
  selected: boolean;
  onCheck: (checked: boolean, entry: JoinlogEntry) => void;
};

export const JoinLogEntry: React.FC<Props> = ({ entry, selected, onCheck }) => {
  return (
    <StyledTr selected={selected} onClick={() => onCheck(!selected, entry)}>
      <td>
        <Embed>{entry.user_id}</Embed>
      </td>
      <td>{entry.tag}</td>
      <td>
        <SinceDate date={entry.account_created} />
      </td>
      <td>{formatDate(entry.timestamp)}</td>
      <td>
        <input
          type="checkbox"
          checked={selected}
          onChange={(v) => onCheck(v.currentTarget.checked, entry)}
        />
      </td>
    </StyledTr>
  );
};
