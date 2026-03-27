import { useState } from "react";

interface FiltersProps {
  filters: string[];
  onAddFilter: (tag: string) => void;
  onRemoveFilter: (tag: string) => void;
}

export default function Filters({
  filters,
  onAddFilter,
  onRemoveFilter,
}: FiltersProps) {
  const [input, setInput] = useState("");

  const handleSubmit = (e: React.SubmitEvent) => {
    e.preventDefault();
    const tag = input.trim();
    if (tag && !filters.includes(tag)) {
      onAddFilter(tag);
      setInput("");
    }
  };

  return (
    <div className="flex flex-col gap-2">
      <form onSubmit={handleSubmit} className="flex justify-center">
        <div className="flex items-center border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-700 overflow-hidden">
          <input
            type="text"
            placeholder="Filter by tag..."
            value={input}
            onChange={(e) => setInput(e.target.value)}
            className="text-sm py-1 px-3 bg-transparent text-gray-900 dark:text-gray-100 placeholder-gray-400 dark:placeholder-gray-500 outline-none"
          />
          <button
            type="submit"
            className="bg-blue-500 text-white px-3 py-1 text-sm hover:bg-blue-600 h-full"
          >
            +
          </button>
        </div>
      </form>

      {filters.length > 0 && (
        <div className="flex flex-wrap gap-2">
          {filters.map((tag) => (
            <span
              key={tag}
              className="inline-flex items-center gap-1 bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300 px-2 py-1 rounded text-sm"
            >
              {tag}
              <button
                onClick={() => onRemoveFilter(tag)}
                className="text-blue-500 dark:text-blue-400 hover:text-blue-700 dark:hover:text-blue-200 font-bold"
              >
                ×
              </button>
            </span>
          ))}
        </div>
      )}
    </div>
  );
}
