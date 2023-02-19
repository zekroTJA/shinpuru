import { ReactComponent as APITokenIcon } from '../../assets/key.svg';
import { ReactComponent as BackIcon } from '../../assets/back.svg';
import { Button } from '../Button';
import { Entry } from './Entry';
import { EntryContainer } from './EntryContainer';
import { GuildSelect } from '../GuildSelect';
import { Navbar } from './Navbar';
import { ReactComponent as ShieldIcon } from '../../assets/shield.svg';
import { ReactComponent as TicketIcon } from '../../assets/ticket.svg';
import styled from 'styled-components';
import { useEffect } from 'react';
import { useNavigate } from 'react-router';
import { useTranslation } from 'react-i18next';

type Props = {};

const StyledEntry = styled(Entry)``;

const StyledGuildSelect = styled(GuildSelect)`
  margin-top: 1em;
`;

const BackButton = styled(Button)`
  margin-top: auto;
  border-radius: 8px;

  svg {
    width: 1em;
    height: 1em;
  }
`;

const StyledNavbar = styled(Navbar)`
  @media (orientation: portrait) {
    ${StyledEntry} {
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

    ${BackButton} > span {
      display: none;
    }
  }
`;

export const NavbarUserSettings: React.FC<Props> = ({}) => {
  const { t } = useTranslation('components', { keyPrefix: 'navbar-usersettings' });
  const nav = useNavigate();

  useEffect(() => {}, []);

  return (
    <StyledNavbar>
      <EntryContainer>
        <div>
          <StyledEntry path={`/usersettings/apitoken`}>
            <APITokenIcon />
            <span>{t('section.default.apitoken')}</span>
          </StyledEntry>
          <StyledEntry path={`/usersettings/ota`}>
            <TicketIcon />
            <span>{t('section.default.ota')}</span>
          </StyledEntry>
          <StyledEntry path={`/usersettings/privacy`}>
            <ShieldIcon />
            <span>{t('section.default.privacy')}</span>
          </StyledEntry>
        </div>
      </EntryContainer>
      <BackButton onClick={() => nav('/db')}>
        <BackIcon />
        <span>{t('back')}</span>
      </BackButton>
    </StyledNavbar>
  );
};
