import { Trans, useTranslation } from 'react-i18next';

import { ReactComponent as ArrowIcon } from '../assets/back.svg';
import { BottomContainer } from '../components/Navbar';
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

const StyledBottomContainer = styled(BottomContainer)`
  margin: 0;
  width: 14rem;
  position: absolute;
  bottom: 2em;
  font-size: 1rem;
  text-align: start;
  font-weight: normal;
`;

const NoGuildsRoute: React.FC<Props> = () => {
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
      <StyledBottomContainer />
    </Container>
  );
};

export default NoGuildsRoute;
