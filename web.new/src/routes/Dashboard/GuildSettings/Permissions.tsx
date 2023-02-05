import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useParams } from 'react-router';
import { MaxWidthContainer } from '../../../components/MaxWidthContainer';
import { PermissionSelector } from '../../../components/PermissionSelector';
import { Small } from '../../../components/Small';
import { useApi } from '../../../hooks/useApi';
import { useGuild } from '../../../hooks/useGuild';
import { useNotifications } from '../../../hooks/useNotifications';
import { useSelfUser } from '../../../hooks/useSelfUser';
import { PermissionsMap } from '../../../lib/shinpuru-ts/src';

type Props = {};

const PermissionsRoute: React.FC<Props> = ({}) => {
  const { t } = useTranslation('routes.guildsettings.permissions');
  const { pushNotification } = useNotifications();
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

    // TODO: Maybe add new route to get all available permissions
    //       instead of depending on allowed permissions for current user.
    fetch((c) => c.guilds.member(guild.id, me.id).permissionsAllowed())
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
      {guild && perms && allowed && (
        <PermissionSelector guild={guild} available={allowed} perms={perms} setPerms={setPerms} />
      )}
    </MaxWidthContainer>
  );
};

export default PermissionsRoute;
