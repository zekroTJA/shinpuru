import { Element, Select } from '../components/Select';
import { Notification, NotificationType } from '../components/Notifications';
import { randomFrom, randomNumber } from '../util/rand';
import { useRef, useState } from 'react';

import { ActionButton } from '../components/ActionButton';
import { Button } from '../components/Button';
import { Heading } from '../components/Heading';
import { Input } from '../components/Input';
import styled from 'styled-components';
import { useModal } from '../hooks/useModal';
import { useNotifications } from '../hooks/useNotifications';

type Props = {};

const DebugContainer = styled.div`
  padding: 2em;
  > section {
    margin-bottom: 2em;
    > * {
      width: 100%;
      display: block;
      margin: 0 0.5em 0.5em 0.5em;
    }
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
  const type = ['INFO', 'WARN', 'ERROR', 'SUCCESS'][randomNumber(3, 0)];
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
  const NOTIFICATION_OPTIONS: Element<NotificationType>[] = [
    { id: 'INFO', display: 'INFO', value: 'INFO' },
    {
      id: 'SUCCESS',
      display: 'SUCCESS',
      value: 'SUCCESS',
    },
    {
      id: 'WARNING',
      display: 'WARNING',
      value: 'WARNING',
    },
    { id: 'ERROR', display: 'ERROR', value: 'ERROR' as NotificationType },
  ];
  const [notType, setNotType] = useState(NOTIFICATION_OPTIONS[0]);

  const refInptTitle = useRef<HTMLInputElement>(null);
  const refInptContent = useRef<HTMLInputElement>(null);
  const refInptDelay = useRef<HTMLInputElement>(null);

  const { pushNotification } = useNotifications();
  const { openModal } = useModal<number>();

  const _openModal = () => {
    openModal({
      content: 'heyo',
      heading: 'test modal',
      controls: [
        { name: 'ok', value: 1 },
        { name: 'cancel', value: 2, variant: 'gray' },
      ],
    }).then(console.log);
  };

  const timeoutPromise = () =>
    new Promise((resolve, _) => {
      setTimeout(resolve, 3000);
    });

  return (
    <DebugContainer>
      <section>
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
          }>
          Send
        </Button>
        <Button
          variant="blue"
          onClick={() =>
            pushNotification(
              getRandomNotification(
                refInptDelay.current?.value ? parseInt(refInptDelay.current?.value) : undefined,
              ),
            )
          }>
          Send Random
        </Button>
      </section>

      <section>
        <Heading>Modals</Heading>
        <Button onClick={_openModal}>Open</Button>
      </section>

      <section>
        <Heading>Buttons</Heading>
        <ActionButton onClick={timeoutPromise}>Some Action!</ActionButton>
      </section>
    </DebugContainer>
  );
};
