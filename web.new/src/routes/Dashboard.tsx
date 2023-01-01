import { useEffect } from 'react';
import { Outlet, useNavigate, useParams, useLocation } from 'react-router';
import styled from 'styled-components';
import { Navbar } from '../components/Navbar';
import { useGuilds } from '../hooks/useGuilds';
import LocalStorageUtil from '../util/localstorage';

type Props = {};

const RouteContainer = styled.div`
  display: flex;
  height: 100%;
`;

const RouterOutlet = styled.main`
  padding: 1em 1em 1em 0em;
  width: 100%;
  height: 100%;
  overflow-y: auto;
`;

export const DashboardRoute: React.FC<Props> = () => {
  const guilds = useGuilds();
  const { guildid } = useParams();
  const loc = useLocation();
  const nav = useNavigate();

  useEffect(() => {
    console.log(loc.pathname.replaceAll('/', ''), guilds, guilds?.length, guildid);
    if (loc.pathname.replaceAll('/', '') === 'db' && !!guilds && guilds.length !== 0 && !guildid) {
      const guild =
        guilds.find((g) => g.id === LocalStorageUtil.get<string>('shnp.selectedguild')) ??
        guilds[0];
      nav(`guilds/${guild.id}/members`);
    }
  }, [guilds, guildid]);

  return (
    <RouteContainer>
      <Navbar />
      <RouterOutlet>
        <Outlet />
      </RouterOutlet>
    </RouteContainer>
  );
};
