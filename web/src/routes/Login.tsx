import styled from 'styled-components';
import { Button } from '../components/Button';
import { Container } from '../components/Container';
import { ReactComponent as DiscordIcon } from '../assets/dc-logo-blurple.svg';
import { LOGIN_ROUTE } from '../services/api';
import { Hider } from '../components/Hider';
import { useApi } from '../hooks/useApi';
import { useEffect, useState } from 'react';
import { getCryptoRandomString } from '../util/crypto';
import { useNavigate } from 'react-router';
import { Embed } from '../components/Embed';
import { Trans, useTranslation } from 'react-i18next';

interface Props {}

const OuterContainer = styled.div`
  display: flex;
  height: 100%;
  padding: 2rem;
  font-size: 1.3rem;

  > div:first-child {
    margin: 0 1rem 0 0;
  }

  > div:last-child {
    margin: 0 0 0 1rem;
  }

  @media (orientation: portrait) {
    flex-direction: column;

    > div:first-child {
      margin: 0 0 1rem 0;
    }

    > div:last-child {
      margin: 1rem 0 0 0;
    }
  }
`;

const Tile = styled(Container)`
  width: 100%;
  padding: 1em;
  font-family: 'Cantarell', sans-serif;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  height: 100%;

  > a {
    text-decoration: none;
  }
`;

const TileDiscord = styled(Tile)`
  background: linear-gradient(
    140deg,
    ${(p) => p.theme.blurple} 0%,
    ${(p) => p.theme.blurpleDarker} 100%
  );

  ${Button} {
    background-color: ${(p) => p.theme.white};
    color: ${(p) => p.theme.blurple};
  }
`;

const TileAlt = styled(Tile)`
  background: linear-gradient(
    140deg,
    ${(p) => p.theme.white} 0%,
    ${(p) => p.theme.whiteDarker} 100%
  );

  color: ${(p) => p.theme.darkGray};
`;

const StyledHider = styled(Hider)`
  color: ${(p) => p.theme.white};

  > input {
    background-color: ${(p) => p.theme.darkGray};
  }
`;

const StyledSmall = styled.small`
  margin-top: 1.4em;
  font-size: 0.7em;
`;

export const LoginRoute: React.FC<Props> = ({}) => {
  const { t } = useTranslation('routes.login');
  const [pushCode, setPushCode] = useState('');
  const fetch = useApi();
  const nav = useNavigate();

  useEffect(() => {
    _pushCodeLoop();
  }, []);

  const _generatePushCode = async () => {
    const code = getCryptoRandomString(16);
    setPushCode(code);
    try {
      await fetch((c) => c.auth.pushCode(code));
      nav('/db');
      return false;
    } catch {}
    return true;
  };

  const _pushCodeLoop = async () => {
    while (await _generatePushCode());
  };

  return (
    <OuterContainer>
      <TileDiscord>
        <p>{t('discord.title')}</p>
        <a href={LOGIN_ROUTE}>
          <Button>
            <DiscordIcon />
            {t('discord.action')}
          </Button>
        </a>
        <StyledSmall>{t('discord.subline')}</StyledSmall>
      </TileDiscord>
      <TileAlt>
        <p>{t('alternative.title')}</p>
        <StyledHider content={pushCode} />
        <StyledSmall>
          <Trans i18nKey="alternative.subline" ns="routes.login">
            Alternatively, you can also use the <Embed>/login</Embed> command.
          </Trans>
        </StyledSmall>
      </TileAlt>
    </OuterContainer>
  );
};