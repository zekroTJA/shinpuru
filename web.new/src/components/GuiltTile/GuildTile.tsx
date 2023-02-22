import { Clickable } from '../styleParts';
import { Container } from '../Container';
import { DiscordImage } from '../DiscordImage';
import { Embed } from '../Embed';
import { Guild } from '../../lib/shinpuru-ts/src';
import styled from 'styled-components';

type Props = React.HTMLAttributes<HTMLDivElement> & {
  guild: Guild;
};

const GuildContainer = styled(Container)`
  display: flex;
  padding: 0;
  cursor: pointer;
  width: fit-content;

  ${Clickable()}

  > img {
    width: 6em;
  }

  > div {
    padding: 1em;
    display: flex;
    flex-direction: column;
    justify-content: center;
    gap: 0.4em;
  }
`;

export const GuildTile: React.FC<Props> = ({ guild, ...props }) => {
  return (
    <GuildContainer {...props}>
      <DiscordImage src={guild.icon_url} />
      <div>
        <strong>{guild.name}</strong>
        <Embed>{guild.id}</Embed>
      </div>
    </GuildContainer>
  );
};
