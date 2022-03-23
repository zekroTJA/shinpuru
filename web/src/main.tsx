import React, { Suspense } from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import './index.scss';

import './i18n';

ReactDOM.render(
  <React.StrictMode>
    {/* TODO: Use better fallback for language suspense */}
    <Suspense fallback="loading...">
      <App />
    </Suspense>
  </React.StrictMode>,
  document.getElementById('root')
);
