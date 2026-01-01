import { NavLink } from 'react-router-dom';
import type { ReactNode } from 'react';
import { motion } from 'framer-motion';

// 🔹 NavItem (для десктопа)

// 🔹 Что такое isActive
// isActive — это флаг (boolean), который React Router сам передаёт в NavLink.
// isActive === true, если текущий URL совпадает с to=""

type Props = {
  to: string;
  text: string;
  icon: ReactNode;
};

const NavItem = ({ to, text, icon }: Props) => {
  return (
    <NavLink
      to={to}
      className={({ isActive }) =>
        `
    flex items-center gap-2 transition relative ${
      isActive
        ? 'text-blue-600 dark:text-blue-400'
        : 'text-gray-600 dark:text-gray-300 hover:text-blue-500'
    }`
      }
    >
      {({ isActive }) => (
        <>
          {icon}
          {text}

          {/* Анимированная линия под активным пунктом */}
          {isActive && (
            <motion.span
              layoutId="underline"
              className="absolute left-0 -bottom-1 h-[2px] w-full bg-blue-500 rounded"
            />
          )}
        </>
      )}
    </NavLink>
  );
};

export default NavItem;
