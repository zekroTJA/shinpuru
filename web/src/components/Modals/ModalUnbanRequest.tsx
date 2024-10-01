import { ModalContainer, ModalTextArea } from './modalParts';
import { Trans, useTranslation } from 'react-i18next';
import { useEffect, useState } from 'react';

import { Button } from '../Button';
import { ControlProps } from '../Modal/Modal';
import { Guild } from '../../lib/shinpuru-ts/src';
import { Modal } from '../Modal';

type Props = ControlProps & {
  guild?: Guild;
  onSend?: (message: string) => void;
};

export const ModalUnbanRequest: React.FC<Props> = ({
  guild,
  show,
  onSend = () => {},
  onClose = () => {},
  ...props
}) => {
  const { t } = useTranslation('components', { keyPrefix: 'modalunbanrequest' });

  const [message, setMessage] = useState('');

  const _onSend = () => {
    onSend(message);
    onClose();
  };

  useEffect(() => {
    if (show) {
      setMessage('');
    }
  }, [show]);

  return (
    <Modal
      show={show}
      onClose={onClose}
      heading={t('heading')}
      controls={
        <>
          <Button disabled={!message} onClick={_onSend}>
            {t('controls.send')}
          </Button>
          <Button variant="gray" onClick={onClose}>
            {t('controls.cancel')}
          </Button>
        </>
      }>
      <ModalContainer>
        <span>
          <Trans
            ns="components"
            i18nKey="modalunbanrequest.sub"
            values={{ guildname: guild?.name }}
            components={{ '1': <strong /> }}
          />
        </span>
        <span>{t('description')}</span>
        <ModalTextArea value={message} onInput={(e) => setMessage(e.currentTarget.value)} />
      </ModalContainer>
    </Modal>
  );
};
