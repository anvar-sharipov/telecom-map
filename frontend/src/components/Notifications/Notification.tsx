import { motion } from 'framer-motion';
import { useDispatch } from 'react-redux';
import { removeNotification } from './notificationSlice';
import { useEffect } from 'react';
import Snowfall from 'react-snowfall';

interface Props {
  id: string;
  message: string;
  type: 'success' | 'error' | 'info';
}

export default function Notification({ id, message, type }: Props) {
  const dispatch = useDispatch();

  useEffect(() => {
    let audio: HTMLAudioElement | null = null;

    if (type === 'success') {
      audio = new Audio('/sounds/success-notification.wav');
    } else if (type === 'error') {
      audio = new Audio('/sounds/error2.mp3');
    }

    if (audio) {
      audio.volume = 0.9;
      audio.play().catch(() => {});
    }

    return () => {
      if (audio) {
        audio.pause();
        audio.currentTime = 0;
      }
    };
  }, [type]);

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, x: 100 }}
      transition={{ duration: 0.7 }}
      className={`px-4 py-3 rounded shadow-md text-white mb-2 cursor-pointer
        ${type === 'success' && 'bg-green-500'}
        ${type === 'error' && 'bg-red-500'}
        ${type === 'info' && 'bg-blue-500'}
      `}
      onClick={() => dispatch(removeNotification(id))}
    >
      {message}
    </motion.div>
  );
}
