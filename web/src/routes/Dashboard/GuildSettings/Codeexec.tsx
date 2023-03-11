import React, { useEffect, useReducer } from 'react';
import { Trans, useTranslation } from 'react-i18next';
import styled, { useTheme } from 'styled-components';

import { ActionButton } from '../../../components/ActionButton';
import { CodeExecSettings } from '../../../lib/shinpuru-ts/src';
import { Controls } from '../../../components/Controls';
import { Hint } from '../../../components/Hint';
import { Input } from '../../../components/Input';
import { Loader } from '../../../components/Loader';
import { MaxWidthContainer } from '../../../components/MaxWidthContainer';
import { Small } from '../../../components/Small';
import { Switch } from '../../../components/Switch';
import { useApi } from '../../../hooks/useApi';
import { useNotifications } from '../../../hooks/useNotifications';
import { useParams } from 'react-router';

type Props = {};

const settingsReducer = (
  state: Partial<CodeExecSettings>,
  [type, payload]:
    | ['set_state', Partial<CodeExecSettings>]
    | ['set_enabled', boolean]
    | ['set_clientid' | 'set_clientsecret', string],
) => {
  switch (type) {
    case 'set_state':
      return { ...state, ...payload };
    case 'set_enabled':
      return { ...state, enabled: payload };
    case 'set_clientid':
      return { ...state, jdoodle_clientid: payload };
    case 'set_clientsecret':
      return { ...state, jdoodle_clientsecret: payload };
    default:
      return state;
  }
};

const InputContainer = styled.div`
  margin-top: 1.5em;

  > input,
  label {
    display: block;
    width: 100%;
  }

  > label {
    margin-bottom: 1em;
  }
`;

const MarginSmall = styled(Small)`
  margin-top: 1.5em;
`;

const VerificationRoute: React.FC<Props> = () => {
  const { t } = useTranslation('routes.guildsettings.codeexec');
  const theme = useTheme();
  const { pushNotification } = useNotifications();
  const { guildid } = useParams();
  const fetch = useApi();
  const [settings, dispatchSettings] = useReducer(settingsReducer, {} as Partial<CodeExecSettings>);

  const _saveSettings = () => {
    if (!guildid) return;

    return fetch((c) => c.guilds.settings(guildid).setCodeexec(settings as CodeExecSettings))
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

    fetch((c) => c.guilds.settings(guildid).codeexec())
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

      {(settings.type === 'ranna' && (
        <>
          <MarginSmall>
            <Trans
              ns="routes.guildsettings.codeexec"
              i18nKey="ranna.explanation"
              components={{
                '1': <a href="https://app.ranna.zekro.de/">_</a>,
              }}
            />
          </MarginSmall>
        </>
      )) || (
        <>
          <h3>{t('jdoodle.heading')}</h3>
          <Small>
            <Trans
              ns="routes.guildsettings.codeexec"
              i18nKey="jdoodle.explanation"
              components={{
                '1': <a href="https://www.jdoodle.com/">_</a>,
                '2': <a href="https://www.jdoodle.com/compiler-api/">_</a>,
              }}
            />
          </Small>
          {settings.enabled && !settings.jdoodle_clientid && !settings.jdoodle_clientsecret && (
            <Hint color={theme.orange}>{t('jdoodle.settingswarn')}</Hint>
          )}
          <InputContainer>
            <label>{t('jdoodle.clientid')}</label>
            {(settings.type !== undefined && (
              <Input
                value={settings.jdoodle_clientid}
                onInput={(e) => dispatchSettings(['set_clientid', e.currentTarget.value])}
              />
            )) || <Loader width="100%" height="2em" />}
          </InputContainer>
          <InputContainer>
            <label>{t('jdoodle.clientsecret')}</label>
            {(settings.type !== undefined && (
              <Input
                value={settings.jdoodle_clientsecret}
                type="password"
                onInput={(e) => dispatchSettings(['set_clientsecret', e.currentTarget.value])}
              />
            )) || <Loader width="100%" height="2em" />}
          </InputContainer>
        </>
      )}

      <Controls>
        <ActionButton variant="green" onClick={_saveSettings}>
          {t('save')}
        </ActionButton>
      </Controls>
    </MaxWidthContainer>
  );
};

export default VerificationRoute;
