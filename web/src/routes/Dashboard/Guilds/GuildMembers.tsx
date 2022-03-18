import { useCallback, useState } from 'react';
import { useParams } from 'react-router';
import styled from 'styled-components';
import { Button } from '../../../components/Button';
import { MemberLarge } from '../../../components/MemberLarge';
import { MemberTile } from '../../../components/MemberTile';
import { SearchBar } from '../../../components/SearchBar';
import { useGuild } from '../../../hooks/useGuild';
import { useMembers } from '../../../hooks/useMembers';
import { useSelfMember } from '../../../hooks/useSelfMember';
import { debounce } from 'debounce';
import { useTranslation } from 'react-i18next';

interface Props {}

const MembersSection = styled.div`
  margin-top: 1em;
`;

const MemberTiles = styled.div`
  margin-top: 1em;
  display: flex;
  flex-wrap: wrap;
  gap: 1em;
`;

const LoadMoreButton = styled(Button)`
  margin: 1em auto 0 auto;
`;

export const GuildMembersRoute: React.FC<Props> = ({}) => {
  const { t } = useTranslation('routes.guildmembers');
  const { guildid } = useParams();
  const selfMember = useSelfMember(guildid);
  const guild = useGuild(guildid);
  const [search, setSearch] = useState('');
  const [members, loadMoreMembers] = useMembers(guildid, 100, search);

  const _onSearchInput = useCallback(debounce(setSearch, 500), []);

  return (
    <>
      <MemberLarge member={selfMember} guild={guild} />
      {members && selfMember && (
        <MembersSection>
          <SearchBar
            onValueChange={_onSearchInput}
            placeholder={t('searchplaceholder')}
          />
          <MemberTiles>
            {members
              .filter((m) => m.user.id !== selfMember.user.id)
              .map((m) => (
                <MemberTile key={`memb-${m.user.id}`} member={m} />
              ))}
          </MemberTiles>
          {members.length > 0 &&
            !search &&
            guild?.member_count! > members.length && (
              <LoadMoreButton onClick={() => loadMoreMembers()}>
                Laod more ...
              </LoadMoreButton>
            )}
        </MembersSection>
      )}
    </>
  );
};
