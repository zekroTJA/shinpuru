import { Link } from 'react-router-dom';
import styled from 'styled-components';
import { useMatch } from 'react-router';

type Props = React.HTMLAttributes<HTMLDivElement> & {
  path: string;
};

const StyledLink = styled(Link)`
  color: ${(p) => p.theme.text};
  text-decoration: none;
`;

const StyledDiv = styled.div<{ activated: boolean }>`
  display: flex;
  gap: 0.5em;
  align-items: center;
  padding: 0.5em;
  background-color: ${(p) => (p.activated ? p.theme.accentDarker : p.theme.background)};
  border-radius: 8px;
  margin-top: 0.5em;
  cursor: pointer;
  transition: background-color 0.2s ease;

  > svg {
    stroke-width: 2;
    height: 1.2em;
    width: 1.2em;
  }
`;

export const Entry: React.FC<Props> = ({ path, children, ...props }) => {
  const match = useMatch(path);
  return (
    <StyledLink to={path}>
      <StyledDiv activated={!!match} {...props}>
        {children}
      </StyledDiv>
    </StyledLink>
  );
};
