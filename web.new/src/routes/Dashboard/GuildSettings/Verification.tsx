import React, { useEffect, useReducer } from 'react';

import { Button } from '../../../components/Button';
import { Controls } from '../../../components/Controls';
import { GuildSettingsVerification } from '../../../lib/shinpuru-ts/src';
import { Loader } from '../../../components/Loader';
import { MaxWidthContainer } from '../../../components/MaxWidthContainer';
import { Small } from '../../../components/Small';
import { Switch } from '../../../components/Switch';
import { useApi } from '../../../hooks/useApi';
import { useNotifications } from '../../../hooks/useNotifications';
import { useParams } from 'react-router';
import { useTranslation } from 'react-i18next';

type Props = {};

const settingsReducer = (
  state: Partial<GuildSettingsVerification>,
  [type, payload]: ['set_state', Partial<GuildSettingsVerification>] | ['set_enabled', boolean],
) => {
  switch (type) {
    case 'set_state':
      return { ...state, ...payload };
    case 'set_enabled':
      return { ...state, enabled: payload };
    default:
      return state;
  }
};

const CodeexecRoute: React.FC<Props> = () => {
  const { t } = useTranslation('routes.guildsettings.verification');
  const { pushNotification } = useNotifications();
  const { guildid } = useParams();
  const fetch = useApi();
  const [settings, dispatchSettings] = useReducer(
    settingsReducer,
    {} as Partial<GuildSettingsVerification>,
  );

  const _saveSettings = () => {
    if (!guildid) return;

    fetch((c) => c.guilds.settings(guildid).setVerification(settings as GuildSettingsVerification))
      .then(() =>
        pushNotification({
          message: t('notifications.saved'),
          type: 'SUCCESS',
        }),
      )
      .catch();
  };

  useEffect(() => {
    if (!guildid) return;

    fetch((c) => c.guilds.settings(guildid).verification())
      .then((res) => dispatchSettings(['set_state', res]))
      .catch();
  }, [guildid]);

  return (
    <MaxWidthContainer>
      <h1>{t('heading')}</h1>
      <Small>{t('explanation')}</Small>

      <h2>Settings</h2>
      {(settings.enabled !== undefined && (
        <Switch
          enabled={settings.enabled}
          onChange={(e) => dispatchSettings(['set_enabled', e])}
          labelAfter={t('enable')}
        />
      )) || <Loader width="20em" height="2em" />}

      <Controls>
        <Button variant="green" onClick={_saveSettings}>
          {t('save')}
        </Button>
      </Controls>
    </MaxWidthContainer>
  );
};

export default CodeexecRoute;
