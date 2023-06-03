import './index.scss';
import './i18n';

import React, { Suspense } from 'react';

import App from './App';
import { createRoot } from 'react-dom/client';

const container = document.getElementById('root');
const root = createRoot(container!);
root.render(
  <React.StrictMode>
    {/* TODO: Use better fallback for language suspense */}
    <Suspense fallback="loading...">
      <App />
    </Suspense>
  </React.StrictMode>,
);
