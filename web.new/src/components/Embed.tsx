import styled from 'styled-components';

export const Embed = styled.span`
  font-family: 'Roboto Mono', monospace;
  font-size: 0.7em;
  font-weight: bolder;
  padding: 0.2em 0.4em;
  border-radius: 3px;
  background-color: rgba(0 0 0 / 10%);
  width: fit-content;
`;

export const EmbedWrapper: React.FC<{ value: string | number | JSX.Element | undefined }> = ({
  value,
}) => <Embed>{value}</Embed>;
