import { Trans, useTranslation } from 'react-i18next';
import { formatSince } from '../../util/date';

type Props = {
  date: string | Date | undefined | null;
};

export const SinceDate: React.FC<Props> = ({ date }) => {
  const { i18n } = useTranslation();
  const dateSince = formatSince(date, i18n.language);
  return (
    <Trans i18nKey="sincedate.since" ns="components">
      <span>{{ dateSince }}</span>
    </Trans>
  );
};
