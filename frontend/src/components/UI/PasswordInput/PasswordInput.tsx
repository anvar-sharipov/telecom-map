import { forwardRef, useState } from 'react';
import { Eye, EyeOff } from 'lucide-react';
import clsx from 'clsx';

type PasswordInputProps = {
  label?: string;
  error?: string;
} & React.InputHTMLAttributes<HTMLInputElement>;

const PasswordInput = forwardRef<HTMLInputElement, PasswordInputProps>(
  ({ label, error, className, ...props }, ref) => {
    const [show, setShow] = useState(false);

    return (
      <div className="w-full space-y-1">
        {label && (
          <label
            className="block text-sm font-medium text-gray-700 dark:text-gray-300"
            htmlFor={props.id}
          >
            {label}
          </label>
        )}

        <div className="relative">
          <input
            id={props.id}
            ref={ref}
            type={show ? 'text' : 'password'}
            {...props}
            className={clsx(
              'w-full rounded-md px-3 py-2 pr-10 text-sm transition',
              'bg-white dark:bg-zinc-800',
              'placeholder-gray-400 dark:placeholder-gray-500',
              'border focus:outline-none focus:ring-2',
              error
                ? 'border-red-500 focus:ring-red-500'
                : 'border-gray-300 dark:border-zinc-700 focus:ring-blue-500',
              className,
            )}
          />

          <button
            type="button"
            onClick={() => setShow(!show)}
            className="absolute inset-y-0 flex items-center text-gray-500 right-2 hover:text-gray-700 dark:hover:text-gray-300"
          >
            {show ? <EyeOff size={18} /> : <Eye size={18} />}
          </button>
        </div>

        {error && <p className="text-sm text-red-500">{error}</p>}
      </div>
    );
  },
);

PasswordInput.displayName = 'PasswordInput';

export default PasswordInput;
