import i18next from 'i18next';
import I18nextBrowserLanguageDetector from 'i18next-browser-languagedetector';
import I18NextHttpBackend from 'i18next-http-backend';
import { initReactI18next } from 'react-i18next';

i18next
  .use(initReactI18next)
  .use(I18NextHttpBackend)
  .use(I18nextBrowserLanguageDetector)
  .init({
    fallbackLng: 'en-US',
    interpolation: { escapeValue: true },
    backend: { loadPath: import.meta.env.BASE_URL + 'locales/{{lng}}/{{ns}}.json' },
  });

export default i18next;
