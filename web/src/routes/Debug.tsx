import { useRef, useState } from 'react';
import styled from 'styled-components';
import { Button } from '../components/Button';
import { Heading } from '../components/Heading';
import { Input } from '../components/Input';
import { Notification, NotificationType } from '../components/Notifications';
import { Element, Select } from '../components/Select';
import { useNotifications } from '../hooks/useNotifications';
import { randomFrom, randomNumber } from '../util/rand';

interface Props {}

const DebugContainer = styled.div`
  > * {
    width: 100%;
    display: block;
    margin: 0 0.5em 0.5em 0.5em;
  }
`;

const getRandomNotification = (delay?: number) => {
  const randomWords = [
    'foo',
    'bar',
    'pog',
    'champ',
    'kek',
    'kekw',
    'poggers',
    'omegalul',
    'lul',
    'lulw',
  ];
  const heading = new Array(randomNumber(5, 2))
    .fill(null)
    .map(() => randomFrom(randomWords))
    .join(' ');
  const message = new Array(randomNumber(30, 5))
    .fill(null)
    .map(() => randomFrom(randomWords))
    .join(' ');
  const type = randomNumber(NotificationType.ERROR);
  return {
    heading,
    message,
    type,
    delay,
  } as Notification;
};

// This component will not be exported in the final
// production build!
export const DebugRoute: React.FC<Props> = () => {
  const NOTIFICATION_OPTIONS = [
    { id: 'INFO', display: 'INFO', value: NotificationType.INFO },
    {
      id: 'SUCCESS',
      display: 'SUCCESS',
      value: NotificationType.SUCCESS,
    },
    {
      id: 'WARNING',
      display: 'WARNING',
      value: NotificationType.WARNING,
    },
    { id: 'ERROR', display: 'ERROR', value: NotificationType.ERROR },
  ];
  const [notType, setNotType] = useState<Element<NotificationType>>(
    NOTIFICATION_OPTIONS[0]
  );

  const refInptTitle = useRef<HTMLInputElement>(null);
  const refInptContent = useRef<HTMLInputElement>(null);
  const refInptDelay = useRef<HTMLInputElement>(null);

  const { pushNotification } = useNotifications();

  return (
    <DebugContainer>
      <Heading>Notifications</Heading>
      <Input ref={refInptTitle} placeholder="title" />
      <Input ref={refInptContent} placeholder="content" />
      <Select
        value={notType}
        onElementSelect={(e) => setNotType(e)}
        options={NOTIFICATION_OPTIONS}
      />
      <Input ref={refInptDelay} placeholder="delay" type="number" min="0" />
      <Button
        onClick={() =>
          pushNotification({
            heading: refInptTitle.current?.value,
            message: refInptContent.current?.value ?? 'empty message',
            type: notType.value,
            delay: refInptDelay.current?.value
              ? parseInt(refInptDelay.current?.value)
              : undefined,
          })
        }
      >
        Send
      </Button>
      <Button
        variant="blue"
        onClick={() =>
          pushNotification(
            getRandomNotification(
              refInptDelay.current?.value
                ? parseInt(refInptDelay.current?.value)
                : undefined
            )
          )
        }
      >
        Send Random
      </Button>
    </DebugContainer>
  );
};
