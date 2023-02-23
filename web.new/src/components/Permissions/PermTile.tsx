import { LinearGradient } from '../styleParts';
import { Tag } from '../Tag';
import { isAllowed } from './util';
import styled from 'styled-components';

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
