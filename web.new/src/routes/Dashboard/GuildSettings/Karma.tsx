import emojis from 'emoji.json';
import React, { useEffect, useReducer, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useParams } from 'react-router';
import { uid } from 'react-uid';
import styled from 'styled-components';
import { Button } from '../../../components/Button';
import { Controls } from '../../../components/Controls';
import { Input } from '../../../components/Input';
import { KarmaRuleEntry, KarmaRuleInput } from '../../../components/KarmaRule';
import { Loader } from '../../../components/Loader';
import { MaxWidthContainer } from '../../../components/MaxWidthContainer';

import { Small } from '../../../components/Small';
import { SplitContainer } from '../../../components/SplitContainer';
import { Switch } from '../../../components/Switch';
import { TagElement, TagsInput } from '../../../components/TagsInput';
import { useApi } from '../../../hooks/useApi';
import { useGuild } from '../../../hooks/useGuild';
import { useNotifications } from '../../../hooks/useNotifications';
import { KarmaRule, KarmaSettings, Member } from '../../../lib/shinpuru-ts/src';
import { listReducer } from '../../../util/reducer';

type Props = {};

const settingsReducer = (
  state: Partial<KarmaSettings>,
  [type, payload]:
    | ['set_state', Partial<KarmaSettings>]
    | ['set_enabled' | 'set_penalty', boolean]
    | ['set_increase' | 'set_decrease', string[]]
    | ['set_tokens', number],
) => {
  switch (type) {
    case 'set_state':
      return { ...state, ...payload };
    case 'set_enabled':
      return { ...state, state: payload };
    case 'set_increase':
      return { ...state, emotes_increase: payload };
    case 'set_decrease':
      return { ...state, emotes_decrease: payload };
    case 'set_tokens':
      return { ...state, tokens: payload > 0 ? payload : 1 };
    case 'set_penalty':
      return { ...state, penalty: payload };
    default:
      return state;
  }
};

const rulesReducer = listReducer<KarmaRule>;

const blocklistReducer = (
  state: Member[],
  [type, payload]: ['set', Member[]] | ['add', Member | Member[]] | ['remove', Member],
) => {
  switch (type) {
    case 'set':
      return payload;
    case 'add':
      return [...state, ...(Array.isArray(payload) ? payload : [payload])];
    case 'remove':
      return state.filter((e) => e.user.id !== payload.user.id);
    default:
      return state;
  }
};

const InputContainer = styled.section`
  > input,
  label {
    display: block;
    width: 100%;
  }

  > label {
    margin-bottom: 1em;
  }
`;

const BlocklistInputContainer = styled.div`
  display: flex;
  gap: 1em;
  white-space: nowrap;

  > ${Input} {
    width: 100%;
  }
`;

const StyledTable = styled.table`
  width: 100%;
  margin-top: 2em;

  th {
    text-align: start;
  }

  ${Button} {
    padding: 0.5em;
  }
`;

