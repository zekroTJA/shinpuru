import { AppTheme, getSystemTheme } from '../../theme/theme';
import { Element, Select } from '../../components/Select';
import React, { useEffect, useState } from 'react';

import { ReactComponent as GearIcon } from '../../assets/settings.svg';
import { MaxWidthContainer } from '../../components/MaxWidthContainer';
import { ReactComponent as MoonIcon } from '../../assets/halfmoon.svg';
import { Small } from '../../components/Small';
import { ReactComponent as SunIcon } from '../../assets/sun.svg';
import styled from 'styled-components';
import { useStore } from '../../services/store';
import { useTranslation } from 'react-i18next';

type Props = {};

const IconOption = styled.div`
  display: flex;
  gap: 0.5em;
  align-items: center;
`;

const ThemeRoute: React.FC<Props> = ({}) => {
  const { t, i18n } = useTranslation('routes.usersettings.theme');
  const [lang, setLang] = useState<string>();
  const [scheme, setScheme] = useStore((s) => [s.theme, s.setTheme]);

  useEffect(() => {
    setLang(i18n.language);
  }, []);

  useEffect(() => {
    i18n.changeLanguage(lang);
  }, [lang]);

  const schemes: Element<AppTheme>[] = [
    {
      id: 'dark',
      display: (
        <IconOption>
          <MoonIcon />
          {t('scheme.dark')}
        </IconOption>
      ),
      value: AppTheme.DARK,
    },
    {
      id: 'light',
      display: (
        <IconOption>
          <SunIcon />
          {t('scheme.light')}
        </IconOption>
      ),
      value: AppTheme.LIGHT,
    },
    {
      id: 'system',
      display: (
        <IconOption>
          <GearIcon />
          {t('scheme.system')}
        </IconOption>
      ),
      value: getSystemTheme(),
    },
  ];

  return (
    <MaxWidthContainer>
      <h1>{t('heading')}</h1>
      <Small>{t('explanation')}</Small>
      <h2>{t('scheme.heading')}</h2>
      <Select
        options={schemes}
        value={schemes.find((s) => s.value === scheme)}
        onElementSelect={(s) => setScheme(s.value)}
      />
    </MaxWidthContainer>
  );
};

export default ThemeRoute;
