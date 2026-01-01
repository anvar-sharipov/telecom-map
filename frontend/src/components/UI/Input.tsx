import { forwardRef } from 'react';
import clsx from 'clsx';
import React from 'react';

type InputProps = {
  label?: string;
  error?: string;
} & React.InputHTMLAttributes<HTMLInputElement>;

const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, className, ...props }, ref) => {
    return (
      <div className="w-full space-y-1">
        {label && (
          <label className="block text-sm font-medium" htmlFor={props.id}>
            {label}
          </label>
        )}

        <input
          id={props.id}
          ref={ref}
          {...props}
          className={clsx(
            'w-full rounded-md border px-3 py-2 text-sm transition',
            'bg-white dark:bg-zinc-800',
            'placeholder-gray-400 dark:placeholder-gray-500',
            'focus:outline-none focus:ring-2 focus:ring-blue-500',
            error ? 'border-red-500 focus:ring-red-500' : 'border-gray-300 dark:border-zinc-700',
            className,
          )}
        />
        {error && <p className="mt-1 text-sm text-red-500">{error}</p>}
      </div>
    );
  },
);

Input.displayName = 'Input';

export default Input;
