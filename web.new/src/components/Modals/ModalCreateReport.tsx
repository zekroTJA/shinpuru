import { useEffect, useState } from 'react';
import { Trans, useTranslation } from 'react-i18next';
import styled from 'styled-components';
import { Member } from '../../lib/shinpuru-ts/src';
import { Button } from '../Button';
import { Filedrop } from '../Filedrop';
import { Heading } from '../Heading';
import { Modal } from '../Modal';
import { ControlProps } from '../Modal/Modal';
import { TextArea } from '../TextArea';

const ALLOWED_ATTACHMENT_TYPES = ['image/jpeg', 'image/jpg', 'image/png'];

export type ReportActionType = 'report' | 'kick' | 'ban' | 'mute';

type Props = ControlProps & {
  type: ReportActionType;
  member: Member;
};

const StyledTextArea = styled(TextArea)`
  background-color: ${(p) => p.theme.background3};
  min-width: 100%;
  max-width: 100%;
`;

const ModalContainer = styled.div`
  width: 30em;
  max-width: 80vw;
  display: flex;
  flex-direction: column;
  gap: 1.5em;
`;

const ACTION_TEXT = {
  kick: 'routes.member:moderation.kick',
  ban: 'routes.member:moderation.ban',
  mute: 'routes.member:moderation.mute',
  report: 'routes.member:moderation.report',
};

export const ModalCreateReport: React.FC<Props> = ({ show, onClose, type }) => {
  const { t } = useTranslation();
  const [reason, setReason] = useState('');
  const [attachment, setAttachment] = useState<File>();

  useEffect(() => {
    if (show) {
      setReason('');
      setAttachment(undefined);
    }
  }, [show]);

  const _setAttachment = (f: File) => {
    if (!ALLOWED_ATTACHMENT_TYPES.includes(f.type))
      throw t('errors.reports.disallowed-attachment-type');
    if (f.size > 50 * 1024 * 1024) throw t('errors.reports.file-too-big');
    setAttachment(f);
  };

  const action = t(ACTION_TEXT[type]);

  return (
    <Modal
      show={show}
      onClose={onClose}
      heading={action}
      controls={
        <>
          <Button disabled={!reason}>
            <Trans ns="components" i18nKey="modalcreatereport.controls.execute">
              {{ action }}
            </Trans>
          </Button>
          <Button variant="gray" onClick={onClose}>
            {t('components:modalcreatereport.controls.cancel')}
          </Button>
        </>
      }>
      <ModalContainer>
        <section>
          <Heading>{t('components:modalcreatereport.reason')}</Heading>
          <StyledTextArea value={reason} onInput={(e) => setReason(e.currentTarget.value)} />
        </section>
        <section>
          <Heading>{t('components:modalcreatereport.attachment')}</Heading>
          <Filedrop file={attachment} onFileInput={_setAttachment} />
        </section>
      </ModalContainer>
    </Modal>
  );
};
