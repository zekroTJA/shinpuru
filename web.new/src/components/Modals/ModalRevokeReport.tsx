import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import styled from 'styled-components';
import { Report } from '../../lib/shinpuru-ts/src';
import { Button } from '../Button';
import { ControlProps, Modal } from '../Modal';
import { TextArea } from '../TextArea';

type Props = ControlProps & {
  report?: Report;
  onConfirm: (report: Report, reason: string) => void;
};

const StyledTextArea = styled(TextArea)`
  background-color: ${(p) => p.theme.background3};
`;

export const ModalRevokeReport: React.FC<Props> = ({ report, onClose = () => {}, onConfirm }) => {
  const { t } = useTranslation('components', { keyPrefix: 'modalrevokereport.' });
  const [reason, setReason] = useState('');

  const _onClose = () => {
    onClose();
    setReason('');
  };

  const _onConfirm = () => {
    if (!report) return;
    onConfirm(report, reason);
    _onClose();
  };

  return (
    <Modal
      show={!!report}
      onClose={_onClose}
      heading={t('heading')}
      controls={
        <>
          <Button disabled={!reason} onClick={_onConfirm}>
            {t('controls.revoke')}
          </Button>
          <Button variant="gray" onClick={_onClose}>
            {t('controls.cancel')}
          </Button>
        </>
      }>
      <p>{t(`description.${report?.type_name?.toLowerCase()}`)}</p>
      <StyledTextArea value={reason} onInput={(e) => setReason(e.currentTarget.value)} />
    </Modal>
  );
};
