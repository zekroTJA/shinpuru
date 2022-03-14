import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import styled from 'styled-components';
import { ReactComponent as SPHeader } from '../../assets/sp-header.svg';
import { useApi } from '../../hooks/useApi';
import { Guild } from '../../lib/shinpuru-ts/src';
import { GuildSelect } from '../GuildSelect';
import { Section } from './Section';

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

export const Navbar: React.FC<Props> = ({}) => {
  const { t } = useTranslation('components');
  const fetch = useApi();
  const [guilds, setGuilds] = useState<Guild[]>([]);
  const [selectedGuild, setSelectedGuild] = useState<Guild>();

  useEffect(() => {
    fetch((c) => c.guilds.list())
      .then((r) => {
        setGuilds(r.data);
        setSelectedGuild(r.data[0]);
      })
      .catch();
  }, []);

  return (
    <StyledNav>
      <SPHeader width="auto" height="auto" />
      <Section title={t('navbar.section.guilds')}>
        <GuildSelect
          guilds={guilds}
          value={selectedGuild}
          onElementSelect={(g) => setSelectedGuild(g)}
        />
      </Section>
      <Section title={t('navbar.section.users')}></Section>
    </StyledNav>
  );
};
