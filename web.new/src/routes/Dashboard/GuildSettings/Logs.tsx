import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useParams } from 'react-router-dom';
import styled from 'styled-components';
import { ReactComponent as ArrowIcon } from '../../../assets/arrow.svg';
import { ReactComponent as DeleteIcon } from '../../../assets/delete.svg';
import { Button } from '../../../components/Button';
import { Loader } from '../../../components/Loader';
import { MaxWidthContainer } from '../../../components/MaxWidthContainer';
import { Element, Select } from '../../../components/Select';
import { Small } from '../../../components/Small';
import { Switch } from '../../../components/Switch';
import { useApi } from '../../../hooks/useApi';
import { useNotifications } from '../../../hooks/useNotifications';
import { GuildLogEntry } from '../../../lib/shinpuru-ts/src';
import { formatDate } from '../../../util/date';

const PAGE_SIZE = 30;

const SEVERITY_OPTIONS: Element<number>[] = ['all', 'debug', 'info', 'warn', 'error', 'fatal'].map(
  (v, i) => ({ id: v, display: v, value: i - 1 }),
);

const TableControls = styled.div`
  padding: 0.6em;
  border-radius: 12px;
  background-color: ${(p) => p.theme.background3};

  &,
  > section {
    display: flex;
    align-items: center;
  }

  > section {
    &:nth-child(1) {
      gap: 1em;
      ${Button}:nth-child(1) > svg {
        transform: rotate(180deg);
      }
    }
    &:nth-child(2) {
      margin-left: 2em;
      > * {
        width: 6em;
      }
    }
    &:nth-child(3) {
      margin-left: auto;
    }
  }

  ${Button} {
    padding: 0.5em;
  }
`;

const Severity = styled.td<{ value: number }>`
  background-color: ${(p) => {
    switch (p.value) {
      case 1: // info
        return p.theme.cyan;
      case 2: // warn
        return p.theme.orange;
      case 3: // error
        return p.theme.red;
      case 4: // fatal
        return p.theme.pink;
      default:
        return '';
    }
  }};
`;

const Table = styled.table`
  margin-top: 1em;

  &,
  tbody,
  tr,
  th,
  td {
    border: solid 1px #6c757d;
    border-collapse: collapse;
  }

  td,
  th {
    padding: 0.4em 0.6em;
  }
`;

type Props = {};

const LogsRoute: React.FC<Props> = ({}) => {
  const { t } = useTranslation('routes.guildsettings.logs');
  const { pushNotification } = useNotifications();
  const { guildid } = useParams();
  const fetch = useApi();

  const [logs, setLogs] = useState<GuildLogEntry[]>();
  const [logsCount, setLogsCount] = useState<number>();
  const [page, setPage] = useState(0);
  const [severity, setSeverity] = useState(-1);
  const [enabled, setEnabled] = useState<boolean>();

  const _setEnabled = (v: boolean) => {};

  const fetchLogs = () => {
    if (!guildid) return;
    fetch((c) => c.guilds.settings(guildid).logs(PAGE_SIZE, page * PAGE_SIZE, severity))
      .then((res) => setLogs(res.data))
      .catch();
    fetch((c) => c.guilds.settings(guildid).logsCount(severity))
      .then((res) => setLogsCount(res.count))
      .catch();
  };

  useEffect(() => {
    if (!guildid) return;
    fetch((c) => c.guilds.settings(guildid).logsEnabled())
      .then((res) => setEnabled(res.state))
      .catch();
    fetchLogs();
  }, [guildid]);

  const pageCountStart = page * PAGE_SIZE + 1;
  let pageCountEnd = pageCountStart + PAGE_SIZE;
  if (pageCountEnd > (logsCount ?? 0)) pageCountEnd = logsCount ?? 0;

  return (
    <MaxWidthContainer>
      <h1>{t('heading')}</h1>
      <Small>{t('explanation')}</Small>

      <h2>{t('enabled.heading')}</h2>
      {(enabled !== undefined && (
        <Switch enabled={enabled} onChange={_setEnabled} labelAfter={t('enabled.enabled')} />
      )) || <Loader width="20em" height="2em" />}

      <h2>{t('entries.heading')}</h2>
      <TableControls>
        <section>
          <Button>
            <ArrowIcon />
          </Button>
          <span>
            {pageCountStart} ... {pageCountEnd} ({logsCount})
          </span>
          <Button>
            <ArrowIcon />
          </Button>
        </section>
        <section>
          <Select
            options={SEVERITY_OPTIONS}
            value={SEVERITY_OPTIONS.find((o) => o.value === severity)}
            onElementSelect={(v) => setSeverity(v.value)}
          />
        </section>
        <section>
          <Button variant="orange">
            <DeleteIcon />
            {t('entries.deleteall')}
          </Button>
        </section>
      </TableControls>
      <Table>
        <tbody>
          <tr>
            <th>{t('entries.table.timestamp')}</th>
            <th>{t('entries.table.severity')}</th>
            <th>{t('entries.table.module')}</th>
            <th>{t('entries.table.message')}</th>
          </tr>
          {logs?.map((l) => (
            <tr>
              <td>{formatDate(l.timestamp)}</td>
              <Severity value={l.severity}>
                {SEVERITY_OPTIONS.find((v) => v.value === l.severity)?.display ?? ''}
              </Severity>
              <td>{l.module}</td>
              <td>{l.message}</td>
            </tr>
          ))}
        </tbody>
      </Table>
    </MaxWidthContainer>
  );
};

export default LogsRoute;
