import { byteFormatter } from 'byte-formatter';
import { useEffect, useRef, useState } from 'react';
import styled, { css } from 'styled-components';
import { ReactComponent as AddFileIcon } from '../../assets/addfile.svg';
import { ReactComponent as FileIcon } from '../../assets/file.svg';

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

const FiledropContainer = styled(Container)<{ isError: boolean }>`
  border: dashed 2px currentColor;
  flex-direction: column;

  ${(p) =>
    p.isError
      ? css`
          color: ${p.theme.red};
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

export const Filedrop: React.FC<Props> = ({ file, onFileInput = () => true }) => {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [error, setError] = useState<string>();

  const _fileInputChange: React.ChangeEventHandler<HTMLInputElement> = (e) =>
    _setFile((e.currentTarget.files ?? [])[0]);

  const _onDrop: React.DragEventHandler<HTMLDivElement> = (e) => {
    e.preventDefault();
    _setFile((e.dataTransfer.files ?? [])[0]);
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
  }, [file]);

  return (
    <div onClick={() => fileInputRef.current?.click()} onDrop={_onDrop}>
      {(file && !error && (
        <FileContainer>
          <FileIcon />
          <div>
            <strong>{file.name}</strong>
            <span>{byteFormatter(file.size)}</span>
          </div>
        </FileContainer>
      )) || (
        <FiledropContainer isError={!!error}>
          <AddFileIcon />
          {error && <span>{error}</span>}
        </FiledropContainer>
      )}
      <input ref={fileInputRef} type="file" hidden onChange={_fileInputChange} />
    </div>
  );
};