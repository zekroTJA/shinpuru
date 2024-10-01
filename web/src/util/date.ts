import {de, enUS, Locale} from 'date-fns/locale';
import { format, formatDistance } from 'date-fns';

const LANG_MAP: { [key: string]: Locale } = {
  'en-US': enUS,
  de: de,
  'de-DE': de,
};

export const formatDate = (date: string | Date | undefined | null, locale?: string) => {
  if (!date) return 'n/a';
  const _date = typeof date === 'string' ? new Date(date) : date;
  return format(_date, 'dd/LL/yyyy HH:mm:ss O');
};

export const formatSince = (date: string | Date | undefined | null, locale?: string) => {
  if (!date) return 'n/a';
  const _date = typeof date === 'string' ? new Date(date) : date;
  return formatDistance(_date, new Date(), { locale: getLocale(locale) });
};

const getLocale = (v?: string): Locale => (!!v ? LANG_MAP[v] : undefined) ?? enUS;

export const parseToDateString = (v: Date | number) => format(v, "yyyy-MM-dd'T'HH:mm:ssxxx");
