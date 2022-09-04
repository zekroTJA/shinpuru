import React, { useEffect, useReducer } from 'react';
import { Trans, useTranslation } from 'react-i18next';
import { useParams } from 'react-router';
import styled from 'styled-components';
import { Button } from '../../../components/Button';
import { Controls } from '../../../components/Controls';
import { Embed } from '../../../components/Embed';
import { Flex } from '../../../components/Flex';
import { Input } from '../../../components/Input';
import { Loader } from '../../../components/Loader';
import { MaxWidthContainer } from '../../../components/MaxWidthContainer';
import { NotificationType } from '../../../components/Notifications';
import { RoleInput } from '../../../components/RoleInput';
import { Element, Select } from '../../../components/Select';
import { Small } from '../../../components/Small';
import { TagElement } from '../../../components/TagsInput/TagsInput';
import { useApi } from '../../../hooks/useApi';
import { useGuild } from '../../../hooks/useGuild';
import { useNotifications } from '../../../hooks/useNotifications';
import { usePerms } from '../../../hooks/usePerms';
import { Channel, ChannelType, GuildSettings, Role } from '../../../lib/shinpuru-ts/src';

type Props = {};

type GuildSettingsVM = {
  autoroles: Role[];
  modlogchannel?: Element<Channel>;
  voicelogchannel?: Element<Channel>;
  joinmessagechannel?: Element<Channel>;
  joinmessagetext?: string;
  leavemessagechannel?: Element<Channel>;
  leavemessagetext?: string;
};

const guildSettingsReducer = (
  state: GuildSettingsVM,
  [type, payload]:
    | ['set_state', Partial<GuildSettingsVM>]
    | ['set_autoroles', Role[]]
    | [
        (
          | 'set_modlogchannel'
          | 'set_voicelogchannel'
          | 'set_joinmessagechannel'
          | 'set_leavemessagechannel'
        ),
        Element<Channel>,
      ]
    | ['set_joinmessagetext' | 'set_leavemessagetext', string]
    | [
        | 'reset_modlogchannel'
        | 'reset_voicelogchannel'
        | 'reset_joinmessagechannel'
        | 'reset_leavemessagechannel'
        | 'reset_joinmessagetext'
        | 'reset_leavemessagetext',
      ],
) => {
  console.log(type, payload);
  switch (type) {
    case 'set_state':
      return { ...state, ...payload };
    case 'set_autoroles':
      return { ...state, autoroles: payload };
    case 'set_modlogchannel':
      return { ...state, modlogchannel: payload };
    case 'set_voicelogchannel':
      return { ...state, voicelogchannel: payload };
    case 'set_joinmessagechannel':
      return { ...state, joinmessagechannel: payload };
    case 'set_joinmessagetext':
      return { ...state, joinmessagetext: payload };
    case 'set_leavemessagechannel':
      return { ...state, leavemessagechannel: payload };
    case 'set_leavemessagetext':
      return { ...state, leavemessagetext: payload };
    case 'reset_modlogchannel':
      return { ...state, modlogchannel: undefined };
    case 'reset_voicelogchannel':
      return { ...state, voicelogchannel: undefined };
    case 'reset_joinmessagechannel':
      return { ...state, joinmessagechannel: undefined };
    case 'reset_leavemessagechannel':
      return { ...state, leavemessagechannel: undefined };
    case 'reset_joinmessagetext':
      return { ...state, joinmessagetext: '' };
    case 'reset_leavemessagetext':
      return { ...state, leavemessagetext: '' };
    default:
      return state;
  }
};

const Container = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.5em;
`;

const Section = styled.section`
  width: 100%;

  > label {
    margin: 1em 0 0.5em 0;
    display: block;
  }

  > div {
    display: flex;
    gap: 1em;
    width: 100%;

    > div,
    ${Input} {
      width: 100%;
    }

    > ${Button} {
      padding-bottom: 0;
      padding-top: 0;
    }
  }
