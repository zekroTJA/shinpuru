import { useEffect, useState } from 'react';
import { useParams } from 'react-router';
import { MemberLarge } from '../../../components/MemberLarge';
import { useApi } from '../../../hooks/useApi';
import { useGuild } from '../../../hooks/useGuild';
import { useMembers } from '../../../hooks/useMembers';
import { useSelfMember } from '../../../hooks/useSelfMember';
import { Member } from '../../../lib/shinpuru-ts/src';

interface Props {}

export const GuildMembersRoute: React.FC<Props> = ({}) => {
  const { guildid } = useParams();
  const selfMember = useSelfMember(guildid);
  const guild = useGuild(guildid);
  const members = useMembers(guildid);

  return (
    <>
      <MemberLarge member={selfMember} guild={guild} />
    </>
  );
};
