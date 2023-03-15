import styled from 'styled-components';

export const LinkButton = styled.button`
  background: transparent;
  color: ${(p) => p.theme.accent};
  border: none;
  outline: none;
  text-decoration: underline;
  cursor: pointer;
  font-size: 1em;
  padding: 0;
`;
