import { tr } from 'date-fns/locale';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import styled from 'styled-components';
import { Button } from '../components/Button';
import { Flex } from '../components/Flex';
import { GuildTile } from '../components/GuiltTile';
import { Heading } from '../components/Heading';
import { Loader } from '../components/Loader';
import { ModalUnbanRequest } from '../components/Modals/ModalUnbanRequest';
import { NotificationType } from '../components/Notifications';
import { UnbanRequestTile } from '../components/UnbanRequestTile';
import { useApi } from '../hooks/useApi';
import { useEffectAsync } from '../hooks/useEffectAsync';
import { setInitRedirect } from '../hooks/useInitRedirect';
import { useNotifications } from '../hooks/useNotifications';
import { Guild, UnbanRequest } from '../lib/shinpuru-ts/src';
import { APIError } from '../lib/shinpuru-ts/src/errors';

type Props = {};

const Container = styled.div`
  display: flex;
  flex-direction: column;
  gap: 1em;
  padding: 2em;

  > span {
    display: block;
  }

  > ${Button} {
    width: fit-content;
  }
`;

const UnbanmeRoute: React.FC<Props> = () => {
  const { t } = useTranslation('routes.unbanme');
  const nav = useNavigate();
  const fetch = useApi();
  const { pushNotification } = useNotifications();

  const [authorized, setAuthorized] = useState<boolean>();
  const [bannedGuilds, setBannedGuilds] = useState<Guild[]>();
  const [queue, setQueue] = useState<UnbanRequest[]>([]);
  const [showUnbanModal, setShowUnbanModal] = useState<boolean>(false);
  const [selectedGuild, setSelectedGuild] = useState<Guild>();

  useEffectAsync(async () => {
    try {
      await fetch((c) => c.auth.check(), true);
      setAuthorized(true);
    } catch (e) {
      if (e instanceof APIError && e.code === 401) {
        setAuthorized(false);
      }
      return;
    }

    const queue = await fetch((c) => c.unbanrequests.list());
    setQueue(queue.data);

    const g = await fetch((c) => c.unbanrequests.guilds());
    setBannedGuilds(g.data.filter((g) => !queue.data.find((r) => r.guild_id === g.id)));
  }, []);

  const _onLogin = () => {
    setInitRedirect('/unbanme');
    nav('/login');
  };

  const _onSelectGuild = (guild: Guild) => {
    setSelectedGuild(guild);
    setShowUnbanModal(true);
  };

  const _onSend = (message: string) => {
    if (!selectedGuild || !message) return;

    fetch((c) =>
      c.unbanrequests.create({
        guild_id: selectedGuild.id,
        message,
      } as UnbanRequest),
    )
      .then((r) => {
        setQueue([r, ...queue]);
        pushNotification({
          message: t('notifications.sent'),
          type: NotificationType.SUCCESS,
        });
      })
      .catch();
  };

  const noguilds = queue.length === 0 ? <i>{t('nobannedguilds')}</i> : <i>{t('noopenguilds')}</i>;

  return (
    <Container>
      <ModalUnbanRequest
        show={showUnbanModal}
        guild={selectedGuild}
        onClose={() => setShowUnbanModal(false)}
        onSend={_onSend}
      />

      {authorized === undefined && (
        <>
          <Loader />
        </>
      )}
      {authorized === false && (
        <Container>
          <span>{t('notloggedin')}</span>
          <Button onClick={_onLogin}>{t('login')}</Button>
        </Container>
      )}
      {authorized === true && (
        <Flex direction="column" gap="3em">
          <section>
            <Heading>{t('applicableguilds')}</Heading>
            {(bannedGuilds === undefined && <Loader height="4em" width="25em" />) ||
              (bannedGuilds?.length === 0 && noguilds) || (
                <Flex gap="2em">
                  {bannedGuilds?.map((g) => (
                    <GuildTile guild={g} onClick={() => _onSelectGuild(g)} />
                  ))}
                </Flex>
              )}
          </section>

          {queue.length > 0 && (
            <section>
              <Heading>{t('requestqueue')}</Heading>
              {queue.map((r) => (
                <UnbanRequestTile request={r} key={r.id} />
              ))}
            </section>
          )}
        </Flex>
      )}
    </Container>
  );
};

export default UnbanmeRoute;
