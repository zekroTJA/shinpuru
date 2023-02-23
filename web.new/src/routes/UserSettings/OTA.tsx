import React, { useEffect, useState } from 'react';
import { Trans, useTranslation } from 'react-i18next';

import { Button } from '../../components/Button';
import { Controls } from '../../components/Controls';
import { Embed } from '../../components/Embed';
import { Loader } from '../../components/Loader';
import { MaxWidthContainer } from '../../components/MaxWidthContainer';
import { Small } from '../../components/Small';
import { Switch } from '../../components/Switch';
import { useApi } from '../../hooks/useApi';
import { useNotifications } from '../../hooks/useNotifications';

type Props = {};

const OTARoute: React.FC<Props> = ({}) => {
  const { t } = useTranslation('routes.usersettings.ota');
  const { pushNotification } = useNotifications();
  const fetch = useApi();
  const [enabled, setEnabled] = useState<boolean>();

  const _saveSettings = () => {
    if (enabled === undefined) return;
    fetch((c) => c.usersettings.setOta({ enabled }))
      .then(() =>
        pushNotification({
          message: t('notifications.saved'),
          type: 'SUCCESS',
        }),
      )
      .catch();
  };

  useEffect(() => {
    fetch((c) => c.usersettings.ota())
      .then((r) => setEnabled(r.enabled))
      .catch();
  }, []);

  return (
    <MaxWidthContainer>
      <h1>{t('heading')}</h1>
      <Small>
        <Trans
          ns="routes.usersettings.ota"
          i18nKey="explanation"
          components={{
            code: <Embed />,
            '1': (
              <a
                href="https://github.com/zekroTJA/shinpuru/wiki/One-Time-Authentication-(OTA)"
                target="_blank"
                rel="noreferrer">
                link
              </a>
            ),
          }}></Trans>
      </Small>
      {(enabled !== undefined && (
        <Switch enabled={enabled} onChange={setEnabled} labelAfter={t('toggle')} />
      )) || <Loader width="15em" height="2.5em" />}
      <Controls>
        <Button variant="green" onClick={_saveSettings}>
          {t('save')}
        </Button>
      </Controls>
    </MaxWidthContainer>
  );
};

export default OTARoute;
