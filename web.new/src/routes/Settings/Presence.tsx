import React, { useEffect } from 'react';

import styled from 'styled-components';
import { useApi } from '../../hooks/useApi';
import { useNotifications } from '../../hooks/useNotifications';
import { useParams } from 'react-router';
import { useTranslation } from 'react-i18next';

type Props = {};

const PresenceRoute: React.FC<Props> = ({}) => {
  const { t } = useTranslation('routes.');
  const { pushNotification } = useNotifications();
  const fetch = useApi();

  useEffect(() => {}, []);

  return <>presence</>;
};

export default PresenceRoute;
