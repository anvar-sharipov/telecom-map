import { NavLink } from 'react-router-dom';
import { motion } from 'framer-motion';
import type { LucideIcon } from 'lucide-react';

type Props = {
  to: string;
  text: string;
  icon?: LucideIcon;
  setOpen: (v: boolean) => void;
};

const MobileLink = ({ to, text, icon: Icon, setOpen }: Props) => {
  return (
    <NavLink to={to} onClick={() => setOpen(false)}>
      {({ isActive }) => (
        <motion.div
          whileTap={{ scale: 0.97 }}
          className={`
            flex items-center gap-3 px-4 py-3 rounded-xl
            text-sm font-medium
            transition-colors
            ${
              isActive
                ? 'bg-indigo-100 text-indigo-700 dark:bg-zinc-800'
                : 'hover:bg-gray-100 dark:hover:bg-zinc-800'
            }
          `}
        >
          {Icon && <Icon className="w-5 h-5" />}
          <span>{text}</span>
        </motion.div>
      )}
    </NavLink>
  );
};

export default MobileLink;
