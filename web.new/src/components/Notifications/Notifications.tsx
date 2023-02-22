import { NotificationTile } from './Notification';
import styled from 'styled-components';
import { uid } from 'react-uid';
import { useNotifications } from '../../hooks/useNotifications';

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
