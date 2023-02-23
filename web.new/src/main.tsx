import './index.scss';
import './i18n';

import React, { Suspense } from 'react';

import App from './App';
import ReactDOM from 'react-dom';

ReactDOM.render(
  <React.StrictMode>
    {/* TODO: Use better fallback for language suspense */}
    <Suspense fallback="loading...">
      <App />
    </Suspense>
  </React.StrictMode>,
  document.getElementById('root'),
);
