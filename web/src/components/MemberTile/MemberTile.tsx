import styled from 'styled-components';
import { Member } from '../../lib/shinpuru-ts/src';
import { Container } from '../Container';
import { DiscordImage } from '../DiscordImage';

interface Props {
  member: Member;
}

const StyledContainer = styled(Container)`
  margin-top: 2em;
  padding: 0.5em;

  > img {
    width: 3em;
  }
`;

export const MemberTile: React.FC<Props> = ({ member }) => {
  return (
    <StyledContainer>
      <DiscordImage src={member.avatar_url} />
    </StyledContainer>
  );
};
