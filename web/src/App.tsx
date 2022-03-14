import { useStoredTheme } from './hooks/useStoredTheme';
import {
  BrowserRouter as Router,
  Navigate,
  Route,
  Routes,
} from 'react-router-dom';
import { StartRoute } from './routes/Start';
import styled, { createGlobalStyle, ThemeProvider } from 'styled-components';
import { LoginRoute } from './routes/Login';
import { DashboardRoute } from './routes/Dashboard';
import { Guild } from './routes/Dashboard/Guild';

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

  return (
    <ThemeProvider theme={theme}>
      <AppContainer>
        <Router>
          <Routes>
            <Route path="start" element={<StartRoute />} />
            <Route path="login" element={<LoginRoute />} />
            <Route path="db" element={<DashboardRoute />}>
              <Route path="guilds/:guildid" element={<Guild />} />
            </Route>
            <Route path="*" element={<Navigate to="db" />} />
          </Routes>
        </Router>
      </AppContainer>
      <GlobalStyle />
    </ThemeProvider>
  );
};

export default App;
