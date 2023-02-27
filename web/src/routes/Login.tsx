import { Trans, useTranslation } from 'react-i18next';
import { useEffect, useState } from 'react';

import { Button } from '../components/Button';
import { Container } from '../components/Container';
import { ReactComponent as DiscordIcon } from '../assets/dc-logo-blurple.svg';
import { Embed } from '../components/Embed';
import { Hider } from '../components/Hider';
import { LinearGradient } from '../components/styleParts';
import { getCryptoRandomString } from '../util/crypto';
import { loginRoute } from '../services/api';
import styled from 'styled-components';
import { useApi } from '../hooks/useApi';
import { useNavigate } from 'react-router';
import { useSearchParams } from 'react-router-dom';

type Props = {};

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
  color: ${(p) => p.theme.textAlt};
  ${(p) => LinearGradient(p.theme.blurple)};

  ${Button} {
    background-color: ${(p) => p.theme.white};
    color: ${(p) => p.theme.blurple};
  }
`;

const TileAlt = styled(Tile)`
  ${(p) => LinearGradient(p.theme.white)};

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

const DiscordButton = styled(Button)`
  background: ${(p) => p.theme.white};
`;

const LoginRoute: React.FC<Props> = () => {
  const { t } = useTranslation('routes.login');
  const [pushCode, setPushCode] = useState('');
  const fetch = useApi();
  const nav = useNavigate();
  const [params] = useSearchParams();

  useEffect(() => {
    _pushCodeLoop();
  }, []);

  const redirect = params.get('redirect');

  const _generatePushCode = async () => {
    const code = getCryptoRandomString(16);
    setPushCode(code);
    try {
      await fetch((c) => c.auth.pushCode(code), true);
      nav('/' + redirect ?? 'db');
      return false;
    } catch {}
    return true;
  };

  const _pushCodeLoop = async () => {
    while (await _generatePushCode());
  };

  // TODO: remove beta redirect when going live
  const _loginRoute = loginRoute(!!redirect ? `/${redirect}/` : '/');

  return (
    <OuterContainer>
      <TileDiscord>
        <p>{t('discord.title')}</p>
        <a href={_loginRoute}>
          <DiscordButton>
            <DiscordIcon />
            {t('discord.action')}
          </DiscordButton>
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

export default LoginRoute;
