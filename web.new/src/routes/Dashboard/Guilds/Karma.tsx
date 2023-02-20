import React, { useEffect } from 'react';

import { useApi } from '../../../hooks/useApi';
import { useNotifications } from '../../../hooks/useNotifications';
import { useParams } from 'react-router';
import { useTranslation } from 'react-i18next';

type Props = {};

const KarmaRoute: React.FC<Props> = ({}) => {
  const { t } = useTranslation('routes.guildstarboard');
  const { pushNotification } = useNotifications();
  const { guildid } = useParams();
  const fetch = useApi();

  useEffect(() => {
    if (!guildid) return;
  }, [guildid]);

  return <>karma</>;
};

export default KarmaRoute;
