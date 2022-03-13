import { useStoredTheme } from './hooks/useStoredTheme';
import { useEffect } from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import { StartRoute } from './routes/Start';
import { APIClient } from './services/api';
import { createGlobalStyle, ThemeProvider } from 'styled-components';

const GlobalStyle = createGlobalStyle`
  body {
    background-color: ${(p) => p.theme.background};
    color: ${(p) => p.theme.text};
  }
`;

export const App: React.FC = () => {
  const { theme } = useStoredTheme();

  useEffect(() => {
    APIClient.etc.sysinfo().then(console.log).catch(console.error);
  }, []);

  return (
    <ThemeProvider theme={theme}>
      <p>poggers</p>
      <Router>
        <Routes>
          <Route path="start" element={<StartRoute />} />
        </Routes>
      </Router>
      <GlobalStyle />
    </ThemeProvider>
  );
};

export default App;
