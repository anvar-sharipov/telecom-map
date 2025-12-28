import { useTranslation } from 'react-i18next';

export function LanguageSwitcher() {
  const { i18n } = useTranslation();

  const changeLang = (lang: 'ru' | 'tm') => {
    i18n.changeLanguage(lang);
    localStorage.setItem('lang', lang);
  };

  return (
    <div className="flex gap-2">
      <button
        className={i18n.language === 'ru' ? 'font-bold' : ''}
        onClick={() => changeLang('ru')}
      >
        RU
      </button>
      <button
        className={i18n.language === 'tm' ? 'font-bold' : ''}
        onClick={() => changeLang('tm')}
      >
        TM
      </button>
    </div>
  );
}