const KarmaRoute: React.FC<Props> = ({}) => {
  const { t } = useTranslation('routes.guildsettings.karma');
  const { pushNotification } = useNotifications();
  const { guildid } = useParams();
  const guild = useGuild(guildid);
  const fetch = useApi();
  const [settings, dispatchSettings] = useReducer(settingsReducer, {} as Partial<KarmaSettings>);
  const [rules, dispatchRules] = useReducer(rulesReducer, []);
  const [blocklist, dispatchBlocklist] = useReducer(blocklistReducer, []);
  const [blocklistInput, setBlocklistInput] = useState('');

  const _saveSettings = () => {
    if (!guildid) return;

    fetch((c) => c.guilds.settings(guildid).setKarma(settings as KarmaSettings))
      .then(() =>
        pushNotification({
          message: t('notifications.saved'),
          type: 'SUCCESS',
        }),
      )
      .catch();
  };

  const _addKarmaRule = (r: KarmaRule) => {
    if (!guildid) return;
    fetch((c) => c.guilds.settings(guildid).addKarmaRule(r))
      .then((res) => {
        dispatchRules(['add', res]);
        pushNotification({
          message: t('notifications.ruleapplied'),
          type: 'SUCCESS',
        });
      })
      .catch();
  };

  const _removeKarmaRule = (r: KarmaRule) => {
    if (!guildid) return;
    fetch((c) => c.guilds.settings(guildid).removeKarmaRule(r.id)).then((res) => {
      dispatchRules(['remove', r]);
      pushNotification({
        message: t('notifications.ruleremoved'),
        type: 'SUCCESS',
      });
    });
  };

  const _addBlocklist = (v: string) => {
    if (!guildid) return;
    fetch((c) => c.guilds.settings(guildid).addKarmaBlocklist(v))
      .then((res) => {
        dispatchBlocklist(['add', res]);
        setBlocklistInput('');
        pushNotification({
          message: t('notifications.blocklistadded'),
          type: 'SUCCESS',
        });
      })
      .catch();
  };

  const _removeBlocklist = (v: Member) => {
    if (!guildid) return;
    fetch((c) => c.guilds.settings(guildid).removeKarmaBlocklist(v.user.id))
      .then((res) => {
        dispatchBlocklist(['remove', v]);
        pushNotification({
          message: t('notifications.blocklistremoved'),
          type: 'SUCCESS',
        });
      })
      .catch();
  };

  useEffect(() => {
    if (!guildid) return;

    fetch((c) => c.guilds.settings(guildid).karma())
      .then((res) => dispatchSettings(['set_state', res]))
      .catch();

    fetch((c) => c.guilds.settings(guildid).karmaRules())
      .then((res) => dispatchRules(['set', res.data]))
      .catch();

    fetch((c) => c.guilds.settings(guildid).karmaBlocklist())
      .then((res) => dispatchBlocklist(['set', res.data]))
      .catch();
  }, [guildid]);

  const emojiOptions: TagElement<string>[] = emojis.map((e) => ({
    id: e.codes,
    value: e.char,
    display: e.char,
    keywords: e.name.split(' '),
  }));

  return (
    <MaxWidthContainer>
      <h1>{t('heading')}</h1>
      <Small>{t('explanation')}</Small>

      <h2>{t('settings')}</h2>
      {(settings.state !== undefined && (
        <Switch
          enabled={settings.state}
          onChange={(e) => dispatchSettings(['set_enabled', e])}
          labelAfter={t('enable')}
        />
      )) || <Loader width="20em" height="2em" />}

      <h3>{t('emotes.heading')}</h3>
      <Small>{t('emotes.description')}</Small>
      <SplitContainer>
        <InputContainer>
          <label>{t('emotes.increase')}</label>
          {(settings.state !== undefined && (
            <TagsInput
              selected={emojiOptions.filter((e) => settings.emotes_increase?.includes(e.value))}
              options={emojiOptions}
              onChange={(e) => dispatchSettings(['set_increase', e.map((e) => e.value)])}
            />
          )) || <Loader width="100%" height="3em" />}
        </InputContainer>
        <InputContainer>
          <label>{t('emotes.decrease')}</label>
          {(settings.state !== undefined && (
            <TagsInput
              selected={emojiOptions.filter((e) => settings.emotes_decrease?.includes(e.value))}
              options={emojiOptions}
              onChange={(e) => dispatchSettings(['set_decrease', e.map((e) => e.value)])}
            />
          )) || <Loader width="100%" height="3em" />}
        </InputContainer>
      </SplitContainer>

      <h3>{t('limit.heading')}</h3>
      <Small>{t('limit.description')}</Small>
      <InputContainer>
        {(settings.state !== undefined && (
          <Input
            type="number"
            min="1"
            value={settings.tokens}
            onInput={(e) => dispatchSettings(['set_tokens', parseInt(e.currentTarget.value)])}
          />
        )) || <Loader width="100%" height="2em" />}
      </InputContainer>

      <h3>{t('penalty.heading')}</h3>
      <Small>{t('penalty.description')}</Small>
      <InputContainer>
        {(settings.state !== undefined && (
          <Switch
            enabled={settings.penalty!}
            onChange={(e) => dispatchSettings(['set_penalty', e])}
            labelAfter={t('penalty.switch')}
          />
        )) || <Loader width="20em" height="2em" />}
      </InputContainer>

      <Controls>
        <Button variant="green" onClick={_saveSettings}>
          {t('save')}
        </Button>
      </Controls>

      <h2>{t('rules.heading')}</h2>
      <Small>{t('rules.description')}</Small>
      {(guild !== undefined && (
        <div>
          <KarmaRuleInput guild={guild!} onApply={(r) => _addKarmaRule(r)} />
          {rules.map((r) => (
            <KarmaRuleEntry
              key={uid(r)}
              guild={guild!}
              rule={r}
              onRemove={() => _removeKarmaRule(r)}
            />
          ))}
        </div>
      )) || <Loader width="20em" height="2em" />}

      <h2>{t('blocklist.heading')}</h2>
      <Small>{t('blocklist.description')}</Small>
      <BlocklistInputContainer>
        <Input value={blocklistInput} onChange={(e) => setBlocklistInput(e.currentTarget.value)} />
        <Button disabled={!blocklistInput} onClick={() => _addBlocklist(blocklistInput)}>
          {t('blocklist.addmember')}
        </Button>
      </BlocklistInputContainer>
      <StyledTable>
        <tbody>
          <tr>
            <th>{t('blocklist.table.id')}</th>
            <th>{t('blocklist.table.name')}</th>
            <th>{t('blocklist.table.nick')}</th>
            <th>{t('blocklist.table.unblock')}</th>
          </tr>
          {blocklist.map((m) => (
            <tr>
              <td>{m.user.id}</td>
              <td>
                {m.user.username}#{m.user.discriminator}
              </td>
              <td>{m.nick || m.user.username}</td>
              <td>
                <Button onClick={() => _removeBlocklist(m)}>{t('blocklist.table.unblock')}</Button>
              </td>
            </tr>
          ))}
        </tbody>
      </StyledTable>
    </MaxWidthContainer>
  );
};

export default KarmaRoute;
