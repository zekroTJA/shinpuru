import React, { useEffect, useRef, useState } from 'react';

import { ActionButton } from '../../../components/ActionButton';
import { Button } from '../../../components/Button';
import { ReactComponent as DownloadIcon } from '../../../assets/download.svg';
import { Embed } from '../../../components/Embed';
import { GuildBackup } from '../../../lib/shinpuru-ts/src';
import { Loader } from '../../../components/Loader';
import { MaxWidthContainer } from '../../../components/MaxWidthContainer';
import { Switch } from '../../../components/Switch';
import { Tag } from '../../../components/Tag';
import { formatDate } from '../../../util/date';
import { range } from '../../../util/utils';
import styled from 'styled-components';
import { uid } from 'react-uid';
import { useApi } from '../../../hooks/useApi';
import { useGuild } from '../../../hooks/useGuild';
import { useNotifications } from '../../../hooks/useNotifications';
import { useParams } from 'react-router';
import { useTranslation } from 'react-i18next';

type Props = {};

const ControlContainer = styled.div`
  display: flex;
  margin-bottom: 2em;

  > ${Button} {
    padding: 0 0.8em;
    margin-left: auto;
  }
`;

const BackupContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.5em;
`;

const BackupTile = styled.div<{ latest: boolean }>`
  display: flex;
  gap: 1em;
  align-items: center;
  background-color: ${(p) => (p.latest ? p.theme.accentDarker : p.theme.background3)};
  border-radius: 3px;
  padding: 0.5em;

  &:hover > svg {
    opacity: 1;
  }

  > svg {
    margin-left: auto;
    opacity: 0;
    transition: opacity 0.2s ease;
    cursor: pointer;
  }
`;

const BackupRoute: React.FC<Props> = () => {
  const { t } = useTranslation('routes.guildsettings.backups');
  const { pushNotification } = useNotifications();
  const { guildid } = useParams();
  const guild = useGuild(guildid);
  const [enabled, setEnabled] = useState(false);
  const [backups, setBackups] = useState<GuildBackup[]>();
  const currStateRef = useRef(false);
  const fetch = useApi();

  const _applyEnabled = () => {
    if (!guildid) return;
    return fetch((c) => c.guilds.backups(guildid).toggle(enabled))
      .then(() => {
        currStateRef.current = enabled;
        pushNotification({
          message: t<string>(enabled ? 'notifications.enabled' : 'notifications.disabled'),
          type: enabled ? 'SUCCESS' : 'WARNING',
        });
      })
      .catch(() => setEnabled(currStateRef.current));
  };

  const _download = async (backup: GuildBackup) => {
    if (!guildid) return;
    const token = await fetch((c) => c.guilds.backups(guildid).download(backup.file_id));
    const link = document.createElement('a');
    link.href = await fetch((c) =>
      c.guilds.backups(guildid).downloadUrl(backup.file_id, token.token),
    );
    link.click();
    document.removeChild(link);
  };

  useEffect(() => {
    if (!guild) return;
    setEnabled(guild.backups_enabled);
    currStateRef.current = guild.backups_enabled;

    fetch((c) => c.guilds.backups(guild.id).list())
      .then((res) => setBackups(res.data))
      .catch();
  }, [guild]);

  const _backups =
    backups?.map((b, i, l) => (
      <BackupTile key={uid(b)} latest={i + 1 === l.length}>
        <Tag>{i}</Tag>
        <span>{formatDate(b.timestamp)}</span>
        <Embed>{b.file_id}</Embed>
        <DownloadIcon onClick={() => _download(b)} />
      </BackupTile>
    )) ?? range(10).map((n) => <Loader key={`loader-${n}`} height="2em" width="100%" />);

  return (
    <MaxWidthContainer>
      <h1>{t('heading')}</h1>
      <h2>{t('settings')}</h2>
      <ControlContainer>
        {(guild && (
          <>
            <Switch enabled={enabled} onChange={setEnabled} labelAfter={t<string>('toggle')} />
            <ActionButton
              disabled={enabled === currStateRef.current}
              variant="green"
              onClick={_applyEnabled}>
              {t('apply')}
            </ActionButton>
          </>
        )) || (
          <>
            <Loader width="4em" height="2em" />
            <Loader width="10em" height="2em" margin="0 0 0 1em" />
            <Loader width="4em" height="2em" margin="0 0 0 auto" />
          </>
        )}
      </ControlContainer>

      <h2>{t('backups')}</h2>
      <BackupContainer>
        {_backups.length === 0 ? <i>No backups available.</i> : _backups}
      </BackupContainer>
    </MaxWidthContainer>
  );
};

export default BackupRoute;
