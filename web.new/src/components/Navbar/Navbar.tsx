import { useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate, useParams } from 'react-router';
import styled from 'styled-components';
import SPBrand from '../../assets/sp-brand.svg';
import SPIcon from '../../assets/sp-icon.png';
import { ReactComponent as UsersIcon } from '../../assets/users.svg';
import { useGuilds } from '../../hooks/useGuilds';
import { Guild } from '../../lib/shinpuru-ts/src';
import { useStore } from '../../services/store';
import { GuildSelect } from '../GuildSelect';
import { Heading } from '../Heading';
import { Entry } from './Entry';
import { Section } from './Section';

interface Props {}

const Brand = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;

  > img {
    width: 100%;
    &:first-child {
      width: 38px;
      height: 38px;
    }
  }
`;

const EntryContainer = styled.div`
  margin-top: 1em;
`;

const StyledEntry = styled(Entry)`
  height: 100%;
`;

const StyledGuildSelect = styled(GuildSelect)``;

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

    ${Brand} > img:last-child {
      display: none;
    }

    ${StyledEntry} span {
      display: none;
    }

    ${StyledGuildSelect} > div > div > span {
      display: none;
    }

    ${Heading} {
      display: none;
    }
  }
`;

export const Navbar: React.FC<Props> = () => {
  const { t } = useTranslation('components');
  const nav = useNavigate();
  const { guildid } = useParams();
  const guilds = useGuilds();
  const [selectedGuild, setSelectedGuild] = useStore((s) => [s.selectedGuild, s.setSelectedGuild]);

  useEffect(() => {
    if (!!guilds && !!guildid) setSelectedGuild(guilds.find((g) => g.id === guildid) ?? guilds[0]);
  }, [guildid, guilds]);

  const _onGuildSelect = (g: Guild) => {
    setSelectedGuild(g);
    nav(`guilds/${g.id}/members`);
  };

  return (
    <StyledNav>
      <Brand>
        <img src={SPIcon} alt="shinpuru Heading" />
        <img src={SPBrand} alt="shinpuru Heading" />
      </Brand>
      <Section title={t('navbar.section.guilds.title')}>
        <StyledGuildSelect
          guilds={guilds ?? []}
          value={selectedGuild}
          onElementSelect={_onGuildSelect}
        />
        <EntryContainer>
          <StyledEntry path={`guilds/${selectedGuild?.id}/members`}>
            <UsersIcon />
            <span>{t('navbar.section.guilds.members')}</span>
          </StyledEntry>
        </EntryContainer>
      </Section>
      {/* <Section title={t('navbar.section.users')}></Section> */}
    </StyledNav>
  );
};
