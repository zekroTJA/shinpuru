import React, { useEffect, useReducer } from 'react';
import { Trans, useTranslation } from 'react-i18next';
import styled, { css } from 'styled-components';

import { Button } from '../../../components/Button';
import { Controls } from '../../../components/Controls';
import { Embed } from '../../../components/Embed';
import { Flex } from '../../../components/Flex';
import { GuildSettingsApi } from '../../../lib/shinpuru-ts/src';
import { Hint } from '../../../components/Hint';
import { Input } from '../../../components/Input';
import { Loader } from '../../../components/Loader';
import { MaxWidthContainer } from '../../../components/MaxWidthContainer';
import { ReactComponent as ShieldIcon } from '../../../assets/shield.svg';
import { Small } from '../../../components/Small';
import { Switch } from '../../../components/Switch';
import { ReactComponent as WarnIcon } from '../../../assets/warning.svg';
import { useApi } from '../../../hooks/useApi';
import { useNotifications } from '../../../hooks/useNotifications';
import { useParams } from 'react-router';

type Props = {};

const Endpoint = styled.a<{ disabled: boolean }>`
  background-color: ${(p) => p.theme.background2};
  color: ${(p) => p.theme.accent};
  cursor: pointer;
  padding: 1em 1.2em;
  border-radius: 8px;
  width: fit-content;
  display: block;

  ${(p) =>
    p.disabled &&
    css`
      cursor: not-allowed;
      opacity: 0.5;
      color: ${(p) => p.theme.text};
    `}
`;

const Settings = styled(Flex)`
  ${Input} {
    width: 100%;
  }
`;

const StyledHint = styled(Hint)`
  margin-bottom: 1em;

  > span {
    display: flex;
    align-items: center;
    height: 100%;
    gap: 0.2em;

    svg {
      height: 1.2em;
      width: 1.2em;
    }
  }
`;

const settingsReducer = (
  state: Partial<GuildSettingsApi>,
  [type, payload]:
    | ['set_state', Partial<GuildSettingsApi>]
    | ['set_enabled', boolean]
    | ['set_allowedorigins' | 'set_token', string],
) => {
  switch (type) {
    case 'set_state':
      return { ...state, ...payload };
    case 'set_enabled':
      return { ...state, enabled: payload };
    case 'set_allowedorigins':
      return { ...state, allowed_origins: payload };
    case 'set_token':
      return { ...state, token: payload };
    default:
      return state;
  }
};

const APIRoute: React.FC<Props> = ({}) => {
  const { t } = useTranslation('routes.guildsettings.api');
  const { pushNotification } = useNotifications();
  const { guildid } = useParams();
  const fetch = useApi();
  const [settings, dispatchSettings] = useReducer(settingsReducer, {} as GuildSettingsApi);

  const _refresh = () => {
    if (!guildid) return;

    dispatchSettings(['set_token', '']);

    return fetch((c) => c.guilds.settings(guildid).api())
      .then((res) => dispatchSettings(['set_state', res]))
      .catch();
  };

  const _saveSettings = () => {
    if (!guildid) return;

    fetch((c) => c.guilds.settings(guildid).setApi(settings as GuildSettingsApi))
      .then(_refresh)
      .then(() =>
        pushNotification({
          message: t('notifications.saved'),
          type: 'SUCCESS',
        }),
      )
      .catch();
  };

  const _resetToken = () => {
    if (!guildid || !settings.protected) return;

    fetch((c) =>
      c.guilds.settings(guildid).setApi({ ...settings, reset_token: true } as GuildSettingsApi),
    )
      .then(_refresh)
      .then(() =>
        pushNotification({
          message: t('notifications.reset'),
          type: 'WARNING',
        }),
      )
      .catch();
  };

  useEffect(() => {
    _refresh();
  }, [guildid]);

  const endpoint = `${window.origin}/api/public/guilds/${guildid}${
    settings?.protected ? 'token={token}' : ''
  }`;

  return (
    <MaxWidthContainer>
      <h1>{t('heading')}</h1>
      <Small>{t('explanation')}</Small>

      {(settings.enabled !== undefined && (
        <Endpoint
          href={settings.enabled ? endpoint : undefined}
          target="_blank"
          disabled={!settings.enabled}>
          GET {endpoint}
        </Endpoint>
      )) || <Loader width="30em" height="3em" />}

      <h2>{t('settings.heading')}</h2>
      <Settings gap="1em" direction="column">
        {(settings.enabled !== undefined && (
          <Switch
            enabled={settings.enabled}
            onChange={(e) => dispatchSettings(['set_enabled', e])}
            labelAfter={t('settings.enable')}
          />
        )) || <Loader width="20em" height="2em" />}
        <div>
          <h3>{t('settings.allowedorigins.heading')}</h3>
          <Small>
            <Trans ns="routes.guildsettings.api" i18nKey="settings.allowedorigins.description">
              <Embed>1</Embed>
              <Embed>2</Embed>
            </Trans>
          </Small>
          {(settings.enabled !== undefined && (
            <Input
              placeholder="*"
              value={settings.allowed_origins}
              onInput={(e) => dispatchSettings(['set_allowedorigins', e.currentTarget.value])}
            />
          )) || <Loader width="100%" height="2em" />}
        </div>
        <div>
          <h3>{t('settings.token.heading')}</h3>
          <Small>
            <Trans ns="routes.guildsettings.api" i18nKey="settings.token.description">
              <Embed>1</Embed>
              <Embed>2</Embed>
              <Embed>3</Embed>
            </Trans>
          </Small>
          {(settings.enabled !== undefined && (
            <>
              {(settings.protected && (
                <StyledHint color="green">
                  <ShieldIcon />
                  <span>{t('settings.token.protected')}</span>
                </StyledHint>
              )) || (
                <StyledHint color="orange">
                  <WarnIcon />
                  <span>{t('settings.token.unprotected')}</span>
                </StyledHint>
              )}
              <Input
                type="password"
                placeholder={t(
                  `settings.token.placeholder.${settings.protected ? 'protected' : 'unprotected'}`,
                )}
                value={settings.token}
                onInput={(e) => dispatchSettings(['set_token', e.currentTarget.value])}
              />
            </>
          )) || (
            <>
              <Loader width="100%" height="2em" margin="0 0 1em 0" />
              <Loader width="100%" height="2em" />
            </>
          )}
        </div>
      </Settings>

      <Controls>
        <Button variant="green" onClick={_saveSettings}>
          {t('save')}
        </Button>
        <Button disabled={!settings.protected} variant="orange" onClick={_resetToken}>
          {t('reset')}
        </Button>
      </Controls>
    </MaxWidthContainer>
  );
};

export default APIRoute;
