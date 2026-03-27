import { forwardRef, type InputHTMLAttributes } from "react";

type ThemedInputProps = InputHTMLAttributes<HTMLInputElement>;

const ThemedInput = forwardRef<HTMLInputElement, ThemedInputProps>(
  ({ className = "", ...props }, ref) => {
    return (
      <input
        ref={ref}
        className={`border border-gray-300 dark:border-gray-600 rounded px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 placeholder-gray-400 dark:placeholder-gray-500 ${className}`}
        {...props}
      />
    );
  },
);

ThemedInput.displayName = "ThemedInput";

export default ThemedInput;
