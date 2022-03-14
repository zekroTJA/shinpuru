import styled from 'styled-components';
import { Guild } from '../../lib/shinpuru-ts/src';
import { ReactComponent as DiscordIcon } from '../../assets/dc-logo.svg';

interface Props {
  guild: Guild;
}

const StyledDiv = styled.div`
  display: flex;
  align-items: center;

  > img,
  svg {
    height: 1.2em;
    aspect-ratio: 1;
    margin-right: 0.5em;
    border-radius: 100%;
  }
`;

export const Option: React.FC<Props> = ({ guild }) => {
  const guildIcon =
    guild.icon_url === '' ? <DiscordIcon /> : <img src={guild.icon_url} />;
  return (
    <StyledDiv>
      {guildIcon}
      <span>{guild.name}</span>
    </StyledDiv>
  );
};
