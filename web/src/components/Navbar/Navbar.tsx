import { useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate, useParams } from 'react-router';
import styled from 'styled-components';
import { Guild } from '../../lib/shinpuru-ts/src';
import { GuildSelect } from '../GuildSelect';
import { Entry } from './Entry';
import { Section } from './Section';
import SPHeader from '../../assets/sp-header.svg'; // Imported as path so that it can be bundled separately because it is kinda big.
import { ReactComponent as UsersIcon } from '../../assets/users.svg';
import { useGuilds } from '../../hooks/useGuilds';
import { useStore } from '../../services/store';

interface Props {}

const StyledNav = styled.nav`
  display: flex;
  flex-direction: column;
  background-color: ${(p) => p.theme.background2};
  margin: 1em;
  padding: 1em;
  border-radius: 12px;
  width: 25vw;
  max-width: 15em;
`;

const EntryContainer = styled.div`
  margin-top: 1em;
`;

export const Navbar: React.FC<Props> = () => {
  const { t } = useTranslation('components');
  const nav = useNavigate();
  const { guildid } = useParams();
  const guilds = useGuilds();
  const [selectedGuild, setSelectedGuild] = useStore((s) => [
    s.selectedGuild,
    s.setSelectedGuild,
  ]);

  useEffect(() => {
    if (!!guilds)
      setSelectedGuild(guilds.find((g) => g.id === guildid) ?? guilds[0]);
  }, [guildid, guilds]);

  const _onGuildSelect = (g: Guild) => {
    setSelectedGuild(g);
    nav(`guilds/${g.id}/members`);
  };

  return (
    <StyledNav>
      <img src={SPHeader} width="auto" height="auto" alt="shinpuru Heading" />
      <Section title={t('navbar.section.guilds.title')}>
        <GuildSelect
          guilds={guilds ?? []}
          value={selectedGuild}
          onElementSelect={_onGuildSelect}
        />
        <EntryContainer>
          {/* <Entry path={`guilds/${selectedGuild?.id}/home`}>
            <HomeIcon />
            {t('navbar.section.guilds.home')}
          </Entry> */}
          <Entry path={`guilds/${selectedGuild?.id}/members`}>
            <UsersIcon />
            {t('navbar.section.guilds.members')}
          </Entry>
        </EntryContainer>
      </Section>
      {/* <Section title={t('navbar.section.users')}></Section> */}
    </StyledNav>
  );
};
