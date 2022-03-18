import { useMatch } from 'react-router';
import { Link } from 'react-router-dom';
import styled from 'styled-components';

interface Props {
  path: string;
}

const StyledDiv = styled.div<{ activated: boolean }>`
  display: flex;
  padding: 0.5em;
  background-color: ${(p) =>
    p.activated ? p.theme.accentDarker : p.theme.background};
  border-radius: 8px;
  margin-top: 0.5em;

  > a {
    color: ${(p) => p.theme.text};
    text-decoration: none;

    > svg {
      margin-right: 0.5em;
    }
  }
`;

export const Entry: React.FC<Props> = ({ path, children }) => {
  const match = useMatch('db/' + path);
  return (
    <StyledDiv activated={!!match}>
      <Link to={path}>{children}</Link>
    </StyledDiv>
  );
};
