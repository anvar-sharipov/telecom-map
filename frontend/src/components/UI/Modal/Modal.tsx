import type { ReactNode } from 'react';
import { useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';

type ModalProps = {
  open: boolean;
  onClose: () => void;

  children: ReactNode;

  closeOnBackdrop?: boolean; // 👈 разрешить закрытие кликом вне
  showBackdrop?: boolean; // 👈 false просто окно без затемнения
  width?: 'sm' | 'md' | 'lg' | 'xl';
};

const widthMap = {
  sm: 'max-w-sm',
  md: 'max-w-md',
  lg: 'max-w-lg',
  xl: 'max-w-xl',
};

const Modal = ({
  open,
  onClose,
  children,
  closeOnBackdrop = true,
  showBackdrop = true,
  width = 'md',
}: ModalProps) => {
  // ESC
  useEffect(() => {
    if (!open) return;
    audio = new Audio('/sounds/open_modal_sound.mp3');
    audio.volume = 1.0;
    audio.play().catch(() => {});

    const handler = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        onClose();
      }
    };

    document.addEventListener('keydown', handler);
    return () => document.removeEventListener('keydown', handler);
  }, [open, onClose]);

  let audio: HTMLAudioElement | null = null;

  return (
    <AnimatePresence>
      {open && (
        <motion.div
          className="fixed inset-0 z-50 flex items-center justify-center"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
        >
          {/* BACKDROP */}
          {showBackdrop && (
            <motion.div
              className="absolute inset-0 bg-black/40"
              onClick={closeOnBackdrop ? onClose : undefined}
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
            />
          )}

          {/* MODAL */}
          <motion.div
            className={`relative z-10 w-full ${widthMap[width]} mx-4 bg-white dark:bg-gray-800 rounded-2xl shadow-xl`}
            initial={{ scale: 0.95, opacity: 0, y: 10 }}
            animate={{ scale: 1, opacity: 1, y: 0 }}
            exit={{ scale: 0.95, opacity: 0, y: 10 }}
            transition={{ duration: 0.2, ease: 'easeOut' }}
            onClick={(e) => e.stopPropagation()} // 🚫 не закрывать при клике внутри
          >
            {children}
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  );
};

export default Modal;
