import { useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate, useParams } from 'react-router';
import styled from 'styled-components';
import { ReactComponent as SPBrand } from '../../assets/sp-brand.svg';
import SPIcon from '../../assets/sp-icon.png';
import { ReactComponent as TriangleIcon } from '../../assets/triangle.svg';
import { ReactComponent as UsersIcon } from '../../assets/users.svg';
import { useApi } from '../../hooks/useApi';
import { useGuilds } from '../../hooks/useGuilds';
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
      span {
        display: none;
      }
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
  const { t } = useTranslation('components', { keyPrefix: 'navbar' });
  const nav = useNavigate();
  const fetch = useApi();
  const { guildid } = useParams();
  const guilds = useGuilds();
  const self = useSelfUser();
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
        </Section>
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
