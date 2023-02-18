import { APIClient } from '../services/api';
import { APIError } from '../lib/shinpuru-ts/src/errors';
import { Client } from '../lib/shinpuru-ts/src';
import { useNavigate } from 'react-router';
import { useNotifications } from './useNotifications';

export const useApi = () => {
  const nav = useNavigate();
  const { pushNotification } = useNotifications();

  async function fetch<T>(
    req: (c: Client) => Promise<T>,
    silenceErrors: boolean = false,
  ): Promise<T> {
    try {
      return await req(APIClient);
    } catch (e) {
      if (!silenceErrors) {
        if (e instanceof APIError) {
          if (e.code === 401) {
            nav('/start');
          } else {
            pushNotification({
              type: 'ERROR',
              delay: 8000,
              heading: 'API Error',
              message: `${e.message} (${e.code})`,
            });
          }
        } else {
          pushNotification({
            type: 'ERROR',
            delay: 8000,
            heading: 'Error',
            message: `Unknown Request Error: ${e}`,
          });
        }
      }
      throw e;
    }
  }

  return fetch;
};
