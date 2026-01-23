import { Loader2 } from 'lucide-react';
import { useTranslation } from 'react-i18next';

const LoadingSpin = () => {
  const { t } = useTranslation('common');
  return (
    <div className="flex flex-col items-center justify-center gap-2 py-6 text-gray-600">
      <Loader2 className="w-6 h-6 text-blue-600 animate-spin" />
      <span className="text-sm font-medium">{t('loading')}</span>
    </div>
  );
};

export default LoadingSpin;
