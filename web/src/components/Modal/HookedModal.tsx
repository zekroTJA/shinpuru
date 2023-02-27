import { Button } from '../Button';
import { Modal } from './Modal';
import { uid } from 'react-uid';
import { useStore } from '../../services/store';

type Props = {};

export const HookedModal: React.FC<Props> = () => {
  const [modal, setModal] = useStore((s) => [s.modal, s.setModal]);

  const _resolve = (v: any | null) => {
    if (!modal.resolver) return;
    setModal({ ...modal, isOpen: false });
    modal.resolver(v);
  };

  const _controls = modal.modal?.controls?.map((c) => (
    <Button variant={c.variant} onClick={() => _resolve(c.value)} key={uid(c)}>
      {c.name}
    </Button>
  ));

  const _onClose = () => _resolve(null);

  return (
    <>
      <Modal
        show={modal.isOpen}
        heading={modal.modal?.heading}
        controls={_controls}
        onClose={_onClose}>
        {modal.modal?.content}
      </Modal>
    </>
  );
};
