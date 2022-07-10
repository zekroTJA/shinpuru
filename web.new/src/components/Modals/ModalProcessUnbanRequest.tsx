import { useEffect, useState } from 'react';
import { Trans, useTranslation } from 'react-i18next';
import styled from 'styled-components';
import { useApi } from '../../hooks/useApi';
import { useNotifications } from '../../hooks/useNotifications';
import { Report, UnbanRequest, UnbanRequestState } from '../../lib/shinpuru-ts/src';
import { Button } from '../Button';
import { Heading } from '../Heading';
import { Loader } from '../Loader';
import { ControlProps, Modal } from '../Modal/Modal';
import { NotificationType } from '../Notifications';
import { ReportTile } from '../Report';
import { ModalContainer, ModalTextArea } from './modalParts';

export type UnbanRequestWrapper = UnbanRequest & {
  isAccept: boolean;
};

type Props = ControlProps & {
  request?: UnbanRequestWrapper;
  onProcessed?: (request: UnbanRequest) => void;
};

const StyledModalContainer = styled(ModalContainer)`
  width: 60em;
`;

const ReportsContainer = styled.div`
  overflow-y: auto;
  max-height: calc(90vh - 30em);
  > div {
    display: flex;
    flex-direction: column;
    gap: 1em;
  }
`;

export const ModalProcessUnbanRequest: React.FC<Props> = ({
  show,
  request,
  onClose = () => {},
  onProcessed = () => {},
}) => {
  const { t } = useTranslation('components', { keyPrefix: 'modalprocessunbanrequest' });
  const fetch = useApi();
  const { pushNotification } = useNotifications();

  const [reason, setReason] = useState('');
  const [reports, setReports] = useState<Report[]>();

  useEffect(() => {
    if (show) {
      setReason('');
      setReports(undefined);
      if (request) {
        fetch((c) => c.guilds.member(request.guild_id, request.creator.id).reports())
          .then((res) => setReports(res.data))
          .catch();
      }
    }
  }, [show]);

  const type = request?.isAccept ? 'accept' : 'decline';

  const _process = () => {
    if (!request) return;
    fetch((c) =>
      c.guilds.respondUnbanrequest(request.guild_id, request.id, {
        status: request.isAccept ? UnbanRequestState.ACCEPTED : UnbanRequestState.DECLINED,
        processed_message: reason,
      } as UnbanRequest),
    )
      .then((res) => {
        pushNotification({
          message: t(`notifications.${type}`),
          type: NotificationType.SUCCESS,
        });
        onProcessed(res);
        onClose();
      })
      .catch();
  };

  return (
    <Modal
      show={show}
      onClose={onClose}
      heading={t(`title.${type}`)}
      controls={
        <>
          <Button disabled={!setReason} onClick={_process}>
            {t(`controls.${type}`)}
          </Button>
          <Button variant="gray" onClick={onClose}>
            {t('controls.cancel')}
          </Button>
        </>
      }>
      <StyledModalContainer>
        <section>
          <Trans ns="components" i18nKey={`modalprocessunbanrequest.description.${type}`}>
            {{ username: `${request?.creator.username}#${request?.creator.discriminator}` }}
            {{ userid: request?.creator.id }}
          </Trans>
        </section>
        <section>
          <Heading>{t('reason')}</Heading>
          <ModalTextArea value={reason} onInput={(e) => setReason(e.currentTarget.value)} />
        </section>
        <section>
          <Heading>{t('reports')}</Heading>
          <ReportsContainer>
            <div>
              {(reports && reports.map((r) => <ReportTile key={r.id} report={r} />)) || (
                <Loader height="10em" />
              )}
            </div>
          </ReportsContainer>
        </section>
      </StyledModalContainer>
    </Modal>
  );
};
