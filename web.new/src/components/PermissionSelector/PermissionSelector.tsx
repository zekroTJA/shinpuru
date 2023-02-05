import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import styled, { useTheme } from 'styled-components';
import { ReactComponent as DisallowIcon } from '../../assets/ban.svg';
import { ReactComponent as AllowIcon } from '../../assets/check.svg';
import { PermissionsMap, Role } from '../../lib/shinpuru-ts/src';
import { AutocompleteInput } from '../AutocompleteInput';
import { Button } from '../Button';
import { Flex } from '../Flex';
import { Switch } from '../Switch';

type Props = {
  perms: PermissionsMap;
  setPerms: (v: PermissionsMap) => void;
  roles: Role[];
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

export const PermissionSelector: React.FC<Props> = ({ perms, setPerms, available }) => {
  const theme = useTheme();
  const { t } = useTranslation('components', { keyPrefix: 'permissionselector' });
  const [permission, setPermission] = useState('');
  const [allow, setAllow] = useState(false);

  return (
    <>
      <Flex gap="1em">
        <StyledSwitch
          enabled={allow}
          onChange={setAllow}
          theaming={{ enabledColor: theme.green, disabledColor: theme.red }}>
          {(allow && <AllowIcon style={{ color: theme.green }} />) || (
            <DisallowIcon style={{ color: theme.red }} />
          )}
        </StyledSwitch>
        <AutocompleteInput value={permission} setValue={setPermission} selections={available} />
        <StyledButton>{t('apply')}</StyledButton>
      </Flex>
    </>
  );
};
