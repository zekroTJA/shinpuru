import React, { useEffect, useState } from 'react';
import { Trans, useTranslation } from 'react-i18next';
import { formatDate, formatSince } from '../../util/date';

import Color from 'color';
import { MaxWidthContainer } from '../../components/MaxWidthContainer';
import { SystemInfo } from '../../lib/shinpuru-ts/src';
import { byteFormatter } from 'byte-formatter';
import styled from 'styled-components';
import { useApi } from '../../hooks/useApi';

type Props = {};

const Table = styled.table`
  width: 100%;
  border-collapse: collapse;

  th {
    text-align: start;
    text-transform: uppercase;
    opacity: 0.6;
    font-weight: normal;
    font-size: 0.85rem;
    padding: 0.8em 2em 0.8em 0;
  }

  td {
    padding: 0.8em 0 0.8em 0;
  }

  tr {
    border-bottom: solid 1px ${(p) => Color(p.theme.text).alpha(0.2).hexa()};
  }
`;

const SystemRoute: React.FC<Props> = ({}) => {
  const { t } = useTranslation('routes.info.system');
  const fetch = useApi();

  const [sysinfo, setSysinfo] = useState<SystemInfo>();

  useEffect(() => {
    fetch((c) => c.etc.sysinfo())
      .then((r) => setSysinfo(r))
      .catch();
  }, []);

  return (
    <MaxWidthContainer>
      <h1>{t('heading')}</h1>
      {sysinfo && (
        <Table>
          <tbody>
            <tr>
              <th>{t('version')}</th>
              <td>
                <a
                  href={`https://github.com/zekroTJA/shinpuru/releases/tag/${sysinfo.version}`}
                  target="_blank"
                  rel="noreferrer">
                  {sysinfo.version}
                </a>
              </td>
            </tr>
            <tr>
              <th>{t('commit')}</th>
              <td>
                <a
                  href={`https://github.com/zekroTJA/shinpuru/commit/${sysinfo.commit_hash}`}
                  target="_blank"
                  rel="noreferrer">
                  {sysinfo.commit_hash.substring(0, 7)}
                </a>
              </td>
            </tr>
            <tr>
              <th>{t('builddate')}</th>
              <td>
                {formatDate(sysinfo.build_date)} (
                <Trans
                  ns="routes.info.system"
                  i18nKey="duration"
                  values={{ time: formatSince(sysinfo.build_date) }}
                />
                )
              </td>
            </tr>
            <tr>
              <th>{t('uptime')}</th>
              <td>{formatSince(sysinfo.build_date)}</td>
            </tr>
            <tr>
              <th>{t('goversion')}</th>
              <td>{sysinfo.go_version}</td>
            </tr>
            <tr>
              <th>{t('os')}</th>
              <td>{sysinfo.os}</td>
            </tr>
            <tr>
              <th>{t('arch')}</th>
              <td>{sysinfo.arch}</td>
            </tr>
            <tr>
              <th>{t('cpus')}</th>
              <td>{sysinfo.cpus}</td>
            </tr>
            <tr>
              <th>{t('goroutines')}</th>
              <td>{sysinfo.go_routines}</td>
            </tr>
            <tr>
              <th>{t('stackuse')}</th>
              <td>{byteFormatter(sysinfo.stack_use)}</td>
            </tr>
            <tr>
              <th>{t('heapuse')}</th>
              <td>{byteFormatter(sysinfo.heap_use)}</td>
            </tr>
            <tr>
              <th>{t('botid')}</th>
              <td>{sysinfo.bot_user_id}</td>
            </tr>
            <tr>
              <th>{t('invite')}</th>
              <td>
                <a href={sysinfo.bot_invite} target="_blank" rel="noreferrer">
                  {sysinfo.bot_invite}
                </a>
              </td>
            </tr>
            <tr>
              <th>{t('guilds')}</th>
              <td>{sysinfo.guilds}</td>
            </tr>
          </tbody>
        </Table>
      )}
    </MaxWidthContainer>
  );
};

export default SystemRoute;
