import { useParams } from 'react-router';

interface Props {}

export const Guild: React.FC<Props> = ({}) => {
  const { guildid } = useParams();
  return <>{guildid}</>;
};
