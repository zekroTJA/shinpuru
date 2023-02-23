import { Notification } from '../components/Notifications';
import { randomNumber } from '../util/rand';
import { useStore } from '../services/store';

const DEFAULT_DELAY = 6000;
export const ANIMATION_DELAY = 500;

export const useNotifications = () => {
  const [notifications, setNotifications] = useStore((s) => [s.notifications, s.setNotifications]);

  const _push = (n: Notification) => {
    console.log('push notification', n);
    notifications.push(n);
    setNotifications(notifications);
  };

  const _remove = (n: Notification) => {
    n.hide = true;
    setNotifications(notifications);
    setTimeout(() => {
      const i = notifications.findIndex((_n) => _n.uid === n.uid);
      if (i !== -1) notifications.splice(i, 1);
      setNotifications(notifications);
    }, ANIMATION_DELAY);
  };

  const pushNotification = (notification: Notification) => {
    notification.delay = notification.delay ?? DEFAULT_DELAY;
    notification.uid = randomNumber(1000);
    _push(notification);
    if (notification.delay > 0) {
      setTimeout(() => _remove(notification), notification.delay + ANIMATION_DELAY * 2);
    }
    return () => _remove(notification);
  };

  const hideNotification = (notification: Notification) => _remove(notification);

  return { pushNotification, hideNotification, notifications };
};
