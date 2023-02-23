import { AntiraidActionType, AntiraidSettings, JoinlogEntry } from '../../../lib/shinpuru-ts/src';
import React, { useEffect, useReducer, useState } from 'react';
import { Trans, useTranslation } from 'react-i18next';

import { ReactComponent as BanIcon } from '../../../assets/ban.svg';
import { Button } from '../../../components/Button';
import { Controls } from '../../../components/Controls';
import { ReactComponent as DeleteIcon } from '../../../assets/delete.svg';
import { ReactComponent as DownloadIcon } from '../../../assets/download.svg';
import { Flex } from '../../../components/Flex';
import { Input } from '../../../components/Input';
import { JoinLogEntry } from '../../../components/JoinLogEntry';
import { ReactComponent as KickIcon } from '../../../assets/kick.svg';
import { Loader } from '../../../components/Loader';
import { MaxWidthContainer } from '../../../components/MaxWidthContainer';
import { ReactComponent as RefreshIcon } from '../../../assets/refresh.svg';
import { Small } from '../../../components/Small';
import { Switch } from '../../../components/Switch';
import styled from 'styled-components';
import { uid } from 'react-uid';
import { useApi } from '../../../hooks/useApi';
import { useNotifications } from '../../../hooks/useNotifications';
import { useParams } from 'react-router';

type Props = {};

