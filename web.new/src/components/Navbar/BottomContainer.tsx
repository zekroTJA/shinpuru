import { Button } from '../Button';
import { DiscordImage } from '../DiscordImage';
import { ReactComponent as GlobalSettingsIcon } from '../../assets/settings.svg';
import { Hoverplate } from '../Hoverplate';
import { ReactComponent as InfoIcon } from '../../assets/info.svg';
import { Loader } from '../Loader';
import { ReactComponent as LogoutIcon } from '../../assets/logout.svg';
import { Styled } from '../props';
import { ReactComponent as TriangleIcon } from '../../assets/triangle.svg';
import { ReactComponent as UserSettingsIcon } from '../../assets/tool.svg';
import styled from 'styled-components';
import { useApi } from '../../hooks/useApi';
import { useNavigate } from 'react-router';
import { useSelfUser } from '../../hooks/useSelfUser';
import { useTranslation } from 'react-i18next';

type Props = Styled & {};

const HoverplateButton = styled(Button)`
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 0.5em;
`;

const StyledHoverplate = styled(Hoverplate)`
  width: 100%;
`;

const HoverplateContent = styled.div`
  display: flex;
  flex-direction: column;
  gap: 1em;
  width: 30vw;
  max-width: 13rem;
`;

const Container = styled.div`
  margin-top: auto;
  gap: 0.5em;
  display: flex;
  padding-top: 1em;

  > ${Button} {
    padding: 0.7em;
    border-radius: 8px;

    > svg {
      width: 1.5em;
      height: 1.5em;
    }
  }

  @media (orientation: portrait) {
    flex-direction: column;
  }
`;

export const SelfContainer = styled.div`
  display: flex;
  align-items: center;
  width: 100%;
  gap: 1em;
  background-color: ${(p) => p.theme.background3};
  border-radius: 8px;
  padding: 0.5em;

  > img {
    width: 2em;
    height: 2em;
  }

  > span {
    align-items: center;
    width: 100%;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  > svg {
    height: 0.5em;
    margin: 0 0.5em 0 auto;
  }
`;

const SettingsButton = styled(Button)`
  color: ${(p) => p.theme.text};
  z-index: 1;
`;

export const BottomContainer: React.FC<Props> = ({ ...props }) => {
  const { t } = useTranslation('components', { keyPrefix: 'navbar' });
  const nav = useNavigate();
  const self = useSelfUser();
  const fetch = useApi();

  const _logout = () => {
    fetch((c) => c.auth.logout())
      .then(() => nav('/start'))
      .catch();
  };

  return (
    <Container {...props}>
      <StyledHoverplate
        hoverContent={
          <HoverplateContent>
            <HoverplateButton onClick={() => nav('/usersettings')}>
              <UserSettingsIcon />
              <span>{t('user-settings')}</span>
            </HoverplateButton>
            <HoverplateButton variant="blue" onClick={() => nav('/info')}>
              <InfoIcon />
              <span>{t('info')}</span>
            </HoverplateButton>
            <HoverplateButton variant="orange" onClick={_logout}>
              <LogoutIcon />
              <span>{t('logout')}</span>
            </HoverplateButton>
          </HoverplateContent>
        }>
        {(self && (
          <SelfContainer>
            <DiscordImage src={self?.avatar_url} round />
            <span>{self?.username}</span>
            <TriangleIcon />
          </SelfContainer>
        )) || <Loader width="100%" height="2em" borderRadius="8px" />}
      </StyledHoverplate>
      {self?.bot_owner && (
        <SettingsButton variant="gray" onClick={() => nav('/settings')}>
          <GlobalSettingsIcon />
        </SettingsButton>
      )}
    </Container>
  );
};
