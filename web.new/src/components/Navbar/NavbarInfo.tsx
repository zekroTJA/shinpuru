import { NAVBAR_WIDTH, Navbar } from './Navbar';

import { ReactComponent as BackIcon } from '../../assets/back.svg';
import { ReactComponent as BookIcon } from '../../assets/book.svg';
import { Button } from '../Button';
import { ReactComponent as CommandsIcon } from '../../assets/command.svg';
import { Entry } from './Entry';
import { EntryContainer } from './EntryContainer';
import { GuildSelect } from '../GuildSelect';
import { MAX_WIDTH } from '../MaxWidthContainer';
import { ReactComponent as SystemIcon } from '../../assets/cpu.svg';
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
  @media (orientation: portrait) and (max-width: calc(${NAVBAR_WIDTH} + ${MAX_WIDTH})) {
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

export const NavbarInfo: React.FC<Props> = () => {
  const { t } = useTranslation('components', { keyPrefix: 'navbar-info' });
  const nav = useNavigate();

  return (
    <StyledNavbar>
      <EntryContainer>
        <div>
          <StyledEntry path={`/info/general`}>
            <BookIcon />
            <span>{t('section.default.general')}</span>
          </StyledEntry>
          <StyledEntry path={`/info/commands`}>
            <CommandsIcon />
            <span>{t('section.default.commands')}</span>
          </StyledEntry>
          <StyledEntry path={`/info/system`}>
            <SystemIcon />
            <span>{t('section.default.system')}</span>
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
