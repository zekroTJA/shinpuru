import { Client } from '../lib/shinpuru-ts/src';

const API_ENDPOINT = import.meta.env.PROD
  ? '/api'
  : 'http://localhost:8080/api';

export const LOGIN_ROUTE = API_ENDPOINT + '/auth/login';

export const APIClient = new Client(API_ENDPOINT);
