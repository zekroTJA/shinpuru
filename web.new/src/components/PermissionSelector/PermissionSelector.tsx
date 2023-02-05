import { useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { uid } from 'react-uid';
import styled, { useTheme } from 'styled-components';
import { ReactComponent as DisallowIcon } from '../../assets/ban.svg';
import { ReactComponent as AllowIcon } from '../../assets/check.svg';
import { Guild, PermissionsMap, Role } from '../../lib/shinpuru-ts/src';
import { AutocompleteInput } from '../AutocompleteInput';
import { Button } from '../Button';
import { Flex } from '../Flex';
import { RoleInput } from '../RoleInput';
import { Switch } from '../Switch';

type Props = {
  perms: PermissionsMap;
  setPerms: (v: PermissionsMap) => void;
  guild: Guild;
  available: string[];
};

const StyledSwitch = styled(Switch)`
  svg {
    width: 100%;
    height: 100%;
  }
`;

const StyledButton = styled(Button)`
  padding: 0 0.8em;
`;

export const PermissionSelector: React.FC<Props> = ({ perms, setPerms, available, guild }) => {
  const theme = useTheme();
  const { t } = useTranslation('components', { keyPrefix: 'permissionselector' });
  const [roles, setRoles] = useState<Role[]>([]);
  const [allow, setAllow] = useState(false);
  const [permission, setPermission] = useState('');

  const isInvalidPermission =
    !!permission &&
    !available.find((a) => {
      if (permission.endsWith('.*'))
        return a.startsWith(permission.substring(0, permission.length - 2));
      return a === permission;
    });

  const permsKV = useMemo(
    () =>
      Object.keys(perms)
        .map((k) => [guild.roles?.find((r) => r.id === k), perms[k]] as [Role, string[]])
        .filter(([r, _]) => !!r),
    [guild, perms],
  );

  return (
    <Flex direction="column" gap="1em">
      {guild && <RoleInput guild={guild} selected={roles} onChange={setRoles} />}
      <Flex gap="1em">
        <StyledSwitch
          enabled={allow}
          onChange={setAllow}
          theaming={{ enabledColor: theme.green, disabledColor: theme.red }}>
          {(allow && <AllowIcon style={{ color: theme.green }} />) || (
            <DisallowIcon style={{ color: theme.red }} />
          )}
        </StyledSwitch>
        <AutocompleteInput
          isInvalid={isInvalidPermission}
          value={permission}
          setValue={setPermission}
          selections={available}
          placeholder={t('placeholder')}
        />
        <StyledButton disabled={roles.length === 0 || !permission || isInvalidPermission}>
          {t('apply')}
        </StyledButton>
      </Flex>
      {permsKV.map(([role, perm]) => (
        <div key={uid(role)}>
          <span>{role.name}</span>
          {perm.map((p) => (
            <span>{p}</span>
          ))}
        </div>
      ))}
    </Flex>
  );
};
