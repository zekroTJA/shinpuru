import React, { useEffect, useReducer } from 'react';
import { useTranslation } from 'react-i18next';
import { useParams } from 'react-router';
import styled from 'styled-components';
import { MaxWidthContainer } from '../../../components/MaxWidthContainer';
import { Small } from '../../../components/Small';
import { Switch } from '../../../components/Switch';
import { useApi } from '../../../hooks/useApi';
import { useNotifications } from '../../../hooks/useNotifications';
import { AntiraidSettings } from '../../../lib/shinpuru-ts/src';

type Props = {};

const settingsReducer = (
  state: AntiraidSettings,
  [type, payload]:
    | ['set_state', Partial<AntiraidSettings>]
    | ['set_enabled' | 'set_verification', boolean]
    | ['set_regeneration' | 'set_burst', number],
) => {
  switch (type) {
    case 'set_state':
      return { ...state, ...payload };
    case 'set_enabled':
      return { ...state, state: payload };
    case 'set_verification':
      return { ...state, verification: payload };
    case 'set_regeneration':
      return { ...state, regeneration_period: payload };
    case 'set_burst':
      return { ...state, burst: payload };
    default:
      return state;
  }
};

const AntiraidRoute: React.FC<Props> = ({}) => {
  const { t } = useTranslation('routes.guildsettings.antiraid');
  const { pushNotification } = useNotifications();
  const { guildid } = useParams();
  const fetch = useApi();
  const [settings, dispatchSettings] = useReducer(settingsReducer, {} as AntiraidSettings);

  useEffect(() => {
    if (!guildid) return;

    fetch((c) => c.guilds.settings(guildid).antiraid())
      .then((res) => dispatchSettings(['set_state', res]))
      .catch();
  }, [guildid]);

  return (
    <MaxWidthContainer>
      <h1>{t('heading')}</h1>
      <Small>{t('explaination')}</Small>
      <Switch
        enabled={settings.state}
        onChange={(e) => dispatchSettings(['set_enabled', e])}
        labelAfter={t('enable')}
      />
    </MaxWidthContainer>
  );
};

export default AntiraidRoute;
