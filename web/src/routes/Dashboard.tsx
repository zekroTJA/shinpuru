import { Outlet, Route, Router, Routes } from 'react-router';
import styled from 'styled-components';
import { Navbar } from '../components/Navbar';

interface Props {}

const RouteContainer = styled.div`
  display: flex;
  height: 100%;
`;

const RouterOutlet = styled.main`
  padding: 1em 1em 0 1em;
  width: 100%;
`;

export const DashboardRoute: React.FC<Props> = ({}) => {
  return (
    <RouteContainer>
      <Navbar />
      <RouterOutlet>
        <Outlet />
      </RouterOutlet>
    </RouteContainer>
  );
};
