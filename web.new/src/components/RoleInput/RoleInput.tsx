import { Guild, Role } from '../../lib/shinpuru-ts/src';
import { TagElement, TagsInput } from '../TagsInput';

type Props = {
  guild: Guild;
  selected?: Role[];
  onChange: (v: Role[]) => void;
  placeholder?: string;
};

export const RoleInput: React.FC<Props> = ({ guild, selected, onChange, placeholder }) => {
  const roleTagOptions =
    guild?.roles?.map(
      (r) =>
        ({ id: r.id, display: r.name, keywords: [r.id, r.name], value: r } as TagElement<Role>),
    ) ?? [];

  const mapSelectedRole = (r: Role): TagElement<Role> => {
    const rto = roleTagOptions.find((e) => e.id === r.id)!;
    if (!rto) {
      return { id: r.id, display: r.name, keywords: [], value: r };
    }
    return rto;
  };

  return (
    <TagsInput
      options={roleTagOptions}
      selected={(selected ?? []).map(mapSelectedRole)}
      onChange={(v) => onChange(v.map((e) => e.value))}
      placeholder={placeholder}
    />
  );
};
