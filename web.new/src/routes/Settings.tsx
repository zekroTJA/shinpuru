import { NavbarSettings } from '../components/Navbar';
import { Outlet } from 'react-router';
import React from 'react';
import styled from 'styled-components';

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

const SettingsRoute: React.FC<Props> = () => {
  return (
    <RouteContainer>
      <NavbarSettings />
      <RouterOutlet>
        <Outlet />
      </RouterOutlet>
    </RouteContainer>
  );
};

export default SettingsRoute;
