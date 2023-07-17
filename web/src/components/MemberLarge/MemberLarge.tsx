import { Guild, Member } from '../../lib/shinpuru-ts/src';

import { Clickable } from '../styleParts';
import { Container } from '../Container';
import { DiscordImage } from '../DiscordImage';
import { Embed } from '../Embed';
import { Flex } from '../Flex';
import { RoleList } from '../RoleList';
import { memberName } from '../../util/users';
import styled from 'styled-components';

type Props = {
  member?: Member;
  guild?: Guild;
  onClick?: (member: Member) => void;
};

const StyledContainer = styled(Container)`
  ${Clickable(1.01)}

  display: flex;
  width: 100%;

  > img {
    height: 4em;
    margin-right: 1em;
  }

  > div {
    > * {
      margin: 0 0 0.6em 0;
      &:last-child {
        margin: 0;
      }
    }
  }

  p,
  h2 {
    margin: 0;
  }
`;

const Header = styled(Flex)`
  @media screen and (max-width: 800px) {
    > h2 {
      display: block;
      width: 100%;
    }
  }

  flex-wrap: wrap;
  gap: 1em;
  align-items: center;
`;

export const MemberLarge: React.FC<Props> = ({ member, guild, onClick = () => {} }) => {
  return !!member && !!guild ? (
    <StyledContainer onClick={() => onClick(member)}>
      <DiscordImage src={member.avatar_url} />
      <div>
        <Header>
          <h2>{memberName(member)}</h2>
          <small>{member.user.username}</small>
          <Embed>{member.user.id}</Embed>
        </Header>
        <RoleList guildroles={guild.roles!} roleids={member.roles} />
      </div>
    </StyledContainer>
  ) : (
    <>loading...</>
  );
};
