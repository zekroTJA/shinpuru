import { Outlet, useParams } from 'react-router';
import React, { useEffect } from 'react';

import { NavbarUserSettings } from '../components/Navbar';
import styled from 'styled-components';
import { useApi } from '../hooks/useApi';
import { useNotifications } from '../hooks/useNotifications';
import { useTranslation } from 'react-i18next';

type Props = {};

const RouteContainer = styled.div`
  display: flex;
  height: 100%;
`;

const RouterOutlet = styled.main`
  padding: 1em;
  width: 100%;
  height: 100%;
  overflow-y: auto;
`;

const UserSettingsRoute: React.FC<Props> = ({}) => {
  const { t } = useTranslation('routes.');
  const { pushNotification } = useNotifications();
  const fetch = useApi();

  useEffect(() => {}, []);

  return (
    <RouteContainer>
      <NavbarUserSettings />
      <RouterOutlet>
        <Outlet />
      </RouterOutlet>
    </RouteContainer>
  );
};

export default UserSettingsRoute;
