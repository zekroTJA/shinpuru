import { Client } from '../lib/shinpuru-ts/src';

const API_ENDPOINT = import.meta.env.PROD ? '/api' : 'http://localhost:8080/api';

export const loginRoute = (redirect?: string): string =>
  `${API_ENDPOINT}/auth/login${!!redirect ? '?redirect=' + redirect : ''}`;

export const APIClient = new Client(API_ENDPOINT);
