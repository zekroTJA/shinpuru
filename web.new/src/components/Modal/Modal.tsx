import styled from 'styled-components';
import { Container } from '../Container';
import { Heading } from '../Heading';
import { ModalControls } from './ModalControls';
import { ModalHeader } from './ModalHeader';

const BACKGROUND_ID = ':modal-background';

export type ControlProps = {
  show?: boolean;
  onClose?: () => void;
};

type Content = {
  heading?: string | JSX.Element | JSX.Element[];
  controls?: JSX.Element | JSX.Element[];
};

type Props = ControlProps & Content & React.HTMLAttributes<HTMLDivElement>;

const ModalOutlet = styled.div<ControlProps>`
  z-index: 100;
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0 0 0 / ${(p) => (p.show ? '75%' : '0')});
  transition: all 0.25s ease;
  pointer-events: ${(p) => (p.show ? 'all' : 'none')};
  display: flex;
  align-items: center;
  justify-content: center;
`;

const ModalContainer = styled(Container)<ControlProps>`
  position: relative;
  opacity: ${(p) => (p.show ? '1' : '0')};
  transform: scale(${(p) => (p.show ? '1' : '0.9')});
  transition: all 0.25s ease;
  padding: 0;
  max-height: 90vh;
  width: 90vw;
  max-width: fit-content;
  overflow-y: auto;

  > section {
    padding: 1em;
  }
`;

export const Modal: React.FC<Props> = ({
  children,
  show,
  controls,
  heading,
  onClose = () => {},
  ...props
}) => {
  const _heading = typeof heading === 'string' ? <Heading>{heading}</Heading> : heading;

  const _onBackgroundClick: React.MouseEventHandler<HTMLDivElement> = (e) => {
    if ((e.target as HTMLElement).id !== BACKGROUND_ID) return;
    onClose();
  };

  return (
    <ModalOutlet id={BACKGROUND_ID} show={show} onMouseDown={_onBackgroundClick} {...props}>
      <ModalContainer show={show}>
        {_heading && <ModalHeader>{_heading}</ModalHeader>}
        <section>{children}</section>
        {controls && <ModalControls>{controls}</ModalControls>}
      </ModalContainer>
    </ModalOutlet>
  );
};
