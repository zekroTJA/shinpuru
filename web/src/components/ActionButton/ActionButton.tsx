import { Button, ButtonProps } from '../Button';
import React, { PropsWithChildren, useState } from 'react';

import { BeatLoader } from 'react-spinners';
import { PropsWithStyle } from '../props';
import { useTheme } from 'styled-components';

type Props = PropsWithChildren &
  ButtonProps &
  PropsWithStyle & {
    onClick: (e: React.MouseEvent<HTMLButtonElement, MouseEvent>) => Promise<any> | undefined;
    disabled?: boolean;
  };

export const ActionButton: React.FC<Props> = ({ onClick, children, disabled, ...props }) => {
  const theme = useTheme();
  const [loading, setLoading] = useState(false);

  const _onClick = (e: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    const prom = onClick(e);
    if (!prom) return;
    setLoading(true);
    prom.finally(() => setLoading(false));
  };

  return (
    <Button onClick={_onClick} disabled={disabled ?? loading} {...props}>
      {(loading && <BeatLoader color={theme.textAlt} margin={0} />) || children}
    </Button>
  );
};
