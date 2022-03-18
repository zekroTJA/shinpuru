import styled from 'styled-components';
import { Member } from '../../lib/shinpuru-ts/src';
import { memberName } from '../../util/users';
import { Container } from '../Container';
import { DiscordImage } from '../DiscordImage';
import { Clickable } from '../styleParts';

interface Props {
  member: Member;
}

const StyledContainer = styled(Container)`
  ${Clickable()}

  display: flex;
  padding: 0.5em;

  > img {
    width: 3em;
    height: 3em;
  }
`;

const Details = styled.div`
  margin-left: 0.5em;

  > h4 {
    margin: 0 0 0.5em 0;
    font-weight: 600;
  }

  > span {
    font-size: 0.8rem;
  }
`;

export const MemberTile: React.FC<Props> = ({ member }) => {
  return (
    <StyledContainer>
      <DiscordImage src={member.avatar_url} />
      <Details>
        <h4>{memberName(member)}</h4>
        <span>
          {member.user.username}#{member.user.discriminator}
        </span>
      </Details>
    </StyledContainer>
  );
};
