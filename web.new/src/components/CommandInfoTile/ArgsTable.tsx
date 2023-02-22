import { CommandOption, CommandOptionType } from '../../lib/shinpuru-ts/src';

import Color from 'color';
import { Embed } from '../Embed';
import styled from 'styled-components';
import { uid } from 'react-uid';
import { useTranslation } from 'react-i18next';

type Props = {
  options: CommandOption[];
};

const StyledTable = styled.table`
  width: 100%;
  border-collapse: collapse;

  th {
    text-transform: uppercase;
    opacity: 0.6;
    font-size: 0.7rem;
    text-align: start;
    padding-right: 2em;
  }

  th,
  td {
    padding: 0.3em 0;
  }

  tr {
    border-bottom: solid 1px ${(p) => Color(p.theme.text).alpha(0.3).hexa()};
  }
`;

export const ArgsTable: React.FC<Props> = ({ options }) => {
  const { t } = useTranslation('components', { keyPrefix: 'commandtile.args' });

  return (
    <StyledTable>
      <tbody>
        <tr>
          <th>{t('argument')}</th>
          <th>{t('type')}</th>
          <th>{t('required')}</th>
          <th>{t('description')}</th>
        </tr>
        {options.map((o) => (
          <tr key={uid(o)}>
            <td>{o.name}</td>
            <td>
              <Embed>{typeToString(o.type)}</Embed>
            </td>
            <td>{t(o.required ? 'yes' : 'no')}</td>
            <td>{o.description}</td>
          </tr>
        ))}
      </tbody>
    </StyledTable>
  );
};

const typeToString = (t: CommandOptionType) => {
  switch (t) {
    case CommandOptionType.BOOLEAN:
      return 'boolean';
    case CommandOptionType.CHANNEL:
      return 'channel';
    case CommandOptionType.INTEGER:
      return 'integer';
    case CommandOptionType.MENTIONABLE:
      return 'mentionable';
    case CommandOptionType.ROLE:
      return 'role';
    case CommandOptionType.STRING:
      return 'string';
    case CommandOptionType.USER:
      return 'user';
    default:
      return 'unknown';
  }
};
