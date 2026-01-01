import { useSelector } from 'react-redux';
import { AnimatePresence } from 'framer-motion';
import Notification from './Notification';
import type { RootState } from '../../app/store';

export default function Notifications() {
  const notifications = useSelector((state: RootState) => state.notification.list);

  return (
    <div className="fixed z-50 space-y-2 bottom-4 left-4">
      <AnimatePresence>
        {notifications.map((n) => (
          <Notification key={n.id} {...n} />
        ))}
      </AnimatePresence>
    </div>
  );
}
