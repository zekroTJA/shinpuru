import { useNavigate, useParams } from 'react-router';

import { ReactComponent as APIIcon } from '../../assets/api.svg';
import { ReactComponent as AntiraidIcon } from '../../assets/antiraid.svg';
import { ReactComponent as BackupIcon } from '../../assets/backup.svg';
import { ReactComponent as BlockIcon } from '../../assets/block.svg';
import { Button } from '../Button';
import { ReactComponent as CodeIcon } from '../../assets/code.svg';
import { ReactComponent as DataIcon } from '../../assets/data.svg';
import { DiscordImage } from '../DiscordImage';
import { Entry } from './Entry';
import { EntryContainer } from './EntryContainer';
import { Flex } from '../Flex';
import { Guild } from '../../lib/shinpuru-ts/src';
import { GuildSelect } from '../GuildSelect';
import { ReactComponent as HammerIcon } from '../../assets/hammer.svg';
import { Hoverplate } from '../Hoverplate';
import { ReactComponent as KarmaIcon } from '../../assets/karma.svg';
import { Loader } from '../Loader';
import { ReactComponent as LogsIcon } from '../../assets/logs.svg';
import { Navbar } from './Navbar';
import { ReactComponent as PermissionsIcon } from '../../assets/lock-open.svg';
import { Section } from './Section';
import { ReactComponent as SettingsIcon } from '../../assets/settings.svg';
import { ReactComponent as TriangleIcon } from '../../assets/triangle.svg';
import { ReactComponent as UsersIcon } from '../../assets/users.svg';
import { ReactComponent as VerificationIcon } from '../../assets/verification.svg';
import styled from 'styled-components';
import { useApi } from '../../hooks/useApi';
import { useEffect } from 'react';
import { useGuilds } from '../../hooks/useGuilds';
import { usePerms } from '../../hooks/usePerms';
import { useSelfUser } from '../../hooks/useSelfUser';
import { useStore } from '../../services/store';
import { useTranslation } from 'react-i18next';

type Props = {};

const StyledHoverplate = styled(Hoverplate)`
  margin-top: auto;
`;

const StyledEntry = styled(Entry)``;

const StyledGuildSelect = styled(GuildSelect)`
  margin-top: 1em;
`;

const SelfContainer = styled.div`
  display: flex;
  align-items: center;
  gap: 1em;
  background-color: ${(p) => p.theme.background3};
  border-radius: 8px;
  padding: 0.5em;
  margin-top: 1em;

  > img {
    width: 2em;
    height: 2em;
  }

  > span {
    display: flex;
    align-items: center;
    width: 100%;

    > svg {
      width: 0.5em;
      height: 0.5em;
      margin: 0 1em 0 auto;
    }
  }
`;

const StyledNavbar = styled(Navbar)`
  @media (orientation: portrait) {
    ${StyledEntry}, ${SelfContainer} {
      justify-content: center;
      span {
        display: none;
      }
    }

    ${StyledGuildSelect} > div > div {
      justify-content: center;
      > span {
        display: none;
      }
    }
  }
`;

