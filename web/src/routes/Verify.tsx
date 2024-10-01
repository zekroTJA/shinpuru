import styled, { useTheme } from 'styled-components';

import { APIError } from '../lib/shinpuru-ts/src/errors';
import { Button } from '../components/Button';
import { Card } from '../components/Card';
import { ReactComponent as CheckIcon } from '../assets/check.svg';
import HCaptcha from '@hcaptcha/react-hcaptcha';
import { Loader } from '../components/Loader';
import { ReactComponent as VerificationIcon } from '../assets/verification.svg';
import { useApi } from '../hooks/useApi';
import { useEffectAsync } from '../hooks/useEffectAsync';
import { useNavigate } from 'react-router';
import { useNotifications } from '../hooks/useNotifications';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';

type Props = {};

const StyledCard = styled(Card)`
  display: flex;
  flex-direction: column;
  gap: 1em;
  align-items: center;
  justify-content: center;
  font-size: 1.2em;
  max-width: 18em;
  text-align: center;
  line-height: 1.6em;

  > svg {
    width: 8em;
    height: 8em;
  }

  > ${Button} {
    width: 100%;
  }
`;

const Container = styled.div`
  display: flex;
  flex-direction: column;
  gap: 1em;
  padding: 2em;
  width: 100%;
  justify-content: center;
  align-items: center;
`;

const VerificationRoute: React.FC<Props> = () => {
  const { t } = useTranslation('routes.verification');
  const nav = useNavigate();
  const fetch = useApi();
  const { pushNotification } = useNotifications();
  const theme = useTheme();

  const [authorized, setAuthorized] = useState<boolean>();
  const [verified, setVerified] = useState<boolean>();
  const [sitekey, setSitekey] = useState<string>();

  const _onLogin = () => {
    nav({ pathname: '/login', search: 'redirect=verify' });
  };

  const _onVerify = (token: string, ekey: string) => {
    fetch((c) => c.verification.verify(token))
      .then(() => {
        setVerified(true);
        pushNotification({ message: t('notifications.verified'), type: 'SUCCESS' });
      })
      .catch();
  };

  useEffectAsync(async () => {
    try {
      // await fetch((c) => c.auth.check(), true);
      setVerified((await fetch((c) => c.etc.me(), true)).captcha_verified);
      setAuthorized(true);
    } catch (e) {
      if (e instanceof APIError && e.code === 401) {
        setAuthorized(false);
      }
      return;
    }

    const sitekey = await fetch((c) => c.verification.sitekey());
    setSitekey(sitekey.sitekey);
  }, []);

  return (
    <Container>
      {authorized === undefined && (
        <>
          <Loader />
        </>
      )}
      {authorized === false && (
        <Container>
          <span>{t('notloggedin')}</span>
          <Button onClick={_onLogin}>{t('login')}</Button>
        </Container>
      )}
      {authorized === true && (
        <Container>
          {(verified && (
            <StyledCard color={theme.lime}>
              <CheckIcon />
              <div>{t('verified')}</div>
              <Button variant="green" onClick={() => nav('/db')}>
                {t('continue')}
              </Button>
            </StyledCard>
          )) || (
            <StyledCard color={theme.orange}>
              <VerificationIcon />
              <span>{t('needverify')}</span>
              {(sitekey && <HCaptcha sitekey={sitekey} onVerify={_onVerify} />) || (
                <Loader width="20em" height="5em" />
              )}
            </StyledCard>
          )}
        </Container>
      )}
    </Container>
  );
};

export default VerificationRoute;
