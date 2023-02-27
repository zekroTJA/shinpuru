import { Guild, PermissionsMap, Role } from '../../lib/shinpuru-ts/src';
import styled, { css, useTheme } from 'styled-components';
import { useMemo, useState } from 'react';

import { ReactComponent as AllowIcon } from '../../assets/check.svg';
import { AutocompleteInput } from '../AutocompleteInput';
import { Button } from '../Button';
import Color from 'color';
import { ReactComponent as DeleteIcon } from '../../assets/delete.svg';
import { ReactComponent as DisallowIcon } from '../../assets/ban.svg';
import { Flex } from '../Flex';
import { LinearGradient } from '../styleParts';
import { RoleInput } from '../RoleInput';
import { Switch } from '../Switch';
import { uid } from 'react-uid';
import { useApi } from '../../hooks/useApi';
import { useTranslation } from 'react-i18next';

type Props = {
  perms: PermissionsMap;
  setPerms: (v: PermissionsMap) => void;
  guild: Guild;
  available: string[];
};

const ControlContainer = styled(Flex)`
  gap: 1em;
  z-index: 5;
`;

const StyledSwitch = styled(Switch)`
  svg {
    width: 100%;
    height: 100%;
  }
`;

const StyledButton = styled(Button)`
  padding: 0 0.8em;
`;

const RoleContainer = styled.div<{ rColor: number }>`
  padding: 1em;
  border-radius: 12px;

  ${(p) =>
    css`
      background: linear-gradient(
        160deg,
        ${Color(p.rColor).alpha(0.3).hexa()},
        ${Color(p.rColor).alpha(0.1).hexa()}
      );
    `}

  > h3 {
    margin: 0 0 1em 0;
    color: ${(p) =>
      Color(!p.rColor ? '#6c6c6c' : p.rColor)
        .mix(Color(p.theme.text))
        .mix(Color(p.theme.text))
        .hexa()};
  }

  > div {
    display: flex;
    flex-direction: column;
    gap: 0.5em;
  }
`;

const PermissionEntry = styled.div<{ isAdditive: boolean }>`
  display: flex;
  align-items: center;
  cursor: pointer;

  > span {
    padding: 0.5em;
    width: fit-content;
    border-radius: 8px;
    z-index: 1;
    ${(p) => LinearGradient(p.isAdditive ? p.theme.green : p.theme.red)}
  }

  > button {
    border: none;
    background-color: ${(p) => p.theme.background3};
    color: ${(p) => p.theme.text};
    cursor: pointer;
    padding: 0.75em 0.75em 0.75em calc(0.75em + 5px);
    border-radius: 0 8px 8px 0;
    opacity: 0.2;
    position: relative;
    left: -10px;
    transition: all 0.2s ease;
  }

  &:hover > button {
    opacity: 1;
    left: -5px;
  }
`;

export const PermissionSelector: React.FC<Props> = ({ perms, setPerms, available, guild }) => {
  const theme = useTheme();
  const fetch = useApi();
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

  const applyRule = () => {
    const sign = allow ? '+' : '-';

    fetch((c) =>
      c.guilds.applyPermission(guild.id, {
        perm: sign + permission,
        role_ids: roles.map((r) => r.id),
        override: true,
      }),
    )
      .then((r) => setPerms(r))
      .catch();
  };

  const removeRule = (roleId: string, rule: string) => {
    const sign = rule[0] === '+' ? '-' : '+';

    fetch((c) =>
      c.guilds.applyPermission(guild.id, {
        perm: sign + rule.substring(1),
        role_ids: [roleId],
        override: false,
      }),
    )
      .then((r) => setPerms(r))
      .catch();
  };

  console.debug(permsKV);

  return (
    <Flex direction="column" gap="1em">
      {guild && (
        <RoleInput
          placeholder={t('placeholder.roles')}
          guild={guild}
          selected={roles}
          onChange={setRoles}
        />
      )}
      <ControlContainer gap="1em">
        <AutocompleteInput
          isInvalid={isInvalidPermission}
          value={permission}
          setValue={setPermission}
          selections={available}
          placeholder={t('placeholder.perms')}
        />
        <StyledSwitch
          enabled={allow}
          onChange={setAllow}
          theaming={{ enabledColor: theme.green, disabledColor: theme.red }}>
          {(allow && <AllowIcon style={{ color: theme.green }} />) || (
            <DisallowIcon style={{ color: theme.red }} />
          )}
        </StyledSwitch>
        <StyledButton
          disabled={roles.length === 0 || !permission || isInvalidPermission}
          onClick={applyRule}>
          {t('apply')}
        </StyledButton>
      </ControlContainer>
      {permsKV.map(([role, perm]) => (
        <RoleContainer key={uid(role)} rColor={role.color}>
          <h3>{role.name}</h3>
          <div>
            {perm.map((p) => (
              <PermissionEntry key={uid(p)} isAdditive={p.startsWith('+')}>
                <span>{p}</span>
                <button onClick={() => removeRule(role.id, p)}>
                  <DeleteIcon />
                </button>
              </PermissionEntry>
            ))}
          </div>
        </RoleContainer>
      ))}
    </Flex>
  );
};
