import { useTranslation } from 'react-i18next';
import { motion } from 'framer-motion';

export function LanguageSwitcher() {
  const { i18n } = useTranslation();

  const changeLang = (lang: 'ru' | 'tm') => {
    i18n.changeLanguage(lang);
    localStorage.setItem('lang', lang);
  };

  return (
    <div className="flex gap-2">
      <div className="relative">
        <button
          className={`${
            i18n.language === 'ru' ? 'font-bold' : ''
          } text-blue-600 dark:text-blue-400 hover:text-blue-500`}
          onClick={() => changeLang('ru')}
        >
          RU
        </button>

        {i18n.language === 'ru' && (
          <motion.span
            className="absolute left-0 -bottom-1 h-[2px] w-full bg-blue-500 rounded"
            layoutId="lang-underline"
          />
        )}
      </div>

      <div className="relative">
        <button
          className={`${
            i18n.language === 'tm' ? 'font-bold' : ''
          } text-blue-600 dark:text-blue-400 hover:text-blue-500`}
          onClick={() => changeLang('tm')}
        >
          TM
        </button>

        {i18n.language === 'tm' && (
          <motion.span
            className="absolute left-0 -bottom-1 h-[2px] w-full bg-blue-500 rounded"
            layoutId="lang-underline"
          />
        )}
      </div>
    </div>
  );
}
