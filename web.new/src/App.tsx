import { Navigate, Route, BrowserRouter as Router, Routes } from 'react-router-dom';
import React, { useEffect } from 'react';
import styled, { ThemeProvider, createGlobalStyle } from 'styled-components';

import { DashboardRoute } from './routes/Dashboard';
import { DebugRoute } from './routes/Debug';
import { HookedModal } from './components/Modal';
import { ModalBetaGreeter } from './components/Modals/ModalBetaGreeter';
import { Notifications } from './components/Notifications';
import { RouteSuspense } from './components/RouteSuspense';
import { StartRoute } from './routes/Start';
import { stripSuffix } from './util/utils';
import { useStoredTheme } from './hooks/useStoredTheme';

const LoginRoute = React.lazy(() => import('./routes/Login'));
const GuildMembersRoute = React.lazy(() => import('./routes/Dashboard/Guilds/GuildMembers'));
const MemberRoute = React.lazy(() => import('./routes/Dashboard/Guilds/Member'));
const GuildModlogRoute = React.lazy(() => import('./routes/Dashboard/Guilds/GuildModlog'));
const UnbanmeRoute = React.lazy(() => import('./routes/Unbanme'));
const VerifyRoute = React.lazy(() => import('./routes/Verify'));
const GuildGeneralRoute = React.lazy(() => import('./routes/Dashboard/GuildSettings/General'));
const GuildBackupsRoute = React.lazy(() => import('./routes/Dashboard/GuildSettings/Backup'));
const GuildAntiraidRoute = React.lazy(() => import('./routes/Dashboard/GuildSettings/Antiraid'));
const GuildCodeexecRoute = React.lazy(() => import('./routes/Dashboard/GuildSettings/Codeexec'));
const GuildVerificationRoute = React.lazy(
  () => import('./routes/Dashboard/GuildSettings/Verification'),
);
const GuildKarmaRoute = React.lazy(() => import('./routes/Dashboard/GuildSettings/Karma'));
const GuildLogsRoute = React.lazy(() => import('./routes/Dashboard/GuildSettings/Logs'));
const GuildDataRoute = React.lazy(() => import('./routes/Dashboard/GuildSettings/Data'));
const GuildAPIRoute = React.lazy(() => import('./routes/Dashboard/GuildSettings/API'));
const GuildPermissionsRoute = React.lazy(
  () => import('./routes/Dashboard/GuildSettings/Permissions'),
);
const GuildLinkBlockingRoute = React.lazy(
  () => import('./routes/Dashboard/GuildSettings/LinkBlocking'),
);

const UserSettingsRoute = React.lazy(() => import('./routes/UserSettings'));
const APITokenRoute = React.lazy(() => import('./routes/UserSettings/APIToken'));
const OTARoute = React.lazy(() => import('./routes/UserSettings/OTA'));
const PrivacyRoute = React.lazy(() => import('./routes/UserSettings/Privacy'));

const GlobalStyle = createGlobalStyle`
  body {
    background-color: ${(p) => p.theme.background};
    color: ${(p) => p.theme.text};
  }

  * {
    box-sizing: border-box;
  }
`;

const AppContainer = styled.div`
  width: 100vw;
  height: 100vh;
`;