const settingsReducer = (
  state: Partial<AntiraidSettings>,
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

const JoinlogContainer = styled.div`
  margin-top: 2.5em;
`;

const JoinlogEntries = styled.div`
  margin-top: 1em;
  display: flex;
  flex-direction: column;
  gap: 0.8em;

  table {
    width: 100%;
    border-collapse: separate;
    border-spacing: 0 0.8em;

    th {
      text-align: start;
      padding: 0.8em;
      cursor: pointer;
      > span {
        text-transform: uppercase;
        opacity: 0.7;
        font-size: 0.9em;
      }
    }
  }
`;

const ActionButton = styled(Button)`
  @keyframes action_anim {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(365deg);
    }
  }

  > svg {
  }
  &:enabled:focus {
    > svg {
      animation: action_anim 1s;
    }
  }
`;

const AntiraidRoute: React.FC<Props> = ({}) => {
  const { t } = useTranslation('routes.guildsettings.antiraid');
  const { pushNotification } = useNotifications();
  const { guildid } = useParams();
  const fetch = useApi();
  const [settings, dispatchSettings] = useReducer(settingsReducer, {} as Partial<AntiraidSettings>);
  const [entries, setEntries] = useState<JoinlogEntry[]>();
  const [selected, setSelected] = useState<string[]>([]);

  const _saveSettings = () => {
    if (!guildid) return;

    fetch((c) => c.guilds.settings(guildid).setAntiraid(settings as AntiraidSettings))
      .then(() =>
        pushNotification({
          message: t('notifications.saved'),
          type: 'SUCCESS',
        }),
      )
      .catch();
  };

  const _onCheck = (checked: boolean, e: JoinlogEntry) => {
    if (checked) setSelected([e.user_id, ...selected]);
    else setSelected(selected.filter((c) => c !== e.user_id));
  };

  const _onCheckAll = (checked: boolean) => {
    if (checked) setSelected(entries?.map((e) => e.user_id) ?? []);
    else setSelected([]);
  };

  const _download = () => {
    const element = document.createElement('a');
    element.setAttribute(
      'href',
      `data:application/json;charset=utf-8,${encodeURIComponent(
        JSON.stringify(entries, undefined, 2),
      )}`,
    );
    element.setAttribute('download', 'joinlog_export');
    element.click();
  };

  const _refresh = () => {
    if (!guildid) return;
    fetch((c) => c.guilds.antiraidJoinlog(guildid))
      .then((r) => setEntries(r.data))
      .catch();
  };

  const _clear = () => {
    if (!guildid) return;
    fetch((c) => c.guilds.deleteAntiraidJoinlog(guildid))
      .then(() => {
        setEntries([]);
        setSelected([]);
      })
      .catch();
  };

  const _action = (type: AntiraidActionType) => () => {
    if (!guildid || !entries) return;
    fetch((c) => c.guilds.settings(guildid).addAntiraidAction({ ids: selected, type }))
      .then(() => {
        pushNotification({
          message: (
            <Trans
              ns="routes.guildsettings.antiraid"
              i18nKey={`notifications.${type === AntiraidActionType.KICK ? 'kicked' : 'banned'}`}
              values={{ count: selected.length }}
            />
          ),
          type: 'SUCCESS',
        });
        setEntries(entries.filter((e) => !selected.includes(e.user_id)));
        setSelected([]);
      })
      .catch();
  };

  useEffect(() => {
    if (!guildid) return;

    fetch((c) => c.guilds.settings(guildid).antiraid())
      .then((res) => dispatchSettings(['set_state', res]))
      .catch();

    _refresh();
  }, [guildid]);

  return (
    <MaxWidthContainer>
      <h1>{t('heading')}</h1>
      <Small>{t('explanation')}</Small>

      <h2>Settings</h2>
      {(settings.state !== undefined && (
        <Switch
          enabled={settings.state}
          onChange={(e) => dispatchSettings(['set_enabled', e])}
          labelAfter={t('enable')}
        />
      )) || <Loader width="20em" height="2em" />}

      <InputContainer>
        <label>{t('regeneration.label')}</label>
        <Small>{t('regeneration.explanation')}</Small>
        {(settings.regeneration_period !== undefined && (
          <Input
            type="number"
            min="1"
            value={settings.regeneration_period}
            onInput={(e) => dispatchSettings(['set_regeneration', parseInt(e.currentTarget.value)])}
          />
        )) || <Loader width="100%" height="2em" />}
      </InputContainer>

      <InputContainer>
        <label>{t('burst.label')}</label>
        <Small>{t('burst.explanation')}</Small>
        {(settings.burst !== undefined && (
          <Input
            type="number"
            min="1"
            value={settings.burst}
            onInput={(e) => dispatchSettings(['set_burst', parseInt(e.currentTarget.value)])}
          />
        )) || <Loader width="100%" height="2em" />}
      </InputContainer>

      <InputContainer>
        <label>{t('verification.label')}</label>
        <Small>{t('verification.explanation')}</Small>
        {(settings.verification !== undefined && (
          <Switch
            enabled={settings.verification}
            onChange={(e) => dispatchSettings(['set_verification', e])}
            labelAfter={t('verification.switch_label')}
          />
        )) || <Loader width="20em" height="2em" />}
      </InputContainer>

      <Controls>
        <Button variant="green" onClick={_saveSettings}>
          {t('save')}
        </Button>
      </Controls>

      <JoinlogContainer>
        <h2>Join Log</h2>
        <Flex gap="1em">
          <Button variant="gray" disabled={!entries?.length} onClick={_download}>
            <DownloadIcon />
            {t('log.controls.download')}
          </Button>
          <ActionButton variant="gray" onClick={_refresh}>
            <RefreshIcon />
            {t('log.controls.refresh')}
          </ActionButton>
          <ActionButton variant="gray" disabled={!entries?.length} onClick={_clear}>
            <DeleteIcon />
            {t('log.controls.clear')}
          </ActionButton>
          <Button
            margin="0 0 0 auto"
            variant="orange"
            disabled={!selected.length}
            onClick={_action(AntiraidActionType.KICK)}>
            <KickIcon />
            {t('log.controls.kick')}
          </Button>
          <Button
            variant="red"
            disabled={!selected.length}
            onClick={_action(AntiraidActionType.BAN)}>
            <BanIcon />
            {t('log.controls.ban')}
          </Button>
        </Flex>

        <JoinlogEntries>
          {(entries === undefined && (
            <>
              <Loader width="100%" height="2em" />
              <Loader width="100%" height="2em" />
              <Loader width="100%" height="2em" />
            </>
          )) ||
            (entries!.length > 0 && (
              <table>
                <tbody>
                  <tr onClick={() => _onCheckAll(!(selected.length === entries?.length))}>
                    <th>
                      <span>{t('log.entries.id')}</span>
                    </th>
                    <th>
                      <span>{t('log.entries.tag')}</span>
                    </th>
                    <th>
                      <span>{t('log.entries.created')}</span>
                    </th>
                    <th>
                      <span>{t('log.entries.joined')}</span>
                    </th>
                    <th>
                      <input
                        type="checkbox"
                        checked={selected.length === entries?.length}
                        onChange={(v) => _onCheckAll(v.currentTarget.checked)}
                      />
                    </th>
                  </tr>
                  {entries!.map((e) => (
                    <JoinLogEntry
                      key={uid(e)}
                      entry={e}
                      selected={selected.includes(e.user_id)}
                      onCheck={_onCheck}
                    />
                  ))}
                </tbody>
              </table>
            )) || <Small textAlign="center">{t('log.empty')}</Small>}
        </JoinlogEntries>
      </JoinlogContainer>
    </MaxWidthContainer>
  );
};

export default AntiraidRoute;