`;

const GeneralRoute: React.FC<Props> = () => {
  const { t } = useTranslation('routes.guildsettings.general');
  const { guildid } = useParams();
  const fetch = useApi();
  const guild = useGuild(guildid);
  const { allowedPerms, isAllowed } = usePerms(guildid);
  const [settings, dispatchSettings] = useReducer(guildSettingsReducer, {} as GuildSettingsVM);
  const { pushNotification } = useNotifications();

  const _saveSettings = () => {
    if (!guildid) return;

    const gs = {} as GuildSettings;

    if (isAllowed('sp.guild.config.autorole')) gs.autoroles = settings.autoroles.map((r) => r.id);

    if (isAllowed('sp.guild.config.modlog'))
      gs.modlogchannel = settings.modlogchannel?.value.id ?? '__RESET__';

    if (isAllowed('sp.guild.config.voicelog'))
      gs.voicelogchannel = settings.voicelogchannel?.value.id ?? '__RESET__';

    if (isAllowed('sp.guild.config.announcements')) {
      gs.joinmessagechannel = settings.joinmessagechannel?.value.id ?? '__RESET__';
      gs.joinmessagetext = settings.joinmessagetext || '__RESET__';
      gs.leavemessagechannel = settings.leavemessagechannel?.value.id ?? '__RESET__';
      gs.leavemessagetext = settings.leavemessagetext || '__RESET__';
    }

    fetch((c) => c.guilds.settings(guildid).setSettings(gs))
      .then(() =>
        pushNotification({
          message: t('notifications.saved'),
          type: NotificationType.SUCCESS,
        }),
      )
      .catch();
  };

  useEffect(() => {
    if (!guild) return;
    fetch((c) => c.guilds.settings(guild.id).settings())
      .then((res) => {
        dispatchSettings([
          'set_autoroles',
          (res.autoroles ?? [])
            .map((rid) => guild.roles!.find((r) => r.id === rid))
            .filter((r) => !!r) as Role[],
        ]);

        const modlogchannel = guild.channels!.find((c) => c.id === res.modlogchannel)!;
        if (modlogchannel)
          dispatchSettings([
            'set_modlogchannel',
            {
              id: modlogchannel.id,
              value: modlogchannel,
              display: modlogchannel.name,
            },
          ]);

        const voicelogchannel = guild.channels!.find((c) => c.id === res.voicelogchannel)!;
        if (voicelogchannel)
          dispatchSettings([
            'set_voicelogchannel',
            {
              id: voicelogchannel.id,
              value: voicelogchannel,
              display: voicelogchannel.name,
            },
          ]);

        const joinmessagechannel = guild.channels!.find((c) => c.id === res.joinmessagechannel)!;
        if (joinmessagechannel)
          dispatchSettings([
            'set_joinmessagechannel',
            {
              id: joinmessagechannel.id,
              value: joinmessagechannel,
              display: joinmessagechannel.name,
            },
          ]);

        const leavemessagechannel = guild.channels!.find((c) => c.id === res.leavemessagechannel)!;
        if (leavemessagechannel)
          dispatchSettings([
            'set_leavemessagechannel',
            {
              id: leavemessagechannel.id,
              value: leavemessagechannel,
              display: leavemessagechannel.name,
            },
          ]);

        dispatchSettings([
          'set_state',
          {
            joinmessagetext: res.joinmessagetext,
            leavemessagetext: res.leavemessagetext,
          },
        ]);
      })
      .catch();
  }, [guild]);

  const textChannelOptions = guild?.channels
    ?.filter((c) => c.type === ChannelType.GUILD_TEXT)
    .map(
      (c) =>
        ({
          id: c.id,
          value: c,
          display: c.name,
        } as Element<Channel>),
    );

  const roleTagOptions =
    guild?.roles?.map(
      (r) =>
        ({ id: r.id, display: r.name, keywords: [r.id, r.name], value: r } as TagElement<Role>),
    ) ?? [];

  return (
    <MaxWidthContainer>
      <h1>{t('heading')}</h1>
      {(guild && allowedPerms && (
        <Container>
          {isAllowed('sp.guild.config.autorole') && (
            <Section>
              <h2>{t('autoroles.title')}</h2>
              <RoleInput
                guild={guild}
                selected={settings.autoroles}
                onChange={(v) => dispatchSettings(['set_autoroles', v])}
                placeholder={t('autoroles.placeholder')}
              />
            </Section>
          )}

          {isAllowed('sp.guild.config.modlog') && (
            <Section>
              <h2>{t('modlog.title')}</h2>
              <label>{t('modlog.channel_label')}</label>
              <div>
                <Select
                  options={textChannelOptions!}
                  value={settings.modlogchannel}
                  onElementSelect={(e) => dispatchSettings(['set_modlogchannel', e])}
                  placeholder={t('modlog.channel_placeholder')}
                />
                <Button onClick={() => dispatchSettings(['reset_modlogchannel'])}>
                  {t('reset')}
                </Button>
              </div>
            </Section>
          )}

          {isAllowed('sp.guild.config.voicelog') && (
            <Section>
              <h2>{t('voicelog.title')}</h2>
              <label>{t('voicelog.channel_label')}</label>
              <div>
                <Select
                  options={textChannelOptions!}
                  value={settings.voicelogchannel}
                  onElementSelect={(e) => dispatchSettings(['set_voicelogchannel', e])}
                  placeholder={t('voicelog.channel_placeholder')}
                />
                <Button onClick={() => dispatchSettings(['reset_voicelogchannel'])}>
                  {t('reset')}
                </Button>
              </div>
            </Section>
          )}

          {isAllowed('sp.guild.config.announcements') && (
            <>
              <Section>
                <h2>{t('joinmessage.title')}</h2>
                <label>{t('joinmessage.channel_label')}</label>
                <div>
                  <Select
                    options={textChannelOptions!}
                    value={settings.joinmessagechannel}
                    onElementSelect={(e) => dispatchSettings(['set_joinmessagechannel', e])}
                    placeholder={t('joinmessage.channel_placeholder')}
                  />
                  <Button onClick={() => dispatchSettings(['reset_joinmessagechannel'])}>
                    {t('reset')}
                  </Button>
                </div>
                <label>{t('joinmessage.message_label')}</label>
                <Small>
                  <Trans ns="routes.guildsettings.general" i18nKey="joinmessage.message_hint">
                    <Embed>[user]</Embed>
                    <Embed>[ment]</Embed>
                  </Trans>
                </Small>
                <div>
                  <Input
                    value={settings.joinmessagetext}
                    placeholder={t('joinmessage.message_placeholder')}
                    onInput={(e) =>
                      dispatchSettings(['set_joinmessagetext', e.currentTarget.value])
                    }
                  />
                  <Button onClick={() => dispatchSettings(['reset_joinmessagechannel'])}>
                    {t('reset')}
                  </Button>
                </div>
              </Section>

              <Section>
                <h2>{t('leavemessage.title')}</h2>
                <label>{t('leavemessage.channel_label')}</label>
                <div>
                  <Select
                    options={textChannelOptions!}
                    value={settings.leavemessagechannel}
                    onElementSelect={(e) => dispatchSettings(['set_leavemessagechannel', e])}
                    placeholder={t('leavemessage.channel_placeholder')}
                  />
                  <Button onClick={() => dispatchSettings(['reset_leavemessagetext'])}>
                    {t('reset')}
                  </Button>
                </div>
                <label>{t('leavemessage.message_label')}</label>
                <Small>
                  <Trans ns="routes.guildsettings.general" i18nKey="leavemessage.message_hint">
                    <Embed>[user]</Embed>
                    <Embed>[ment]</Embed>
                  </Trans>
                </Small>
                <div>
                  <Input
                    value={settings.leavemessagetext}
                    placeholder={t('leavemessage.message_placeholder')}
                    onInput={(e) =>
                      dispatchSettings(['set_leavemessagetext', e.currentTarget.value])
                    }
                  />
                  <Button onClick={() => dispatchSettings(['reset_leavemessagetext'])}>
                    {t('reset')}
                  </Button>
                </div>
              </Section>
            </>
          )}

          <Controls>
            <Button variant="green" onClick={_saveSettings}>
              {t('save')}
            </Button>
          </Controls>
        </Container>
      )) || (
        <>
          <Loader width="30ch" height="2em" margin="1em 0 0 0" />
          <Flex gap="1em">
            <Loader width="100%" height="2em" margin="1em 0 0 0" />
            <Loader width="10ch" height="2em" margin="1em 0 0 0" />
          </Flex>
          <Loader width="30ch" height="2em" margin="2em 0 0 0" />
          <Flex gap="1em">
            <Loader width="100%" height="2em" margin="1em 0 0 0" />
            <Loader width="10ch" height="2em" margin="1em 0 0 0" />
          </Flex>
          <Flex gap="1em">
            <Loader width="100%" height="2em" margin="1em 0 0 0" />
            <Loader width="10ch" height="2em" margin="1em 0 0 0" />
          </Flex>
          <Loader width="30ch" height="2em" margin="2em 0 0 0" />
          <Flex gap="1em">
            <Loader width="100%" height="2em" margin="1em 0 0 0" />
            <Loader width="10ch" height="2em" margin="1em 0 0 0" />
          </Flex>
          <Flex gap="1em">
            <Loader width="100%" height="2em" margin="1em 0 0 0" />
            <Loader width="10ch" height="2em" margin="1em 0 0 0" />
          </Flex>
        </>
      )}
    </MaxWidthContainer>
  );
};

export default GeneralRoute;
