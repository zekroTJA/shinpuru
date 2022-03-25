import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import ReactMarkdown from 'react-markdown';
import { useLocalStorage } from '../../hooks/useLocalStorage';
import { Button } from '../Button';
import { Modal } from '../Modal/Modal';

interface Props {}

export const ModalBetaGreeter: React.FC<Props> = () => {
  const { t } = useTranslation('components');
  const [show, setShow] = useState(false);
  const [dismissed, setDismissed] = useLocalStorage('shnp.betagreeter.dismissed');

  useEffect(() => {
    if (dismissed) return;
    setTimeout(() => setShow(true), 1000);
  }, []);

  const _dismiss = () => {
    setDismissed(true);
    setShow(false);
  };

  const _backToStable = () => {
    window.location.assign('/guilds');
  };

  return (
    <Modal
      show={show}
      heading={t('modalbetagreeter.heading')}
      controls={
        <>
          <Button variant="gray" onClick={_backToStable}>
            {t('modalbetagreeter.controls.back')}
          </Button>
          <Button onClick={_dismiss}>{t('modalbetagreeter.controls.accept')}</Button>
        </>
      }>
      <ReactMarkdown children={t('modalbetagreeter.message')} />
    </Modal>
  );
};