export const NavbarDashboard: React.FC<Props> = ({}) => {
  const { t } = useTranslation('components', { keyPrefix: 'navbar' });
  const nav = useNavigate();
  const fetch = useApi();
  const { guildid } = useParams();
  const guilds = useGuilds();
  const self = useSelfUser();
  const { isAllowed } = usePerms(guildid);
  const [selectedGuild, setSelectedGuild] = useStore((s) => [s.selectedGuild, s.setSelectedGuild]);

  useEffect(() => {
    if (!!guilds && !!guildid) setSelectedGuild(guilds.find((g) => g.id === guildid) ?? guilds[0]);
  }, [guildid, guilds]);

  const _onGuildSelect = (g: Guild) => {
    setSelectedGuild(g);
    nav(`guilds/${g.id}/members`);
  };

  const _logout = () => {
    fetch((c) => c.auth.logout())
      .then(() => nav('/start'))
      .catch();
  };

  const _userSettings = () => {
    nav('/usersettings');
  };

  return (
    <StyledNavbar>
      <StyledGuildSelect
        guilds={guilds ?? []}
        value={selectedGuild}
        onElementSelect={_onGuildSelect}
      />

      <EntryContainer>
        <Section title={t('section.guilds.title')}>
          <StyledEntry path={`/db/guilds/${selectedGuild?.id}/members`}>
            <UsersIcon />
            <span>{t('section.guilds.members')}</span>
          </StyledEntry>
          <StyledEntry path={`/db/guilds/${selectedGuild?.id}/modlog`}>
            <HammerIcon />
            <span>{t('section.guilds.modlog')}</span>
          </StyledEntry>
        </Section>

        {isAllowed('sp.guild.config.*') && (
          <Section title={t('section.guildsettings.title')}>
            <StyledEntry path={`/db/guilds/${selectedGuild?.id}/settings/general`}>
              <SettingsIcon />
              <span>{t('section.guildsettings.general')}</span>
            </StyledEntry>
            {isAllowed('sp.guild.config.perms') && (
              <StyledEntry path={`/db/guilds/${selectedGuild?.id}/settings/permissions`}>
                <PermissionsIcon />
                <span>{t('section.guildsettings.permissions')}</span>
              </StyledEntry>
            )}
            {isAllowed('sp.guild.admin.backup') && (
              <StyledEntry path={`/db/guilds/${selectedGuild?.id}/settings/backups`}>
                <BackupIcon />
                <span>{t('section.guildsettings.backup')}</span>
              </StyledEntry>
            )}
            {isAllowed('sp.guild.config.antiraid') && (
              <StyledEntry path={`/db/guilds/${selectedGuild?.id}/settings/antiraid`}>
                <AntiraidIcon />
                <span>{t('section.guildsettings.antiraid')}</span>
              </StyledEntry>
            )}
            {isAllowed('sp.guild.mod.inviteblock') && (
              <StyledEntry path={`/db/guilds/${selectedGuild?.id}/settings/linkblocking`}>
                <BlockIcon />
                <span>{t('section.guildsettings.linkblocking')}</span>
              </StyledEntry>
            )}
            {isAllowed('sp.guild.config.exec') && (
              <StyledEntry path={`/db/guilds/${selectedGuild?.id}/settings/codeexec`}>
                <CodeIcon />
                <span>{t('section.guildsettings.codeexec')}</span>
              </StyledEntry>
            )}
            {isAllowed('sp.guild.config.verification') && (
              <StyledEntry path={`/db/guilds/${selectedGuild?.id}/settings/verification`}>
                <VerificationIcon />
                <span>{t('section.guildsettings.verification')}</span>
              </StyledEntry>
            )}
            {isAllowed('sp.guild.config.karma') && (
              <StyledEntry path={`/db/guilds/${selectedGuild?.id}/settings/karma`}>
                <KarmaIcon />
                <span>{t('section.guildsettings.karma')}</span>
              </StyledEntry>
            )}
            {isAllowed('sp.guild.config.logs') && (
              <StyledEntry path={`/db/guilds/${selectedGuild?.id}/settings/logs`}>
                <LogsIcon />
                <span>{t('section.guildsettings.logs')}</span>
              </StyledEntry>
            )}
            {isAllowed('sp.guild.admin.flushdata') && (
              <StyledEntry path={`/db/guilds/${selectedGuild?.id}/settings/data`}>
                <DataIcon />
                <span>{t('section.guildsettings.data')}</span>
              </StyledEntry>
            )}
            {isAllowed('sp.guild.config.api') && (
              <StyledEntry path={`/db/guilds/${selectedGuild?.id}/settings/api`}>
                <APIIcon />
                <span>{t('section.guildsettings.api')}</span>
              </StyledEntry>
            )}
          </Section>
        )}
      </EntryContainer>

      <StyledHoverplate
        hoverContent={
          <Flex gap="1em">
            <Button variant="orange" onClick={_logout}>
              {t('logout')}
            </Button>
            <Button onClick={_userSettings}>{t('user-settings')}</Button>
          </Flex>
        }>
        {(self && (
          <SelfContainer>
            <DiscordImage src={self?.avatar_url} round />
            <span>
              {self?.username}
              <TriangleIcon />
            </span>
          </SelfContainer>
        )) || <Loader width="100%" height="2em" borderRadius="8px" />}
      </StyledHoverplate>
    </StyledNavbar>
  );
};
