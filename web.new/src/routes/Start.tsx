import { Trans, useTranslation } from 'react-i18next';
import styled, { useTheme } from 'styled-components';

import { Button } from '../components/Button';
import Color from 'color';
import { LinearGradient } from '../components/styleParts';
import { ReactComponent as LoginIcon } from '../assets/login.svg';
import MockupCodeExecDark from '../assets/mockups/dark/code-execution.svg';
import MockupCodeExecLight from '../assets/mockups/light/code-execution.svg';
import { ReactComponent as SPBrand } from '../assets/sp-brand.svg';
import SPIcon from '../assets/sp-icon.png';
import Wave from 'react-wavify';
import { useNavigate } from 'react-router';

type Props = {};

const StartContainer = styled.div``;

const Header = styled.header`
  display: flex;
  flex-direction: column;
  gap: 3em;
  align-items: center;
  padding-top: 20vh;

  > span {
    font-family: 'Cantarell';
    font-size: 1.1rem;
    font-weight: lighter;
    text-align: center;
    max-width: 20em;
    opacity: 0.9;
  }
`;

const HeaderButtons = styled.div`
  display: flex;
  gap: 2em;

  ${Button} {
    transition: all 0.25s ease;
    padding: 0.8em 2em;
    box-shadow: 0 0 2em 0 ${(p) => Color(p.theme.accent).alpha(0.2).hexa()};
    &:hover {
      box-shadow: 0 0 2em 0 ${(p) => Color(p.theme.accent).alpha(0.4).hexa()};
    }
  }
`;

const GlowLink = styled.a`
  ${(p) => LinearGradient(p.theme.accent)}
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  display: inline-block;
  text-decoration: none;
  text-shadow: 0 0 0.8em ${(p) => Color(p.theme.accent).alpha(0.8).hexa()};
`;

const Brand = styled.div`
  display: flex;
  gap: 1em;
  align-items: center;

  width: 80vw;
  height: 15vw;

  max-height: 6rem;
  max-width: 30rem;

  > img {
    height: 100%;
  }

  > svg {
    width: 100%;
    height: 100%;
  }
`;

const Background = styled(Wave)`
  position: fixed;
  bottom: 0;
  width: 100%;
  height: 30em;
  z-index: -1;
`;

const Features = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2em;
  margin: 8em 4em 4em 4em;

  h1 {
    text-transform: uppercase;
    opacity: 0.8;
  }

  > div {
    display: flex;
    gap: 2em;
    width: 100%;
    max-width: 80em;
    padding: 2em;
    border-radius: 12px;
    background-color: ${(p) => p.theme.background2};

    > img {
      max-width: 20em;
      width: 40vw;
      height: auto;
      z-index: 5;
      border-radius: 8px;
      box-shadow: 0 1em 2em 0 rgba(0 0 0 / 25%);
    }

    span {
      font-size: 1.4rem;
      line-height: 1.5em;
      font-weight: lighter;
    }

    &:nth-child(2) {
      flex-direction: row-reverse;
    }
  }

  @media (max-width: 50em) {
    margin: 8em 1em 4em 1em;
  }

  @media (max-width: 40em) {
    > div {
      flex-direction: column !important;
      align-items: center;

      img {
        max-width: 100%;
        width: 100%;
        height: auto;
      }
    }
  }
`;

const LoginButton = styled(Button)`
  position: fixed;
  top: 1.5em;
  right: 1.5em;

  width: 3em;
  height: 3em;
  padding: 0 0.6em;
  display: flex;
  justify-content: flex-start;
  gap: 1em;
  overflow: hidden;
  background: ${(p) => p.theme.background3};
  opacity: 0.5;

  transition: all 0.25s ease;
  transform: none !important;

  > svg {
    min-height: 2em;
    min-width: 2em;
  }

  &:hover {
    width: 8em;
    background: ${(p) => p.theme.accent};
    opacity: 1;
  }
`;

const Footer = styled.footer``;

export const StartRoute: React.FC<Props> = () => {
  const { t } = useTranslation('routes.start');
  const nav = useNavigate();
  const theme = useTheme();

  return (
    <StartContainer>
      <LoginButton onClick={() => nav('/login')}>
        <LoginIcon />
        {t('login')}
      </LoginButton>
      <Background
        fill="url(#gradient)"
        options={{
          height: 200,
          amplitude: 200,
          speed: 0.05,
          points: 4,
        }}>
        <defs>
          <linearGradient id="gradient" gradientTransform="rotate(90)">
            <stop offset="10%" stopColor={theme.accent} />
            <stop offset="90%" stopColor={Color(theme.accent).darken(0.3).hexa()} />
          </linearGradient>
        </defs>
      </Background>
      <Header>
        <Brand>
          <img src={SPIcon} alt="shinpuru icon" />
          <SPBrand />
        </Brand>
        <span>
          <Trans
            ns="routes.start"
            i18nKey="header.under"
            components={{
              '1': (
                <GlowLink
                  href="https://github.com/zekrotja/shinpuru"
                  target="_blank"
                  rel="noreferrer">
                  _
                </GlowLink>
              ),
            }}
          />
        </span>
        <HeaderButtons>
          <a href="/invite">
            <Button>{t('header.invite')}</Button>
          </a>
          <a href="https://github.com/zekroTJA/shinpuru/wiki/Self-Hosting">
            <Button>{t('header.selfhost')}</Button>
          </a>
        </HeaderButtons>
      </Header>
      <main>
        <Features>
          <div>
            <img src={theme._isDark ? MockupCodeExecDark : MockupCodeExecLight} alt="" />
            <div>
              <h1>{t('features.codeexecution.heading')}</h1>
              <span>{t('features.codeexecution.description')}</span>
            </div>
          </div>
          <div>
            <img src={theme._isDark ? MockupCodeExecDark : MockupCodeExecLight} alt="" />
            <div>
              <h1>{t('features.codeexecution.heading')}</h1>
              <span>{t('features.codeexecution.description')}</span>
            </div>
          </div>
          <div>
            <img src={theme._isDark ? MockupCodeExecDark : MockupCodeExecLight} alt="" />
            <div>
              <h1>{t('features.codeexecution.heading')}</h1>
              <span>{t('features.codeexecution.description')}</span>
            </div>
          </div>
        </Features>
      </main>
      <Footer></Footer>
    </StartContainer>
  );
};
