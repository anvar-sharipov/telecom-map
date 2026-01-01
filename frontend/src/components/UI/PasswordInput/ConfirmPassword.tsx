import { useEffect, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import PasswordInput from './PasswordInput';
import { useTranslation } from 'react-i18next';

type Props = {
  password: string;
  value: string;
  onChange: (v: string) => void;
};

export default function ConfirmPassword({ password, value, onChange }: Props) {
  const { t } = useTranslation('auth');
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!value) {
      setError(null);
    } else if (value !== password) {
      setError('passwords do not match');
    } else {
      setError(null);
    }
  }, [value, password]);

  return (
    <div className="space-y-1">
      <PasswordInput
        label={t('confirm password')}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        autoComplete="new-password"
        // className={className}
      />

      <AnimatePresence>
        {error && (
          <motion.p
            initial={{ opacity: 0, y: -4 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -4 }}
            className="text-sm text-red-500"
          >
            {error}
          </motion.p>
        )}
      </AnimatePresence>
    </div>
  );
}
