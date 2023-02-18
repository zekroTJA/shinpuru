import React, { useEffect, useState } from 'react';
import { Trans, useTranslation } from 'react-i18next';
import styled, { useTheme } from 'styled-components';

import { APIError } from '../../lib/shinpuru-ts/src/errors';
import { APIToken } from '../../lib/shinpuru-ts/src';
import { Button } from '../../components/Button';
import { Card } from '../../components/Card';
import { Controls } from '../../components/Controls';
import { Flex } from '../../components/Flex';
import { Hider } from '../../components/Hider';
import { ReactComponent as InfoIcon } from '../../assets/info.svg';
import { Loader } from '../../components/Loader';
import { MaxWidthContainer } from '../../components/MaxWidthContainer';
import { Small } from '../../components/Small';
import { ReactComponent as WarnIcon } from '../../assets/warning.svg';
import { formatDate } from '../../util/date';
import { useApi } from '../../hooks/useApi';
import { useNotifications } from '../../hooks/useNotifications';

type Props = {};

type TokenState = Partial<APIToken> & {
  hasToken: boolean;
};

const TokenHider = styled(Hider)`
  > input {
    width: 100%;
    max-width: 100%;
    background-color: ${(p) => p.theme.background3};
    color: ${(p) => p.theme.text};
  }
`;

const TokenControls = styled(Controls)`
  margin: 2em 0 1em 0;
`;

const StyledCard = styled(Card)`
  width: 100%;
  gap: 1em;
  display: flex;
  justify-content: flex-start;
  align-items: center;

  svg {
    width: 2em;
    height: 2em;
  }
`;

const APITokenRoute: React.FC<Props> = ({}) => {
  const { t } = useTranslation('routes.usersettings.apitoken');
  const { pushNotification } = useNotifications();
  const theme = useTheme();
  const fetch = useApi();
  const [token, setToken] = useState<TokenState>();

  useEffect(() => {
    fetch((c) => c.tokens.info(), 404)
      .then((r) => setToken({ ...r, hasToken: true }))
      .catch((e) => {
        if (e instanceof APIError && e.code === 404) setToken({ hasToken: false });
      });
  }, []);

  const _generateToken = () => {
    fetch((c) => c.tokens.generate())
      .then((r) => {
        setToken({ ...r, hasToken: true });
        pushNotification({
          message: t('notifications.generated'),
          type: 'SUCCESS',
        });
      })
      .catch();
  };

  const _resetToken = () => {
    fetch((c) => c.tokens.delete())
      .then((r) => {
        setToken({ hasToken: false });
        pushNotification({
          message: t('notifications.reset'),
          type: 'SUCCESS',
        });
      })
      .catch();
  };

  return (
    <MaxWidthContainer>
      <h1>{t('heading')}</h1>
      <Small>
        <Trans
          ns="routes.usersettings.apitoken"
          i18nKey="explanation"
          components={{
            '1': (
              <a
                href="https://github.com/zekroTJA/shinpuru/wiki/REST-API-Docs"
                target="_blank"
                rel="noreferrer">
                link
              </a>
            ),
          }}></Trans>
      </Small>
      <TokenControls>
        {(token && (
          <>
            <Button onClick={_generateToken}>
              {token.hasToken ? t('regenerate') : t('generate')}
            </Button>
            <Button onClick={_resetToken} disabled={!token.hasToken} variant="orange">
              {t('reset')}
            </Button>
          </>
        )) || <Loader height="3em" />}
      </TokenControls>
      {token?.token && (
        <Flex gap="1em" direction="column">
          <StyledCard color={theme.orange}>
            <div>
              <WarnIcon />
            </div>
            <div>{t('tokenwarning')}</div>
          </StyledCard>
          <TokenHider content={token.token} />
        </Flex>
      )}
      {token?.hasToken && !token.token && (
        <StyledCard color={theme.blurple}>
          <div>
            <InfoIcon />
          </div>
          <div>
            <Trans
              ns="routes.usersettings.apitoken"
              i18nKey="tokeninfo"
              shouldUnescape={true}
              values={{
                created: formatDate(token.created),
                expires: formatDate(token.expires),
                hits: token.hits,
              }}
              components={{ b: <b />, br: <br /> }}
            />
          </div>
        </StyledCard>
      )}
    </MaxWidthContainer>
  );
};

export default APITokenRoute;
