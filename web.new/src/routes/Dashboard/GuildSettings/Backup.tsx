import React from 'react';
import { useTranslation } from 'react-i18next';
import { useParams } from 'react-router';
import { useTheme } from 'styled-components';

type Props = {};

const BackupRoute: React.FC<Props> = ({}) => {
  const { t } = useTranslation('routes.member');
  const { guildid } = useParams();
  const theme = useTheme();

  return (
    <>
      <i>soon™</i>
    </>
  );
};

export default BackupRoute;
