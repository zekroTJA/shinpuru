import { BrowserRouter as Router } from 'react-router-dom';
import { createGlobalStyle } from 'styled-components';
import { Client } from 'lib/shinpuru-ts/src';

const GlobalStyle = createGlobalStyle``;

const client = new Client('http://localhost:8080/api');

function App() {
  const test = async () => {
    console.log(await client.req('GET', 'me'));
    console.log(client.util.color('abc', 4));
  };

  return (
    <div>
      <button onClick={() => test()}>TEST</button>
      <Router></Router>
      <GlobalStyle />
    </div>
  );
}

export default App;
