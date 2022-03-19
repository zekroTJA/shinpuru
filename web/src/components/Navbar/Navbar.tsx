import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate, useParams } from 'react-router';
import styled from 'styled-components';
import { useApi } from '../../hooks/useApi';
import { Guild } from '../../lib/shinpuru-ts/src';
import { GuildSelect } from '../GuildSelect';
import { Entry } from './Entry';
import { Section } from './Section';
import { ReactComponent as SPHeader } from '../../assets/sp-header.svg';
import { ReactComponent as HomeIcon } from '../../assets/home.svg';
import { ReactComponent as UsersIcon } from '../../assets/users.svg';

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

export const Navbar: React.FC<Props> = ({}) => {
  const { t } = useTranslation('components');
  const fetch = useApi();
  const nav = useNavigate();
  const { guildid } = useParams();

  const [guilds, setGuilds] = useState<Guild[]>([]);
  const [selectedGuild, setSelectedGuild] = useState<Guild>();

  useEffect(() => {
    fetch((c) => c.guilds.list())
      .then((r) => {
        setGuilds(r.data);
        const guild = r.data.find((g) => g.id === guildid) ?? r.data[0];
        setSelectedGuild(guild);
      })
      .catch();
  }, []);

  const _onGuildSelect = (g: Guild) => {
    setSelectedGuild(g);
    nav(`guilds/${g.id}/members`);
  };

  return (
    <StyledNav>
      <SPHeader width="auto" height="auto" />
      <Section title={t('navbar.section.guilds.title')}>
        <GuildSelect
          guilds={guilds}
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
