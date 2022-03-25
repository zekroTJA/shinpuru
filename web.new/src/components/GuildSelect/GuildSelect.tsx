import { Guild } from '../../lib/shinpuru-ts/src';
import { Select, Element } from '../Select';
import { Option } from './Option';

type Props = React.HTMLAttributes<HTMLDivElement> & {
  guilds: Guild[];
  value?: Guild;
  onElementSelect?: (g: Guild) => void;
};

export const GuildSelect: React.FC<Props> = ({
  guilds,
  value,
  onElementSelect = () => {},
  ...props
}) => {
  const options = guilds.map(
    (g) =>
      ({
        id: g.id,
        display: <Option guild={g} />,
        value: g,
      } as Element<Guild>),
  );
  return (
    <Select
      options={options}
      value={options.find((o) => o.id === value?.id)}
      onElementSelect={(e) => onElementSelect(e.value)}
      {...props}
    />
  );
};
