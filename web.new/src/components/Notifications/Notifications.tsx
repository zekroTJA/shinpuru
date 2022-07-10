import { uid } from 'react-uid';
import styled from 'styled-components';
import { useNotifications } from '../../hooks/useNotifications';
import { NotificationTile } from './Notification';

type Props = {};

const NotificationsContainer = styled.div`
  position: fixed;
  top: 0;
  right: 0;
  display: flex;
  flex-direction: column;
  gap: 1em;
  padding: 1em;
`;

export const Notifications: React.FC<Props> = () => {
  const { notifications, hideNotification } = useNotifications();
  return (
    <NotificationsContainer>
      {notifications.map((n) => (
        <NotificationTile key={uid(n)} n={n} onHide={() => hideNotification(n)} />
      ))}
    </NotificationsContainer>
  );
};
