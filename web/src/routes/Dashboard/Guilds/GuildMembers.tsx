import { useParams } from 'react-router';
import styled from 'styled-components';
import { MemberLarge } from '../../../components/MemberLarge';
import { MemberTile } from '../../../components/MemberTile';
import { useGuild } from '../../../hooks/useGuild';
import { useMembers } from '../../../hooks/useMembers';
import { useSelfMember } from '../../../hooks/useSelfMember';

interface Props {}

const MembersContainer = styled.div`
  display: flex;
  flex-wrap: wrap;
`;

export const GuildMembersRoute: React.FC<Props> = ({}) => {
  const { guildid } = useParams();
  const selfMember = useSelfMember(guildid);
  const guild = useGuild(guildid);
  const [members, loadMoreMembers] = useMembers(guildid);

  return (
    <>
      <MemberLarge member={selfMember} guild={guild} />
      {members && (
        <MembersContainer>
          {members.map((m) => (
            <MemberTile member={m} />
          ))}
        </MembersContainer>
      )}
    </>
  );
};
