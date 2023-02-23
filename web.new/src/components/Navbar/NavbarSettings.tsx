import { ReactComponent as BackIcon } from '../../assets/back.svg';
import { Button } from '../Button';
import { Entry } from './Entry';
import { EntryContainer } from './EntryContainer';
import { GuildSelect } from '../GuildSelect';
import { Navbar } from './Navbar';
import { ReactComponent as PresenceIcon } from '../../assets/halfmoon.svg';
import styled from 'styled-components';
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

export const NavbarSettings: React.FC<Props> = ({}) => {
  const { t } = useTranslation('components', { keyPrefix: 'navbar-settings' });
  const nav = useNavigate();

  return (
    <StyledNavbar>
      <EntryContainer>
        <div>
          <StyledEntry path={`/settings/presence`}>
            <PresenceIcon />
            <span>{t('section.default.presence')}</span>
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
