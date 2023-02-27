import { PropsWithChildren } from 'react';
import styled from 'styled-components';

type Props = PropsWithChildren & {
  icon: JSX.Element;
};

const NoEntries = styled.div`
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2em;
  margin-top: 2em;
  font-weight: 300;
  opacity: 0.5;
  text-align: center;
  line-height: 1.5em;

  > span {
    max-width: 30em;
  }

  > svg {
    width: 10em;
    height: 10em;
    stroke-width: 0.5px;
    opacity: 0.5;
  }
`;

export const EmptyPlaceholder: React.FC<Props> = ({ icon, children }) => {
  return (
    <NoEntries>
      {icon}
      <span>{children}</span>
    </NoEntries>
  );
};
