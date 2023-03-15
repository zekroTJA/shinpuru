import { ButtonVariant } from '../components/Button';
import { useStore } from '../services/store';

export type DisplayElement = string | JSX.Element | JSX.Element[];

export type Control<T extends unknown> = {
  name: string;
  value: T;
  variant?: ButtonVariant;
};

export type OpenModal<T> = {
  content: DisplayElement;
  heading?: DisplayElement;
  controls?: Control<T>[];
};

export type ModalState<T> = {
  modal?: OpenModal<T>;
  isOpen: boolean;
  resolver?: (value: T | null) => void;
};

export const useModal = <TResult>() => {
  const [modal, setModal] = useStore((s) => [s.modal, s.setModal]);

  const openModal = (m: OpenModal<TResult>) => {
    if (modal?.isOpen) return Promise.reject('another modal is already open');
    return new Promise<TResult | null>((resolver) => {
      setModal({
        modal: m,
        isOpen: true,
        resolver,
      });
    });
  };

  return { openModal };
};
