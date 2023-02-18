import React, { useEffect, useState } from 'react';
import { Trans, useTranslation } from 'react-i18next';

import { Button } from '../../../components/Button';
import { Flex } from '../../../components/Flex';
import { Input } from '../../../components/Input';
import { MaxWidthContainer } from '../../../components/MaxWidthContainer';
import { Modal } from '../../../components/Modal';
import { ReactMarkdown } from 'react-markdown/lib/react-markdown';
import { Switch } from '../../../components/Switch';
import styled from 'styled-components';
import { useApi } from '../../../hooks/useApi';
import { useGuild } from '../../../hooks/useGuild';
import { useNotifications } from '../../../hooks/useNotifications';
import { useParams } from 'react-router';

type Props = {};

const Controls = styled(Flex)`
  margin-top: 2em;
`;

const GuildNameInput = styled(Input)`
  width: 100%;
  background-color: ${(p) => p.theme.background3};
`;

const DataRoute: React.FC<Props> = ({}) => {
  const { t } = useTranslation('routes.guildsettings.data');
  const { pushNotification } = useNotifications();
  const { guildid } = useParams();
  const guild = useGuild(guildid);
  const fetch = useApi();
  const [kick, setKick] = useState(false);
  const [showModal, setShowModal] = useState(false);
  const [guildName, setGuildName] = useState('');

  const _confirmDelete = () => {
    if (!guildid || !guildName) return;
    fetch((c) => c.guilds.settings(guildid).flushData(kick, guildName))
      .then(() => {
        pushNotification({ message: t('notifications.success'), type: 'SUCCESS' });
      })
      .catch()
      .finally(() => setShowModal(false));
  };

  useEffect(() => {
    if (showModal) setGuildName('');
  }, [showModal]);

  return (
    <>
      <Modal
        show={showModal}
        onClose={() => setShowModal(false)}
        heading={t('modal.heading')}
        controls={
          <>
            <Button disabled={guild?.name !== guildName} onClick={_confirmDelete}>
              Confirm
            </Button>
            <Button variant="gray" onClick={() => setShowModal(false)}>
              Cancel
            </Button>
          </>
        }>
        <Trans ns="routes.guildsettings.data" i18nKey="modal.content">
          <strong>all</strong>
        </Trans>
        <p>
          <GuildNameInput
            placeholder={guild?.name}
            value={guildName}
            onInput={(e) => setGuildName(e.currentTarget.value)}
          />
        </p>
      </Modal>
      <MaxWidthContainer>
        <h1>{t('heading')}</h1>
        <ReactMarkdown children={t('explanation')} />

        <Controls gap="1em">
          <Button variant="red" onClick={() => setShowModal(true)}>
            {t('delete')}
          </Button>
          <Switch labelAfter={t('remove')} enabled={kick} onChange={setKick} />
        </Controls>
      </MaxWidthContainer>
    </>
  );
};

export default DataRoute;
