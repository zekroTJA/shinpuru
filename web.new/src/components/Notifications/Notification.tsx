import { Notification, NotificationMeta, NotificationType } from './models';
import styled, { Keyframes, css, keyframes } from 'styled-components';

import { ANIMATION_DELAY } from '../../hooks/useNotifications';
import { ReactComponent as CloseIcon } from '../../assets/close.svg';
import { Container } from '../Container';
import { Heading } from '../Heading';
import { LinearGradient } from '../styleParts';

type Props = {
  n: Notification;
  onHide: () => void;
};

const AnimateInKF = keyframes`
  from {
    transform: translateX(10em);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
`;

const Amimate = (animation: Keyframes, delay: number) => css`
  animation: ${animation} ${ANIMATION_DELAY}ms ease ${delay}ms;
`;

const NotificationContainer = styled(Container)<NotificationMeta>`
  position: relative;
  min-width: 15em;
  width: 90vw;
  max-width: 25em;

  display: flex;
  justify-content: space-between;

  > div > svg {
    cursor: pointer;
    opacity: 0.6;
    transition: opacity 0.2s ease;
    margin-left: 0.2em;

    &:hover {
      opacity: 1;
    }
  }

  ${Amimate(AnimateInKF, 0)}
  ${(p) =>
    p.hide
      ? css`
          transform: translateX(10em);
          opacity: 0;
        `
      : ''}

  transition: all ${ANIMATION_DELAY}ms ease;

  ${(p) => {
    switch (p.type) {
      case 'ERROR':
        return LinearGradient(p.theme.red);
      case 'WARNING':
        return LinearGradient(p.theme.orange);
      case 'SUCCESS':
        return LinearGradient(p.theme.green);
      default:
        return LinearGradient(p.theme.blurple);
    }
  }}
`;

export const NotificationTile: React.FC<Props> = ({ n, onHide }) => {
  return (
    <NotificationContainer {...n}>
      <div>
        {n.heading && <Heading>{n.heading}</Heading>}
        <span>{n.message}</span>
      </div>
      <div>
        <CloseIcon onClick={onHide} />
      </div>
    </NotificationContainer>
  );
};
