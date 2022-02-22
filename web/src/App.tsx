import { useStoredTheme } from 'hooks/useStoredTheme';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import { Start } from 'routes/Start';
import { createGlobalStyle, ThemeProvider } from 'styled-components';

const GlobalStyle = createGlobalStyle`
  body {
    background-color: ${(p) => p.theme.background};
    color: ${(p) => p.theme.text};
  }
`;

function App() {
  const { theme } = useStoredTheme();

  return (
    <ThemeProvider theme={theme}>
      <Router>
        <Routes>
          <Route path="start" element={<Start />} />
        </Routes>
      </Router>
      <GlobalStyle />
    </ThemeProvider>
  );
}

export default App;
