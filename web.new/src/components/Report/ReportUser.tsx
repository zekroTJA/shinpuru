import styled from 'styled-components';
import { ReactComponent as HammerIcon } from '../../assets/hammer.svg';
import { ReactComponent as HammerTarget } from '../../assets/target.svg';
import { FlatUser } from '../../lib/shinpuru-ts/src';
import { Container } from '../Container';
import { DiscordImage } from '../DiscordImage';

type Props = {
  fallbackId: string;
  user: FlatUser;
  isEcecutor?: boolean;
};

const ReportUserContainer = styled(Container)`
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

export const ReportUser: React.FC<Props> = ({ fallbackId, user, isEcecutor }) => {
  return (
    <ReportUserContainer>
      {isEcecutor ? <HammerIcon /> : <HammerTarget />}
      <DiscordImage src={user?.avatar_url} round />
      <span>{user ? `${user.username}#${user.discriminator}` : <i>{fallbackId}</i>}</span>
    </ReportUserContainer>
  );
};
