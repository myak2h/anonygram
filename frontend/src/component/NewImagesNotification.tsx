interface NewImagesNotificationProps {
  count: number;
  onClick: () => void;
}

export default function NewImagesNotification({
  count,
  onClick,
}: NewImagesNotificationProps) {
  if (count === 0) return null;

  return (
    <button
      onClick={onClick}
      className="fixed top-20 left-1/2 -translate-x-1/2 z-40 bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-full shadow-lg transition-all animate-bounce"
    >
      {count} new {count === 1 ? "image" : "images"} ↑
    </button>
  );
}
