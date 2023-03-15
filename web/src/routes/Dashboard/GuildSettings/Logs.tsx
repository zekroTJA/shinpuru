import { Element, Select } from '../../../components/Select';
import { useEffect, useState } from 'react';

import { ReactComponent as ArrowIcon } from '../../../assets/arrow.svg';
import { Button } from '../../../components/Button';
import { ReactComponent as DeleteIcon } from '../../../assets/delete.svg';
import { GuildLogEntry } from '../../../lib/shinpuru-ts/src';
import { Loader } from '../../../components/Loader';
import { MaxWidthContainer } from '../../../components/MaxWidthContainer';
import { Modal } from '../../../components/Modal';
import { Small } from '../../../components/Small';
import { Switch } from '../../../components/Switch';
import { formatDate } from '../../../util/date';
import styled from 'styled-components';
import { uid } from 'react-uid';
import { useApi } from '../../../hooks/useApi';
import { useNotifications } from '../../../hooks/useNotifications';
import { useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

const PAGE_SIZE = 14;

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
  width: 100%;

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

const LogsRoute: React.FC<Props> = () => {
  const { t } = useTranslation('routes.guildsettings.logs');
  const { pushNotification } = useNotifications();
  const { guildid } = useParams();
  const fetch = useApi();

  const [logs, setLogs] = useState<GuildLogEntry[]>();
  const [logsCount, setLogsCount] = useState<number>();
  const [page, setPage] = useState(0);
  const [severity, setSeverity] = useState(-1);
  const [enabled, setEnabled] = useState<boolean>();
  const [showDeleteModal, setShowDeleteModal] = useState(false);

  const _setEnabled = (v: boolean) => {
    if (!guildid) return;
    setEnabled(v);
    fetch((c) => c.guilds.settings(guildid).setLogsEnabled(v))
      .then(() =>
        pushNotification({
          message: v ? t('notifications.enabled') : t('notifications.disabled'),
          type: 'SUCCESS',
        }),
      )
      .catch(() => setEnabled(!v));
  };

  const _setPage = (rel: -1 | 1) => {
    if (!guildid || !logsCount) return;

    let newPage = page + rel;
    if (newPage < 0) newPage = 0;
    else if (newPage * PAGE_SIZE >= logsCount) newPage = page;

    setPage(newPage);
  };

  const _setSeverity = (severity: number) => {
    setPage(0);
    setSeverity(severity);
    setLogs([]);
    setLogsCount(0);
  };

  const _deleteEntries = () => {
    setShowDeleteModal(false);
    if (!guildid) return;
    fetch((c) => c.guilds.settings(guildid).flushLogs())
      .then(() => {
        setPage(0);
        pushNotification({
          message: t('notifications.deleted'),
          type: 'SUCCESS',
        });
      })
      .catch();
  };

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
  }, [guildid]);

  useEffect(() => {
    if (!guildid) return;
    fetchLogs();
  }, [guildid, page, severity]);

  const pageCountStart = page * PAGE_SIZE + 1;
  let pageCountEnd = pageCountStart + PAGE_SIZE;
  if (pageCountEnd > (logsCount ?? 0)) pageCountEnd = logsCount ?? 0;

  return (
    <>
      <Modal
        show={showDeleteModal}
        onClose={() => setShowDeleteModal(false)}
        heading={t('deletemodal.heading')}
        controls={
          <>
            <Button onClick={_deleteEntries}>{t('deletemodal.controls.delete')}</Button>
            <Button onClick={() => setShowDeleteModal(false)} variant="gray">
              {t('deletemodal.controls.cancel')}
            </Button>
          </>
        }>
        <span>{t('deletemodal.content')}</span>
      </Modal>

      <MaxWidthContainer>
        <h1>{t('heading')}</h1>
        <Small>{t('explanation')}</Small>

        <h2>{t('enabled.heading')}</h2>
        {(enabled !== undefined && (
          <Switch enabled={enabled} onChange={_setEnabled} labelAfter={t('enabled.enabled')} />
        )) || <Loader width="20em" height="2em" />}

        <h2>{t('entries.heading')}</h2>
        {(logs !== undefined && logsCount !== undefined && (
          <>
            <TableControls>
              <section>
                <Button onClick={() => _setPage(-1)}>
                  <ArrowIcon />
                </Button>
                <span>
                  {pageCountStart} ... {pageCountEnd} ({logsCount})
                </span>
                <Button onClick={() => _setPage(1)}>
                  <ArrowIcon />
                </Button>
              </section>
              <section>
                <Select
                  options={SEVERITY_OPTIONS}
                  value={SEVERITY_OPTIONS.find((o) => o.value === severity)}
                  onElementSelect={(v) => _setSeverity(v.value)}
                />
              </section>
              <section>
                <Button variant="orange" onClick={() => setShowDeleteModal(true)}>
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
                  <tr key={uid(l)}>
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
          </>
        )) || (
          <>
            <Loader width="100%" height="3em" />
            <Loader width="100%" height="2.5em" margin="1em 0 0 0" />
            <Loader width="100%" height="2.5em" margin="0.5em 0 0 0" />
            <Loader width="100%" height="2.5em" margin="0.5em 0 0 0" />
            <Loader width="100%" height="2.5em" margin="0.5em 0 0 0" />
            <Loader width="100%" height="2.5em" margin="0.5em 0 0 0" />
          </>
        )}
      </MaxWidthContainer>
    </>
  );
};

export default LogsRoute;
