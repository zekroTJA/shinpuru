import { Element, Select } from '../../components/Select';
import { Presence, Status } from '../../lib/shinpuru-ts/src';
import React, { useEffect, useReducer } from 'react';

import { ActionButton } from '../../components/ActionButton';
import { Controls } from '../../components/Controls';
import { Flex } from '../../components/Flex';
import { Input } from '../../components/Input';
import { MaxWidthContainer } from '../../components/MaxWidthContainer';
import { Small } from '../../components/Small';
import styled from 'styled-components';
import { useApi } from '../../hooks/useApi';
import { useNotifications } from '../../hooks/useNotifications';
import { useTranslation } from 'react-i18next';

type Props = {};

const Section = styled.section`
  > label {
    margin: 0 0 0.5em 0;
    display: block;
  }

  > ${Input} {
    width: 100%;
  }
`;

const presenceReducer = (
  state: Partial<Presence>,
  [type, payload]: ['set_state', Partial<Presence>] | ['set_status', Status] | ['set_game', string],
) => {
  switch (type) {
    case 'set_state':
      return { ...state, ...payload };
    case 'set_status':
      return { ...state, status: payload };
    case 'set_game':
      return { ...state, game: payload };
    default:
      return state;
  }
};

const PresenceRoute: React.FC<Props> = () => {
  const { t } = useTranslation('routes.settings.presence');
  const { pushNotification } = useNotifications();
  const fetch = useApi();
  const [state, dispatchState] = useReducer(presenceReducer, {});

  const _saveSettings = () => {
    return fetch((c) => c.settings.setPresence(state as Presence))
      .then(() => {
        pushNotification({
          message: t('notifications.saved'),
          type: 'SUCCESS',
        });
      })
      .catch();
  };

  useEffect(() => {
    fetch((c) => c.settings.presence())
      .then((r) => dispatchState(['set_state', r]))
      .catch();
  }, []);

  const statusOptions: Element<Status>[] = (['online', 'dnd', 'idle', 'invisible'] as Status[]).map(
    (s) => ({
      id: s,
      display: <span>{t(`status.${s}`)}</span>,
      value: s,
    }),
  );

  return (
    <MaxWidthContainer>
      <h1>{t('heading')}</h1>
      <Small>{t('explanation')}</Small>
      {state.status !== undefined && (
        <Flex direction="column" gap="1em">
          <Section>
            <label>{t('status.label')}</label>
            <Select
              options={statusOptions}
              value={statusOptions.find((s) => s.value === state.status)}
              onElementSelect={(v) => dispatchState(['set_status', v.value])}
            />
          </Section>
          <Section>
            <label>{t('game')}</label>
            <Input
              value={state.game}
              onInput={(e) => dispatchState(['set_game', e.currentTarget.value])}
            />
          </Section>
          <Controls>
            <ActionButton variant="green" onClick={_saveSettings}>
              {t('save')}
            </ActionButton>
          </Controls>
        </Flex>
      )}
    </MaxWidthContainer>
  );
};

export default PresenceRoute;
