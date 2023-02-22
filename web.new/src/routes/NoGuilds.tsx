import { Trans, useTranslation } from 'react-i18next';

import { ReactComponent as ArrowIcon } from '../assets/back.svg';
import { Button } from '../components/Button';
import React from 'react';
import styled from 'styled-components';

type Props = {};

const Container = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  text-align: center;
  font-size: 1.4rem;
  line-height: 2.5rem;
  font-weight: 300;
  gap: 1.5em;
`;

const InviteButton = styled(Button)`
  > svg {
    transform: rotate(180deg);
  }
`;

const NoGuildsRoute: React.FC<Props> = ({}) => {
  const { t } = useTranslation('routes.noguilds');

  return (
    <Container>
      <Trans
        ns="routes.noguilds"
        i18nKey="info"
        components={{
          br: <br />,
        }}
      />
      <a href="/invite">
        <InviteButton>
          <ArrowIcon />
          <span>{t('invite')}</span>
        </InviteButton>
      </a>
    </Container>
  );
};

export default NoGuildsRoute;
