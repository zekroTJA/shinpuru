import { Trans, useTranslation } from 'react-i18next';
import styled, { css, useTheme } from 'styled-components';

import { Button } from '../components/Button';
import Color from 'color';
import { LinearGradient } from '../components/styleParts';
import { ReactComponent as LoginIcon } from '../assets/login.svg';
import Marquee from 'react-fast-marquee';
import MockupCodeExecDark from '../assets/mockups/dark/code-execution.svg';
import MockupCodeExecLight from '../assets/mockups/light/code-execution.svg';
import MockupKarmaDark from '../assets/mockups/dark/karma.svg';
import MockupKarmaLight from '../assets/mockups/light/karma.svg';
import MockupReportsDark from '../assets/mockups/dark/reports.svg';
import MockupReportsLight from '../assets/mockups/light/reports.svg';
import MockupRoleselectDark from '../assets/mockups/dark/roleselect.svg';
import MockupRoleselectLight from '../assets/mockups/light/roleselect.svg';
import MockupStarboardDark from '../assets/mockups/dark/starboard.svg';
import MockupStarboardLight from '../assets/mockups/light/starboard.svg';
import MockupVotesDark from '../assets/mockups/dark/votes.svg';
import MockupVotesLight from '../assets/mockups/light/votes.svg';
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

const Card = styled.div`
  display: flex;
  gap: 2em;
  width: 100%;
  max-width: 80em;
  padding: 2em;
  border-radius: 12px;
  background-color: ${(p) => Color(p.theme.background2).alpha(0.8).hexa()};
  backdrop-filter: blur(5em);
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

  > ${Card} {
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

    &:nth-child(odd) {
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

const DiscoverMore = styled.div`
  margin: 0em 4em 4em 4em;
  display: flex;
  justify-content: center;

  @media (max-width: 50em) {
    margin: 0em 1em 4em 1em;
  }

  > ${Card} {
    > h1 {
      margin: 0;
    }

    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;

    > div {
      width: 100%;
      background-clip: text;
      display: inline-block;
      text-decoration: none;
    }

    span {
      text-transform: uppercase;
      font-size: 2rem;
      font-weight: lighter;
      white-space: nowrap;
      margin-right: 1em;

      @keyframes scroll {
        from {
          transform: translateX(-151ch);
        }
        to {
          transform: translateX(0ch);
        }
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
  color: ${(p) => p.theme.text};

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
    color: ${(p) => p.theme.textAlt};
  }
`;

const Footer = styled.footer`
  display: flex;
  gap: 5em;
  padding: 2em;
  justify-content: center;
  color: ${(p) => p.theme.text};
  background-color: ${(p) => Color(p.theme.background2).alpha(0.5).hexa()};
  backdrop-filter: blur(5em);

  a {
    color: inherit;
    text-decoration: underline;
  }

  > div {
    > span,
    a {
      display: block;
      line-height: 1.8rem;
    }
  }
`;

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
          <Card>
            <img src={theme._isDark ? MockupCodeExecDark : MockupCodeExecLight} alt="" />
            <div>
              <h1>{t('features.codeexecution.heading')}</h1>
              <span>{t('features.codeexecution.description')}</span>
            </div>
          </Card>
          <Card>
            <img src={theme._isDark ? MockupRoleselectDark : MockupRoleselectLight} alt="" />
            <div>
              <h1>{t('features.roleselect.heading')}</h1>
              <span>{t('features.roleselect.description')}</span>
            </div>
          </Card>
          <Card>
            <img src={theme._isDark ? MockupKarmaDark : MockupKarmaLight} alt="" />
            <div>
              <h1>{t('features.karma.heading')}</h1>
              <span>{t('features.karma.description')}</span>
            </div>
          </Card>
          <Card>
            <img src={theme._isDark ? MockupReportsDark : MockupReportsLight} alt="" />
            <div>
              <h1>{t('features.reports.heading')}</h1>
              <span>{t('features.reports.description')}</span>
            </div>
          </Card>
          <Card>
            <img src={theme._isDark ? MockupVotesDark : MockupVotesLight} alt="" />
            <div>
              <h1>{t('features.votes.heading')}</h1>
              <span>{t('features.votes.description')}</span>
            </div>
          </Card>
          <Card>
            <img src={theme._isDark ? MockupStarboardDark : MockupStarboardLight} alt="" />
            <div>
              <h1>{t('features.starboard.heading')}</h1>
              <span>{t('features.starboard.description')}</span>
            </div>
          </Card>
        </Features>
        <DiscoverMore>
          <Card>
            <h1>{t('discover.header')}</h1>
            <div>
              <Marquee gradient={false} speed={150}>
                {(t('discover.features', { returnObjects: true }) as string[]).map((v) => (
                  <span>{v}</span>
                ))}
              </Marquee>
            </div>
          </Card>
        </DiscoverMore>
      </main>
      <Footer>
        <div>
          <span>shinpuru</span>
          <span>© {new Date().getFullYear()} Ringo Hoffmann</span>
          <a
            href="https://github.com/zekroTJA/shinpuru/blob/master/LICENCE"
            target="_blank"
            rel="noreferrer">
            Covered by the MIT Licence.
          </a>
          <a href="https://github.com/zekroTJA/shinpuru" target="_blank" rel="noreferrer">
            GitHub Repository
          </a>
        </div>
        <div>
          <a href="https://shnp.de/invite" target="_blank" rel="noreferrer">
            Invite Stable
          </a>
          <a href="https://c.shnp.de/invite" target="_blank" rel="noreferrer">
            Invite Canary
          </a>
        </div>
        <div>
          <a href="https://github.com/zekroTJA/shinpuru/wiki" target="_blank" rel="noreferrer">
            Wiki
          </a>
          <a
            href="https://github.com/zekroTJA/shinpuru/wiki/Self-Hosting"
            target="_blank"
            rel="noreferrer">
            Self Host
          </a>
          <a
            href="https://github.com/zekroTJA/shinpuru/wiki/Commands"
            target="_blank"
            rel="noreferrer">
            Commands
          </a>
          <a
            href="https://github.com/zekroTJA/shinpuru/wiki/Permissions-Guide"
            target="_blank"
            rel="noreferrer">
            Permissions Guide
          </a>
          <a
            href="https://github.com/zekroTJA/shinpuru/wiki/REST-API-Docs"
            target="_blank"
            rel="noreferrer">
            REST API
          </a>
        </div>
      </Footer>
    </StartContainer>
  );
};
