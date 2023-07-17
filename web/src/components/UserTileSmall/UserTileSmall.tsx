import { Container } from '../Container';
import { DiscordImage } from '../DiscordImage';
import { FlatUser } from '../../lib/shinpuru-ts/src';
import styled from 'styled-components';

type Props = {
  fallbackId: string;
  user: FlatUser;
  icon?: JSX.Element;
  hideAvatar?: boolean;
};

const UserContainer = styled(Container)`
  display: flex;
  padding: 0.8em;
  align-items: center;
  gap: 1em;
  border-radius: 3px;
  width: 100%;

  > svg {
    width: 1.2em;
    height: auto;
  }

  > img {
    width: 2em;
    height: auto;
  }
`;

export const UserTileSmall: React.FC<Props> = ({ fallbackId, user, icon, hideAvatar }) => {
  return (
    <UserContainer>
      {icon}
      {hideAvatar || <DiscordImage src={user?.avatar_url} round />}
      <span>{user ? user.username : <i>{fallbackId}</i>}</span>
    </UserContainer>
  );
};
