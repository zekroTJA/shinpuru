import { useSelfUser } from '../hooks/useSelfUser';

interface Props {}

export const HomeRoute: React.FC<Props> = ({}) => {
  const selfUser = useSelfUser();
  return <>Hello {selfUser?.username}</>;
};
