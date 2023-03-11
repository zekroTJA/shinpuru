import { useEffect, useState } from 'react';

import { Flex } from '../../../components/Flex';
import { Loader } from '../../../components/Loader';
import { MaxWidthContainer } from '../../../components/MaxWidthContainer';
import { PermissionSelector } from '../../../components/PermissionSelector';
import { PermissionsMap } from '../../../lib/shinpuru-ts/src';
import { Small } from '../../../components/Small';
import { useApi } from '../../../hooks/useApi';
import { useGuild } from '../../../hooks/useGuild';
import { useParams } from 'react-router';
import { useSelfUser } from '../../../hooks/useSelfUser';
import { useTranslation } from 'react-i18next';

type Props = {};

const PermissionsRoute: React.FC<Props> = () => {
  const { t } = useTranslation('routes.guildsettings.permissions');
  const { guildid } = useParams();
  const guild = useGuild(guildid);
  const fetch = useApi();
  const me = useSelfUser();

  const [perms, setPerms] = useState<PermissionsMap>();
  const [allowed, setAllowed] = useState<string[]>();

  useEffect(() => {
    if (!guild || !me) return;

    fetch((c) => c.guilds.permissions(guild.id))
      .then(setPerms)
      .catch();

    fetch((c) => c.etc.allpermissions())
      .then((r) => {
        const rules = r.data.filter(
          (p) => p.startsWith('sp.guild.') || p.startsWith('sp.etc.') || p.startsWith('sp.chat.'),
        );
        setAllowed(rules);
      })
      .catch();
  }, [guild, me]);

  return (
    <MaxWidthContainer>
      <h1>{t('heading')}</h1>
      <Small>{t('explanation')}</Small>
      {(guild && perms && allowed && (
        <PermissionSelector guild={guild} available={allowed} perms={perms} setPerms={setPerms} />
      )) || (
        <>
          <Loader width="100%" height="2.8em" />
          <Flex gap="1em">
            <Loader width="100%" height="2.5em" margin="1em 0 0 0" />
            <Loader width="5em" height="2.5em" margin="1em 0 0 0" />
            <Loader width="5em" height="2.5em" margin="1em 0 0 0" />
          </Flex>
          <Loader width="100%" height="9em" margin="1em 0 0 0" />
          <Loader width="100%" height="6em" margin="1em 0 0 0" />
        </>
      )}
    </MaxWidthContainer>
  );
};

export default PermissionsRoute;
