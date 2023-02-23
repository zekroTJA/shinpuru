import styled, { css } from 'styled-components';
import { useEffect, useRef, useState } from 'react';

import { ReactComponent as AddFileIcon } from '../../assets/addfile.svg';
import { ReactComponent as ErrorFileIcon } from '../../assets/errorfile.svg';
import { ReactComponent as FileIcon } from '../../assets/file.svg';
import { byteFormatter } from 'byte-formatter';
import { useTranslation } from 'react-i18next';

type Props = {
  file?: File;
  onFileInput?: (file: File) => void;
};

const Container = styled.div`
  width: 100%;
  height: 4em;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  cursor: pointer;
`;

const FiledropContainer = styled(Container)<{ isError: boolean; isDragging: boolean }>`
  justify-content: flex-start;

  > svg {
    height: 100%;
    width: auto;
    padding: 0.8em;
    border-radius: 8px;
  }

  > div > * {
    display: block;
    margin: 0.3em 0;
  }

  border: dashed 2px currentColor;
  transition: all 0.25s ease;

  ${(p) =>
    p.isError
      ? css`
          color: ${p.theme.red};
        `
      : ''}

  ${(p) =>
    p.isDragging
      ? css`
          color: ${p.theme.green};
        `
      : ''}
`;

const FileContainer = styled(Container)`
  justify-content: flex-start;

  > svg {
    height: 100%;
    width: auto;
    margin-right: 1em;
    padding: 0.8em;
    background-color: ${(p) => p.theme.background3};
    border-radius: 8px;
  }

  > div > * {
    display: block;
    margin: 0.3em 0;
  }
`;

const Smol = styled.small``;

export const Filedrop: React.FC<Props> = ({ file, onFileInput = () => true }) => {
  const { t } = useTranslation('components', { keyPrefix: 'filedrop' });
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [error, setError] = useState<string>();
  const [dragging, setDragging] = useState(false);

  const _fileInputChange: React.ChangeEventHandler<HTMLInputElement> = (e) =>
    _setFile((e.currentTarget.files ?? [])[0]);

  const _onDragOver: React.DragEventHandler<HTMLDivElement> = (e) => {
    e.stopPropagation();
    e.preventDefault();
    e.dataTransfer.dropEffect = 'copy';
    setDragging(true);
    setError(undefined);
  };

  const _onDragEnd: React.DragEventHandler<HTMLDivElement> = (e) => {
    setDragging(false);
  };

  const _onDrop: React.DragEventHandler<HTMLDivElement> = (e) => {
    e.preventDefault();
    setDragging(false);
    _setFile((e.dataTransfer.files ?? [])[0]);
  };

  const _onPaste = (e: ClipboardEvent) => {
    const item = (e.clipboardData?.items ?? [])[0];
    if (!item || item.kind !== 'file') return;
    const file = item.getAsFile();
    if (file) _setFile(file);
  };

  const _setFile = (file?: File) => {
    if (!!file) {
      try {
        onFileInput(file);
      } catch (e) {
        setError(e instanceof Error ? e.message : (e as string));
      }
    }
  };

  useEffect(() => {
    setError(undefined);
    setDragging(false);
  }, [file]);

  useEffect(() => {
    document.addEventListener('paste', _onPaste);
    return () => document.removeEventListener('paste', _onPaste);
  }, []);

  return (
    <div
      onClick={() => fileInputRef.current?.click()}
      onDrop={_onDrop}
      onDragOver={_onDragOver}
      onDragExit={_onDragEnd}>
      {(file && !error && (
        <FileContainer>
          <FileIcon />
          <div>
            <strong>{file.name}</strong>
            <span>{byteFormatter(file.size)}</span>
          </div>
        </FileContainer>
      )) || (
        <FiledropContainer isError={!!error} isDragging={dragging}>
          {(error && <ErrorFileIcon />) || <AddFileIcon />}
          {(error && <span>{error}</span>) || <Smol>{t('info')}</Smol>}
        </FiledropContainer>
      )}
      <input ref={fileInputRef} type="file" hidden onChange={_fileInputChange} />
    </div>
  );
};
