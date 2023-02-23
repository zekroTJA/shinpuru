import React, { useEffect, useReducer } from 'react';

import { Button } from '../../../components/Button';
import { Controls } from '../../../components/Controls';
import { Loader } from '../../../components/Loader';
import { MaxWidthContainer } from '../../../components/MaxWidthContainer';
import { Small } from '../../../components/Small';
import { Switch } from '../../../components/Switch';
import { useApi } from '../../../hooks/useApi';
import { useGuild } from '../../../hooks/useGuild';
import { useNotifications } from '../../../hooks/useNotifications';
import { useParams } from 'react-router';
import { useTranslation } from 'react-i18next';

type Props = {};

type State = {
  enabled: boolean;
};

const linkBlockingReducer = (
  state: State,
  [type, payload]: ['set_state', Partial<State>] | ['set_enabled', boolean],
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

const LinkBlockingRoute: React.FC<Props> = ({}) => {
  const { t } = useTranslation('routes.guildsettings.linkblocking');
  const { pushNotification } = useNotifications();
  const { guildid } = useParams();
  const guild = useGuild(guildid);
  const fetch = useApi();
  const [state, dispatchState] = useReducer(linkBlockingReducer, {} as State);

  const _saveSettings = () => {
    if (!guildid) return;
    fetch((c) => c.guilds.setInviteBlock(guildid, state.enabled))
      .then(() =>
        pushNotification({
          message: t('notifications.saved'),
          type: 'SUCCESS',
        }),
      )
      .catch();
  };

  useEffect(() => {
    if (!guild) return;
    dispatchState(['set_enabled', guild.invite_block_enabled]);
  }, [guild]);

  return (
    <MaxWidthContainer>
      <h1>{t('heading')}</h1>
      <Small>{t('explanation')}</Small>
      {(guild && (
        <Switch
          labelAfter={t('toggle')}
          enabled={state.enabled}
          onChange={(v) => dispatchState(['set_enabled', v])}
        />
      )) || <Loader height="4em" />}
      <Controls>
        <Button variant="green" onClick={_saveSettings}>
          {t('save')}
        </Button>
      </Controls>
    </MaxWidthContainer>
  );
};

export default LinkBlockingRoute;
