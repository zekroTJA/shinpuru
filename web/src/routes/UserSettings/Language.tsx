import { Element, Select } from '../../components/Select';
import React, { useEffect, useState } from 'react';

import { MaxWidthContainer } from '../../components/MaxWidthContainer';
import { Small } from '../../components/Small';
import { useTranslation } from 'react-i18next';

type Props = {};

const LANGUAGE_OPTIONS: Element<string>[] = [
  {
    id: 'en-US',
    display: <span>🇺🇸&nbsp;&nbsp;English (US)</span>,
    value: 'en-US',
  },
  {
    id: 'de',
    display: <span>🇩🇪&nbsp;&nbsp;German</span>,
    value: 'de',
  },
];

const LanguageRoute: React.FC<Props> = () => {
  const { t, i18n } = useTranslation('routes.usersettings.language');
  const [lang, setLang] = useState<string>();

  useEffect(() => {
    setLang(i18n.language);
  }, []);

  useEffect(() => {
    i18n.changeLanguage(lang);
  }, [lang]);

  return (
    <MaxWidthContainer>
      <h1>{t('heading')}</h1>
      <Small>{t('explanation')}</Small>
      <Select
        options={LANGUAGE_OPTIONS}
        value={LANGUAGE_OPTIONS.find((l) => l.value === lang)}
        onElementSelect={(v) => setLang(v.value)}
      />
    </MaxWidthContainer>
  );
};

export default LanguageRoute;
