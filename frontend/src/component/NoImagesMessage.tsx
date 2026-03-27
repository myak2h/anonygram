import { ThemedText } from "./themed";

export default function NoImagesMessage() {
  return (
    <div className="flex flex-col items-center justify-center py-16 text-center">
      <span className="text-6xl mb-4">📷</span>
      <ThemedText
        as="h2"
        variant="primary"
        className="text-xl font-semibold mb-2"
      >
        No images yet
      </ThemedText>
      <ThemedText variant="muted">
        Click the + button to upload your first image
      </ThemedText>
    </div>
  );
}
