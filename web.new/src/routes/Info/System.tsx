import React, { useEffect } from 'react';

import styled from 'styled-components';
import { useApi } from '../../hooks/useApi';
import { useNotifications } from '../../hooks/useNotifications';
import { useParams } from 'react-router';
import { useTranslation } from 'react-i18next';

type Props = {};

const SystemRoute: React.FC<Props> = ({}) => {
  const { t } = useTranslation('routes.info.system.json');
  const { pushNotification } = useNotifications();
  const { guildid } = useParams();
  const fetch = useApi();

  useEffect(() => {}, []);

  return <>system</>;
};

export default SystemRoute;
