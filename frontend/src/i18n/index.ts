import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

import ruCommon from './ru/common.json';
import ruAuth from './ru/auth.json';
import tmCommon from './tm/common.json';
import tmAuth from './tm/auth.json';

const resources = {
  ru: {
    common: ruCommon,
    auth: ruAuth,
  },
  tm: {
    common: tmCommon,
    auth: tmAuth,
  },
};

const getLanguage = () => {
  const saved = localStorage.getItem('lang');
  if (saved) return saved;

  const browser = navigator.language.split('-')[0];
  return ['ru', 'tm'].includes(browser) ? browser : 'ru';
};

i18n.use(initReactI18next).init({
  resources,
  lng: getLanguage(),
  fallbackLng: 'ru',
  ns: ['common', 'auth'],
  defaultNS: 'common',
  interpolation: {
    escapeValue: false,
  },
});

export default i18n;
