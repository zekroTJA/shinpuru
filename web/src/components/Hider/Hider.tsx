import { useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import styled from 'styled-components';

type Props = React.HTMLAttributes<HTMLDivElement> & {
  content: string;
};

const StyledDiv = styled.div<{ ln: number }>`
  position: relative;

  > input {
    font-size: 1rem;
    font-family: 'Roboto', sans-serif;
    color: currentColor;
    width: calc(${(p) => p.ln}ch + 2em);
    text-align: center;
    background-color: ${(p) => p.theme.white};
    border: none;
    padding: 0.8em 1em;
    border-radius: 3px;
  }
`;

const Hint = styled.p`
  @keyframes hint-anim {
    0% {
      opacity: 0;
      transform: translateY(-2em);
    }
    50% {
      opacity: 1;
    }
    100% {
      opacity: 0;
      transform: translateY(-3.5em);
    }
  }

  position: absolute;
  top: 0;
  left: 0;
  font-size: 0.6em;
  color: ${(p) => p.theme.white};
  background-color: ${(p) => p.theme.green};
  padding: 0.2em;
  border-radius: 3px;
  width: 100%;
  opacity: 0;
  transform: translateY(-1.5em);
  animation: hint-anim 1.5s ease;
`;

export const Hider: React.FC<Props> = ({ content, ...props }) => {
  const { t } = useTranslation('components');
  const [isHover, setIsHover] = useState(false);
  const [showCopyHint, setShowCopyHint] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);
  const timeoutRef = useRef<ReturnType<typeof setTimeout>>();

  const _mouseEnter = () => {
    setIsHover(true);
    inputRef.current?.focus();
    inputRef.current?.setSelectionRange(0, content.length);
    navigator.clipboard.writeText(content).then(() => {
      clearTimeout(timeoutRef.current);
      setShowCopyHint(true);
      timeoutRef.current = setTimeout(() => setShowCopyHint(false), 1500);
    });
  };

  const _mouseLeave = () => {
    setIsHover(false);
  };

  const htsTxt = t('hider.placeholder');

  return (
    <StyledDiv
      onMouseEnter={_mouseEnter}
      onMouseLeave={_mouseLeave}
      ln={content.length < htsTxt.length ? htsTxt.length : content.length}
      {...props}
    >
      <input ref={inputRef} value={isHover ? content : htsTxt} readOnly />
      {showCopyHint && <Hint>Copied to clipboard!</Hint>}
    </StyledDiv>
  );
};
