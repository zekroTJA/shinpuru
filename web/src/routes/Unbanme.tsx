import { Guild, UnbanRequest, UnbanRequestState } from '../lib/shinpuru-ts/src';

import { APIError } from '../lib/shinpuru-ts/src/errors';
import { Button } from '../components/Button';
import { Flex } from '../components/Flex';
import { GuildTile } from '../components/GuiltTile';
import { Heading } from '../components/Heading';
import { Loader } from '../components/Loader';
import { ModalUnbanRequest } from '../components/Modals/ModalUnbanRequest';
import { UnbanRequestTile } from '../components/UnbanRequestTile';
import styled from 'styled-components';
import { useApi } from '../hooks/useApi';
import { useEffectAsync } from '../hooks/useEffectAsync';
import { useNavigate } from 'react-router-dom';
import { useNotifications } from '../hooks/useNotifications';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';

type Props = {};

const Container = styled.div`
  display: flex;
  flex-direction: column;
  gap: 1em;
  padding: 2em;
  align-items: center;
  padding-top: 5em;

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
    setQueue(queue.data.reverse());

    const g = await fetch((c) => c.unbanrequests.guilds());
    setBannedGuilds(
      g.data.filter(
        (g) =>
          !queue.data.find((r) => r.guild_id === g.id && r.status === UnbanRequestState.PENDING),
      ),
    );
  }, []);

  const _onLogin = () => {
    nav({ pathname: '/login', search: 'redirect=unbanme' });
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
        setBannedGuilds(bannedGuilds?.filter((g) => g.id !== selectedGuild.id));
        pushNotification({
          message: t<string>('notifications.sent'),
          type: 'SUCCESS',
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
          <Loader width="30em" height="5em" />
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
            <Flex direction="column" gap="1em">
              <Heading>{t('requestqueue')}</Heading>
              {queue.map((r) => (
                <UnbanRequestTile request={r} key={r.id} />
              ))}
            </Flex>
          )}
        </Flex>
      )}
    </Container>
  );
};

export default UnbanmeRoute;
