import React, { useEffect, useReducer, useState } from 'react';
import styled, { useTheme } from 'styled-components';

import { Button } from '../../components/Button';
import { Card } from '../../components/Card';
import { Controls } from '../../components/Controls';
import { Embed } from '../../components/Embed';
import { Input } from '../../components/Input';
import { Loader } from '../../components/Loader';
import { MaxWidthContainer } from '../../components/MaxWidthContainer';
import { Modal } from '../../components/Modal';
import ReactMarkdown from 'react-markdown';
import { Small } from '../../components/Small';
import { Switch } from '../../components/Switch';
import { UserSettingsPrivacy } from '../../lib/shinpuru-ts/src';
import { useApi } from '../../hooks/useApi';
import { useNavigate } from 'react-router';
import { useNotifications } from '../../hooks/useNotifications';
import { useSelfUser } from '../../hooks/useSelfUser';
import { useTranslation } from 'react-i18next';

type Props = {};

const PurgeCard = styled(Card)`
  width: 100%;
  margin-top: 2em;

  > h2 {
    margin-top: 0;
  }
`;

const PurgeModalInput = styled(Input)`
  display: block;
  width: 100%;
  margin-top: 1em;
  background-color: ${(p) => p.theme.background3};
`;

const privacyReducer = (
  state: Partial<UserSettingsPrivacy>,
  [type, payload]: ['set_state', Partial<UserSettingsPrivacy>] | ['set_optout', boolean],
) => {
  switch (type) {
    case 'set_state':
      return { ...state, ...payload };
    case 'set_optout':
      return { ...state, starboard_optout: payload };
    default:
      return state;
  }
};

const PrivacyRoute: React.FC<Props> = () => {
  const { t } = useTranslation('routes.usersettings.privacy');
  const { pushNotification } = useNotifications();
  const theme = useTheme();
  const nav = useNavigate();
  const self = useSelfUser();
  const fetch = useApi();
  const [state, dispatchState] = useReducer(privacyReducer, {});
  const [showPurgeModal, setShowPurgeModal] = useState(false);
  const [purgeUsername, setPurgeUsername] = useState('');

  const _saveSettings = () => {
    fetch((c) => c.usersettings.setPrivacy(state as UserSettingsPrivacy))
      .then(() =>
        pushNotification({
          message: t('notifications.saved'),
          type: 'SUCCESS',
        }),
      )
      .catch();
  };

  const _showPurgeModal = () => {
    setPurgeUsername('');
    setShowPurgeModal(true);
  };

  const _purgeUserData = () => {
    fetch((c) => c.usersettings.flush())
      .then(() => {
        nav('/login');
        pushNotification({
          message: t('notifications.purged'),
          type: 'SUCCESS',
        });
      })
      .catch();
  };

  useEffect(() => {
    fetch((c) => c.usersettings.privacy())
      .then((r) => dispatchState(['set_state', r]))
      .catch();
  }, []);

  return (
    <>
      <Modal
        show={showPurgeModal}
        onClose={() => setShowPurgeModal(false)}
        heading={t('purge.modal.heading')}
        controls={
          <>
            <Button
              variant="red"
              disabled={!purgeUsername || purgeUsername !== self?.username}
              onClick={_purgeUserData}>
              {t('purge.modal.confirm')}
            </Button>
            <Button variant="gray" onClick={() => setShowPurgeModal(false)}>
              {t('purge.modal.cancel')}
            </Button>
          </>
        }>
        <ReactMarkdown children={t('purge.modal.explanation')} />
        <Embed>{self?.username}</Embed>
        <PurgeModalInput
          placeholder={t('purge.modal.placeholder')}
          value={purgeUsername}
          onInput={(e) => setPurgeUsername(e.currentTarget.value)}
        />
      </Modal>

      <MaxWidthContainer>
        <h1>{t('heading')}</h1>
        <Small>{t('explanation')}</Small>

        <section>
          <h2>{t('starboard.heading')}</h2>
          <Small>
            <ReactMarkdown children={t('starboard.explanation')} />
          </Small>
          {(state.starboard_optout !== undefined && (
            <Switch
              enabled={state.starboard_optout}
              onChange={(v) => dispatchState(['set_optout', v])}
              labelAfter={t('starboard.toggle')}
            />
          )) || <Loader width="15em" height="2.5em" />}
        </section>

        <Controls>
          <Button variant="green" onClick={_saveSettings}>
            {t('starboard.save')}
          </Button>
        </Controls>

        <section>
          <PurgeCard color={theme.red}>
            <h2>{t('purge.heading')}</h2>
            <ReactMarkdown children={t('purge.explanation')} />
            <Controls>
              <Button variant="red" onClick={_showPurgeModal}>
                {t('purge.purge')}
              </Button>
            </Controls>
          </PurgeCard>
        </section>
      </MaxWidthContainer>
    </>
  );
};

export default PrivacyRoute;
