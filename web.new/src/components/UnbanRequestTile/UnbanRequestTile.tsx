import { useTranslation } from 'react-i18next';
import styled from 'styled-components';
import { ReactComponent as HammerIcon } from '../../assets/hammer.svg';
import { ReactComponent as PrayIcon } from '../../assets/pray.svg';
import { UnbanRequest, UnbanRequestState } from '../../lib/shinpuru-ts/src';
import { formatDate } from '../../util/date';
import { Button } from '../Button';
import { Container } from '../Container';
import { Embed } from '../Embed';
import { Heading } from '../Heading';
import { UserTileSmall } from '../UserTileSmall';
import { LinearGradient } from '../styleParts';

type Props = {
  request: UnbanRequest;
  showControls?: boolean;
  onProcess?: (accepted: boolean) => void;
};

const STATUS = {
  [UnbanRequestState.PENDING]: 'pending',
  [UnbanRequestState.DECLINED]: 'declined',
  [UnbanRequestState.ACCEPTED]: 'accepted',
};

const RequestContainer = styled(Container)`
  padding: 0;
  background-color: ${(p) => p.theme.background3};
`;

const StatusBar = styled.div<{ state: UnbanRequestState }>`
  font-size: 0.8rem;
  letter-spacing: 0.2ch;
  text-transform: uppercase;
  width: 100%;
  text-align: center;
  padding: 0.4em;
  border-radius: 8px;
  color: ${(p) => p.theme.background2};

  ${(p) => {
    switch (p.state) {
      case UnbanRequestState.ACCEPTED:
        return LinearGradient(p.theme.green);
      case UnbanRequestState.DECLINED:
        return LinearGradient(p.theme.red);
      default:
        return LinearGradient(p.theme.blurple);
    }
  }}
`;

const ContentContainer = styled.div`
  display: flex;
  flex-direction: column;
  padding: 1em;
  gap: 1em;

  ${Heading} {
    font-size: 0.8rem;
  }
`;

const Controls = styled.div`
  display: flex;
  gap: 1em;

  > * {
    width: 100%;
  }
`;

const UserConatienr = styled.div`
  display: flex;
  gap: 1.5em;
  justify-content: space-between;
  flex-wrap: wrap;

  @media (max-width: 50em) {
    flex-direction: column;
  }
`;

const Footer = styled.section`
  display: flex;
  align-items: center;
  border-top: solid 2px rgb(128, 128, 128);
  padding-top: 0.8em;

  font-size: 0.8rem;

  ${Embed} {
    font-size: 0.8em;
  }

  > a {
    text-decoration: underline;
    color: ${(p) => p.theme.accent};
    cursor: pointer;
  }
`;

const Spacer = styled.div`
  height: 1em;
  width: 1px;
  background-color: ${(p) => p.theme.text};
  margin: 0 0.5em;
`;

export const UnbanRequestTile: React.FC<Props> = ({
  request,
  showControls,
  onProcess = () => {},
}) => {
  const { t } = useTranslation('components', { keyPrefix: 'unbanrequesttile' });

  return (
    <RequestContainer>
      <StatusBar state={request.status}>{t(`state.${STATUS[request.status]}`)}</StatusBar>
      <ContentContainer>
        <UserConatienr>
          <UserTileSmall fallbackId={request.user_id} user={request.creator} icon={<PrayIcon />} />
          {request.processor && (
            <UserTileSmall
              fallbackId={request.processed_by}
              user={request.processor}
              icon={<HammerIcon />}
            />
          )}
        </UserConatienr>
        <section>
          <Heading>{t('requestermessage')}</Heading>
          <span>{request.message}</span>
        </section>
        {request.processed_by && (
          <section>
            <Heading>{t('decisionreason')}</Heading>
            <span>{request.message}</span>
          </section>
        )}
        {showControls && request.status === UnbanRequestState.PENDING && (
          <Controls>
            <Button variant="green" onClick={() => onProcess(true)}>
              {t('actions.accept')}
            </Button>
            <Button variant="red" onClick={() => onProcess(false)}>
              {t('actions.decline')}
            </Button>
          </Controls>
        )}
        <Footer>
          <span>
            ID: <Embed>{request.id}</Embed>
          </span>
          <Spacer />
          <span>{formatDate(request.created)}</span>
        </Footer>
      </ContentContainer>
    </RequestContainer>
  );
};
