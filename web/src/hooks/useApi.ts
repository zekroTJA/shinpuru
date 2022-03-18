import { useNavigate } from 'react-router';
import { Client } from '../lib/shinpuru-ts/src';
import { APIError } from '../lib/shinpuru-ts/src/errors';
import { APIClient } from '../services/api';

export const useApi = () => {
  const nav = useNavigate();

  async function fetch<T>(req: (c: Client) => Promise<T>): Promise<T> {
    try {
      return await req(APIClient);
    } catch (e) {
      if (e instanceof APIError) {
        if (e.code === 401) {
          nav('/start');
        }
      }
      throw e;
    }
  }

  return fetch;
};
