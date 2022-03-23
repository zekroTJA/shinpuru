import { useNavigate } from 'react-router';
import { NotificationType } from '../components/Notifications';
import { Client } from '../lib/shinpuru-ts/src';
import { APIError } from '../lib/shinpuru-ts/src/errors';
import { APIClient } from '../services/api';
import { useNotifications } from './useNotifications';

export const useApi = () => {
  const nav = useNavigate();
  const { pushNotification } = useNotifications();

  async function fetch<T>(req: (c: Client) => Promise<T>): Promise<T> {
    try {
      return await req(APIClient);
    } catch (e) {
      if (e instanceof APIError) {
        if (e.code === 401) {
          nav('/start');
        } else if (e.code !== 410) {
          pushNotification({
            type: NotificationType.ERROR,
            delay: 8000,
            heading: 'API Error',
            message: `${e.message} (${e.code})`,
          });
        }
      } else {
        pushNotification({
          type: NotificationType.ERROR,
          delay: 8000,
          heading: 'Error',
          message: `Unknown Request Error: ${e}`,
        });
      }
      throw e;
    }
  }

  return fetch;
};
