import styled from 'styled-components';
import { Guild } from '../../lib/shinpuru-ts/src';
import { DiscordImage } from '../DiscordImage';

interface Props {
  guild: Guild;
}

const StyledDiv = styled.div`
  display: flex;
  align-items: center;
  gap: 0.5em;

  > img,
  svg {
    height: 1.2em;
    aspect-ratio: 1;
  }
`;

export const Option: React.FC<Props> = ({ guild }) => {
  return (
    <StyledDiv>
      <DiscordImage src={guild.icon_url} round />
      <span>{guild.name}</span>
    </StyledDiv>
  );
};
