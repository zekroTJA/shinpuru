import React, { useEffect, useState } from 'react';
import { Trans, useTranslation } from 'react-i18next';

import { APIError } from '../../lib/shinpuru-ts/src/errors';
import { APIToken } from '../../lib/shinpuru-ts/src';
import { MaxWidthContainer } from '../../components/MaxWidthContainer';
import { Small } from '../../components/Small';
import styled from 'styled-components';
import { useApi } from '../../hooks/useApi';
import { useNotifications } from '../../hooks/useNotifications';

type Props = {};

const APITokenRoute: React.FC<Props> = ({}) => {
  const { t } = useTranslation('routes.usersettings.apitoken');
  const { pushNotification } = useNotifications();
  const fetch = useApi();
  const [token, setToken] = useState<APIToken>();

  useEffect(() => {
    fetch((c) => c.tokens.info(), 404)
      .then((r) => setToken(r))
      .catch();
  }, []);

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
    </MaxWidthContainer>
  );
};

export default APITokenRoute;
