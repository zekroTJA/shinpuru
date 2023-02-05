import { useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate, useParams } from 'react-router';
import styled from 'styled-components';
import { ReactComponent as AntiraidIcon } from '../../assets/antiraid.svg';
import { ReactComponent as APIIcon } from '../../assets/api.svg';
import { ReactComponent as BackupIcon } from '../../assets/backup.svg';
import { ReactComponent as CodeIcon } from '../../assets/code.svg';
import { ReactComponent as DataIcon } from '../../assets/data.svg';
import { ReactComponent as HammerIcon } from '../../assets/hammer.svg';
import { ReactComponent as KarmaIcon } from '../../assets/karma.svg';
import { ReactComponent as PermissionsIcon } from '../../assets/lock-open.svg';
import { ReactComponent as LogsIcon } from '../../assets/logs.svg';
import { ReactComponent as SettingsIcon } from '../../assets/settings.svg';
import { ReactComponent as SPBrand } from '../../assets/sp-brand.svg';
import SPIcon from '../../assets/sp-icon.png';
import { ReactComponent as TriangleIcon } from '../../assets/triangle.svg';
import { ReactComponent as UsersIcon } from '../../assets/users.svg';
import { ReactComponent as VerificationIcon } from '../../assets/verification.svg';
import { useApi } from '../../hooks/useApi';
import { useGuilds } from '../../hooks/useGuilds';
import { usePerms } from '../../hooks/usePerms';
import { useSelfUser } from '../../hooks/useSelfUser';
import { Guild } from '../../lib/shinpuru-ts/src';
import { useStore } from '../../services/store';
import { Button } from '../Button';
import { DiscordImage } from '../DiscordImage';
import { Flex } from '../Flex';
import { GuildSelect } from '../GuildSelect';
import { Heading } from '../Heading';
import { Hoverplate } from '../Hoverplate';
import { Loader } from '../Loader';
import { Entry } from './Entry';
import { Section } from './Section';

type Props = {};

const Brand = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
  table-layout: fixed;

  > img {
    width: 38px;
    height: 38px;
  }

  > svg {
    width: 100%;
    height: 38px;
    justify-content: flex-start;
  }
`;

const EntryContainer = styled.div`
  display: flex;
  flex-direction: column;
  margin-top: 1em;
  overflow-y: auto;
  height: 100%;
  gap: 1em;
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

const StyledHoverplate = styled(Hoverplate)`
  margin-top: auto;
`;

const StyledNav = styled.nav`
  display: flex;
  flex-direction: column;
  background-color: ${(p) => p.theme.background2};
  margin: 1rem;
  padding: 1rem;
  border-radius: 12px;
  width: 30vw;
  max-width: 15rem;

  @media (orientation: portrait) {
    width: fit-content;

    ${Brand} > svg {
      display: none;
    }

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

    ${Heading} {
      display: none;
    }
  }
`;

export const Navbar: React.FC<Props> = () => {
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

  return (
    <StyledNav>
      <Brand>
        <img src={SPIcon} alt="shinpuru Heading" />
        <SPBrand />
      </Brand>

      <StyledGuildSelect
        guilds={guilds ?? []}
        value={selectedGuild}
        onElementSelect={_onGuildSelect}
      />

      <EntryContainer>
        <Section title={t('section.guilds.title')}>
          <StyledEntry path={`guilds/${selectedGuild?.id}/members`}>
            <UsersIcon />
            <span>{t('section.guilds.members')}</span>
          </StyledEntry>
          <StyledEntry path={`guilds/${selectedGuild?.id}/modlog`}>
            <HammerIcon />
            <span>{t('section.guilds.modlog')}</span>
          </StyledEntry>
        </Section>

        {isAllowed('sp.guild.config.*') && (
          <Section title={t('section.guildsettings.title')}>
            <StyledEntry path={`guilds/${selectedGuild?.id}/settings/general`}>
              <SettingsIcon />
              <span>{t('section.guildsettings.general')}</span>
            </StyledEntry>
            {isAllowed('sp.guild.config.perms') && (
              <StyledEntry path={`guilds/${selectedGuild?.id}/settings/permissions`}>
                <PermissionsIcon />
                <span>{t('section.guildsettings.permissions')}</span>
              </StyledEntry>
            )}
            {isAllowed('sp.guild.admin.backup') && (
              <StyledEntry path={`guilds/${selectedGuild?.id}/settings/backups`}>
                <BackupIcon />
                <span>{t('section.guildsettings.backup')}</span>
              </StyledEntry>
            )}
            {isAllowed('sp.guild.config.antiraid') && (
              <StyledEntry path={`guilds/${selectedGuild?.id}/settings/antiraid`}>
                <AntiraidIcon />
                <span>{t('section.guildsettings.antiraid')}</span>
              </StyledEntry>
            )}
            {isAllowed('sp.guild.config.exec') && (
              <StyledEntry path={`guilds/${selectedGuild?.id}/settings/codeexec`}>
                <CodeIcon />
                <span>{t('section.guildsettings.codeexec')}</span>
              </StyledEntry>
            )}
            {isAllowed('sp.guild.config.verification') && (
              <StyledEntry path={`guilds/${selectedGuild?.id}/settings/verification`}>
                <VerificationIcon />
                <span>{t('section.guildsettings.verification')}</span>
              </StyledEntry>
            )}
            {isAllowed('sp.guild.config.karma') && (
              <StyledEntry path={`guilds/${selectedGuild?.id}/settings/karma`}>
                <KarmaIcon />
                <span>{t('section.guildsettings.karma')}</span>
              </StyledEntry>
            )}
            {isAllowed('sp.guild.config.logs') && (
              <StyledEntry path={`guilds/${selectedGuild?.id}/settings/logs`}>
                <LogsIcon />
                <span>{t('section.guildsettings.logs')}</span>
              </StyledEntry>
            )}
            {isAllowed('sp.guild.admin.flushdata') && (
              <StyledEntry path={`guilds/${selectedGuild?.id}/settings/data`}>
                <DataIcon />
                <span>{t('section.guildsettings.data')}</span>
              </StyledEntry>
            )}
            {isAllowed('sp.guild.config.api') && (
              <StyledEntry path={`guilds/${selectedGuild?.id}/settings/api`}>
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
            <Button>{t('user-settings')}</Button>
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
    </StyledNav>
  );
};
