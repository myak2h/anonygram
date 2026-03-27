import type { HTMLAttributes } from "react";

type TextVariant = "primary" | "secondary" | "muted" | "error" | "link";
type TextElement = "span" | "p" | "h1" | "h2" | "h3" | "label";

interface ThemedTextProps extends HTMLAttributes<HTMLElement> {
  variant?: TextVariant;
  as?: TextElement;
}

const variantClasses: Record<TextVariant, string> = {
  primary: "text-gray-900 dark:text-gray-100",
  secondary: "text-gray-700 dark:text-gray-200",
  muted: "text-gray-500 dark:text-gray-400",
  error: "text-red-500 dark:text-red-400",
  link: "text-blue-500 dark:text-blue-400 hover:text-blue-700 dark:hover:text-blue-300",
};

function ThemedText({
  variant = "primary",
  as: Component = "span",
  className = "",
  children,
  ...props
}: ThemedTextProps) {
  return (
    <Component
      className={`${variantClasses[variant]} ${className}`}
      {...props}
    >
      {children}
    </Component>
  );
}

export default ThemedText;
