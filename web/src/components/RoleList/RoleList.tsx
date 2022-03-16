import styled from 'styled-components';
import { Role } from '../../lib/shinpuru-ts/src';
import { Flex } from '../Flex';
import { Tag } from '../Tag';

interface Props {
  roleids: string[];
  guildroles: Role[];
}

const RolesContainer = styled(Flex)`
  flex-wrap: wrap;
  gap: 0.4em;
`;

export const RoleList: React.FC<Props> = ({ roleids, guildroles }) => {
  const roles = roleids
    .map((rid) => guildroles.find((r) => r.id === rid))
    .filter((r) => !!r)
    .sort((ra, rb) => rb!.position - ra!.position)
    .map((r) => (
      <Tag key={r!.id} colors={r!.color} borderRadius="8px">
        {r!.name}
      </Tag>
    ));
  return <RolesContainer>{roles}</RolesContainer>;
};
