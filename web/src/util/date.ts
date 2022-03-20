import { format, formatDistance } from 'date-fns';
import { enUS, de } from 'date-fns/locale';

const LANG_MAP: { [key: string]: Locale } = {
  'en-US': enUS,
  'de': de,
  'de-DE': de,
};

export const formatDate = (
  date: string | Date | undefined | null,
  locale?: string
) => {
  if (!date) return 'n/a';
  if (typeof date === 'string') date = new Date(date);
  return format(date, 'dd/LL/yyyy HH:mm:ss O');
};

export const formatSince = (
  date: string | Date | undefined | null,
  locale?: string
) => {
  if (!date) return 'n/a';
  if (typeof date === 'string') date = new Date(date);
  return formatDistance(date, new Date(), { locale: getLocale(locale) });
};

const getLocale = (v?: string): Locale =>
  (!!v ? LANG_MAP[v] : undefined) ?? enUS;