export const App: React.FC = () => {
  const { theme } = useStoredTheme();

  useEffect(() => {
    if (
      import.meta.env.BASE_URL.length > 0 &&
      window.location.pathname === stripSuffix(import.meta.env.BASE_URL, '/')
    ) {
      window.location.assign(import.meta.env.BASE_URL);
    }
  }, []);

  return (
    <ThemeProvider theme={theme}>
      <AppContainer>
        <HookedModal />
        <ModalBetaGreeter />
        <Router basename={import.meta.env.BASE_URL}>
          <Routes>
            <Route path="start" element={<StartRoute />} />
            <Route path="login" element={<LoginRoute />} />
            <Route
              path="unbanme"
              element={
                <RouteSuspense>
                  <UnbanmeRoute />
                </RouteSuspense>
              }
            />
            <Route
              path="verify"
              element={
                <RouteSuspense>
                  <VerifyRoute />
                </RouteSuspense>
              }
            />

            <Route path="db" element={<DashboardRoute />}>
              <Route
                path="guilds/:guildid/members"
                element={
                  <RouteSuspense>
                    <GuildMembersRoute />
                  </RouteSuspense>
                }
              />
              <Route
                path="guilds/:guildid/members/:memberid"
                element={
                  <RouteSuspense>
                    <MemberRoute />
                  </RouteSuspense>
                }
              />
              <Route
                path="guilds/:guildid/modlog"
                element={
                  <RouteSuspense>
                    <GuildModlogRoute />
                  </RouteSuspense>
                }
              />
              <Route
                path="guilds/:guildid/settings/general"
                element={
                  <RouteSuspense>
                    <GuildGeneralRoute />
                  </RouteSuspense>
                }
              />
              <Route
                path="guilds/:guildid/settings/backups"
                element={
                  <RouteSuspense>
                    <GuildBackupsRoute />
                  </RouteSuspense>
                }
              />
              <Route
                path="guilds/:guildid/settings/antiraid"
                element={
                  <RouteSuspense>
                    <GuildAntiraidRoute />
                  </RouteSuspense>
                }
              />
              <Route
                path="guilds/:guildid/settings/codeexec"
                element={
                  <RouteSuspense>
                    <GuildCodeexecRoute />
                  </RouteSuspense>
                }
              />
              <Route
                path="guilds/:guildid/settings/verification"
                element={
                  <RouteSuspense>
                    <GuildVerificationRoute />
                  </RouteSuspense>
                }
              />
              <Route
                path="guilds/:guildid/settings/karma"
                element={
                  <RouteSuspense>
                    <GuildKarmaRoute />
                  </RouteSuspense>
                }
              />
              <Route
                path="guilds/:guildid/settings/logs"
                element={
                  <RouteSuspense>
                    <GuildLogsRoute />
                  </RouteSuspense>
                }
              />
              <Route
                path="guilds/:guildid/settings/data"
                element={
                  <RouteSuspense>
                    <GuildDataRoute />
                  </RouteSuspense>
                }
              />
              <Route
                path="guilds/:guildid/settings/api"
                element={
                  <RouteSuspense>
                    <GuildAPIRoute />
                  </RouteSuspense>
                }
              />
              <Route
                path="guilds/:guildid/settings/permissions"
                element={
                  <RouteSuspense>
                    <GuildPermissionsRoute />
                  </RouteSuspense>
                }
              />
              <Route
                path="guilds/:guildid/settings/linkblocking"
                element={
                  <RouteSuspense>
                    <GuildLinkBlockingRoute />
                  </RouteSuspense>
                }
              />
              {import.meta.env.DEV && <Route path="debug" element={<DebugRoute />} />}
            </Route>

            <Route
              path="usersettings"
              element={
                <RouteSuspense>
                  <UserSettingsRoute />
                </RouteSuspense>
              }>
              <Route
                path="apitoken"
                element={
                  <RouteSuspense>
                    <APITokenRoute />
                  </RouteSuspense>
                }
              />
              <Route
                path="ota"
                element={
                  <RouteSuspense>
                    <OTARoute />
                  </RouteSuspense>
                }
              />
              <Route
                path="privacy"
                element={
                  <RouteSuspense>
                    <PrivacyRoute />
                  </RouteSuspense>
                }
              />
              <Route path="" element={<Navigate to="apitoken" />} />
            </Route>

            <Route path="*" element={<Navigate to="db" />} />
          </Routes>
        </Router>
        <Notifications />
      </AppContainer>
      <GlobalStyle />
    </ThemeProvider>
  );
};

export default App;
