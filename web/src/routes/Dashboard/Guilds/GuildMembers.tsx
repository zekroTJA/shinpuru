import { useSelfUser } from '../../../hooks/useSelfUser';

interface Props {}

export const GuildMembersRoute: React.FC<Props> = ({}) => {
  const selfUser = useSelfUser();

  return <>members</>;
};
