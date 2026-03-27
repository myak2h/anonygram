import { forwardRef, type HTMLAttributes } from "react";

type ViewVariant = "primary" | "secondary" | "card" | "overlay";

interface ThemedViewProps extends HTMLAttributes<HTMLDivElement> {
  variant?: ViewVariant;
}

const variantClasses: Record<ViewVariant, string> = {
  primary: "bg-white dark:bg-gray-900",
  secondary: "bg-gray-50 dark:bg-gray-800",
  card: "bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700",
  overlay: "bg-black/50",
};

const ThemedView = forwardRef<HTMLDivElement, ThemedViewProps>(
  ({ variant = "primary", className = "", children, ...props }, ref) => {
    return (
      <div
        ref={ref}
        className={`${variantClasses[variant]} ${className}`}
        {...props}
      >
        {children}
      </div>
    );
  },
);

ThemedView.displayName = "ThemedView";

export default ThemedView;
