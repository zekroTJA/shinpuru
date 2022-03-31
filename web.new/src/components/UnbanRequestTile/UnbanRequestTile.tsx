import { useTranslation } from 'react-i18next';
import styled from 'styled-components';
import { ReactComponent as HammerIcon } from '../../assets/hammer.svg';
import { ReactComponent as PrayIcon } from '../../assets/pray.svg';
import { UnbanRequest, UnbanRequestState } from '../../lib/shinpuru-ts/src';
import { Button } from '../Button';
import { Container } from '../Container';
import { Flex } from '../Flex';
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
  width: 100%;
  border-radius: 12px;
  text-align: center;
  text-transform: uppercase;
  padding: 0.4em;

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
`;

const Controls = styled.div`
  display: flex;
  gap: 1em;

  > * {
    width: 100%;
  }
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
        <Flex gap="1em">
          <UserTileSmall fallbackId={request.user_id} user={request.creator} icon={<PrayIcon />} />
          {request.processor && (
            <UserTileSmall
              fallbackId={request.processed_by}
              user={request.processor}
              icon={<HammerIcon />}
            />
          )}
        </Flex>
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
        {showControls && (
          <Controls>
            <Button variant="green" onClick={() => onProcess(true)}>
              {t('actions.accept')}
            </Button>
            <Button variant="red" onClick={() => onProcess(false)}>
              {t('actions.decline')}
            </Button>
          </Controls>
        )}
      </ContentContainer>
    </RequestContainer>
  );
};
