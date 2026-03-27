interface HeaderProps {
  onAddClick: () => void;
  isDark: boolean;
  onToggleTheme: () => void;
}

export default function Header({
  onAddClick,
  isDark,
  onToggleTheme,
}: HeaderProps) {
  return (
    <header className="fixed top-0 left-1/2 -translate-x-1/2 z-50 max-w-lg w-full h-16 flex items-center justify-between p-2 border-b border-gray-300 dark:border-gray-700 bg-white dark:bg-gray-900">
      <h1 className="text-2xl font-bold text-purple-700 dark:text-purple-400">
        Anony<span className="text-blue-500 dark:text-blue-400">Gram</span>
      </h1>
      <div className="flex items-center gap-2">
        <button
          onClick={onToggleTheme}
          className="p-2 rounded hover:bg-gray-100 dark:hover:bg-gray-800 text-gray-600 dark:text-gray-300"
          aria-label="Toggle theme"
        >
          {isDark ? "☀️" : "🌙"}
        </button>
        <button
          onClick={onAddClick}
          className="bg-purple-500 text-white text-2xl px-4 p-1 rounded cursor-pointer"
        >
          +
        </button>
      </div>
    </header>
  );
}
