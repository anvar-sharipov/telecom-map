import { forwardRef } from 'react';
import type { ReactNode } from 'react';
import clsx from 'clsx';
import { motion, useAnimation } from 'framer-motion';
import type { HTMLMotionProps } from 'framer-motion';
import { useEffect } from 'react';

type ButtonProps = {
  children: ReactNode;
  variant?: 'primary' | 'secondary' | 'outline';
  size?: 'sm' | 'md' | 'lg';
  icon?: ReactNode;
  isLoading?: boolean;
} & HTMLMotionProps<'button'>;

const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ children, variant = 'primary', size = 'md', icon, className, isLoading, ...props }, ref) => {
    const baseStyles =
      'inline-flex items-center justify-center font-medium rounded-md transition-all focus:outline-none focus:ring-2 focus:ring-offset-1 relative';

    const variants = {
      primary: 'bg-blue-600 text-white hover:bg-blue-700 dark:bg-blue-500 dark:hover:bg-blue-600',
      secondary:
        'bg-gray-200 text-gray-900 hover:bg-gray-300 dark:bg-zinc-700 dark:text-gray-100 dark:hover:bg-zinc-600',
      outline:
        'border border-gray-300 text-gray-900 hover:bg-gray-100 dark:border-zinc-700 dark:text-gray-100 dark:hover:bg-zinc-800',
    };

    const sizes = {
      sm: 'px-3 py-1 text-sm',
      md: 'px-4 py-2 text-base',
      lg: 'px-6 py-3 text-lg',
    };

    // Анимация для loading
    const controls = useAnimation();

    useEffect(() => {
      if (isLoading) {
        controls.start({
          scale: [1, 1.05, 1],
          transition: { duration: 0.8, repeat: Infinity, repeatType: 'loop', ease: 'easeInOut' },
        });
      } else {
        controls.stop();
        controls.set({ scale: 1 });
      }
    }, [isLoading, controls]);

    return (
      <motion.button
        ref={ref}
        {...props}
        whileTap={{ scale: 0.95 }}
        animate={controls}
        className={clsx(baseStyles, variants[variant], sizes[size], className)}
        disabled={isLoading || props.disabled}
      >
        {icon && <span className="mr-2">{icon}</span>}
        {isLoading ? 'Loading...' : children}
      </motion.button>
    );
  },
);

Button.displayName = 'Button';

export default Button;
