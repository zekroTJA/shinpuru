import styled from 'styled-components';
import { Tag } from '../Tag';
import { LinearGradient } from '../styleParts';
import { isAllowed } from './util';

type Props = {
  perm: string;
};

export const StyledTag = styled(Tag)<{ allowed: boolean }>`
  ${(p) => LinearGradient(p.allowed ? p.theme.green : p.theme.red)};
  border: none;
`;

export const PermTile: React.FC<Props> = ({ perm }) => {
  return <StyledTag allowed={isAllowed(perm)}>{perm}</StyledTag>;
};
