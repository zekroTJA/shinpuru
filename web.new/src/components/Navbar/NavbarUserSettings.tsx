import { ReactComponent as APITokenIcon } from '../../assets/key.svg';
import { ReactComponent as BackIcon } from '../../assets/back.svg';
import { Button } from '../Button';
import { Entry } from './Entry';
import { EntryContainer } from './EntryContainer';
import { GuildSelect } from '../GuildSelect';
import { Navbar } from './Navbar';
import styled from 'styled-components';
import { useApi } from '../../hooks/useApi';
import { useEffect } from 'react';
import { useNavigate } from 'react-router';
import { useTranslation } from 'react-i18next';

type Props = {};

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
`;

const BackButton = styled(Button)`
  margin-top: auto;
`;

export const NavbarUserSettings: React.FC<Props> = ({}) => {
  const { t } = useTranslation('components', { keyPrefix: 'navbar-usersettings' });
  const nav = useNavigate();
  const fetch = useApi();

  useEffect(() => {}, []);

  return (
    <StyledNavbar>
      <EntryContainer>
        <StyledEntry path={`/usersettings/apitoken`}>
          <APITokenIcon />
          <span>{t('section.default.apitoken')}</span>
        </StyledEntry>
      </EntryContainer>
      <BackButton onClick={() => nav('/db')}>
        <BackIcon />
        <span>{t('back')}</span>
      </BackButton>
    </StyledNavbar>
  );
};
