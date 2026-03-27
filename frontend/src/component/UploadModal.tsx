import { useState, useRef } from "react";
import { ThemedView, ThemedText, ThemedInput } from "./themed";

interface UploadModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
  onUpload: (title: string, tags: string[], file: File) => Promise<void>;
}

export default function UploadModal({
  isOpen,
  onClose,
  onSuccess,
  onUpload,
}: UploadModalProps) {
  const [title, setTitle] = useState("");
  const [tags, setTags] = useState("");
  const [file, setFile] = useState<File | null>(null);
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (!file) {
      setError("Please select an image");
      return;
    }

    setIsLoading(true);
    try {
      const tagList = tags
        .split(",")
        .map((t) => t.trim())
        .filter((t) => t);
      await onUpload(title, tagList, file);
      resetForm();
      onSuccess();
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Upload failed");
    } finally {
      setIsLoading(false);
    }
  };

  const resetForm = () => {
    setTitle("");
    setTags("");
    setFile(null);
    setError("");
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  const handleClose = () => {
    resetForm();
    onClose();
  };

  return (
    <ThemedView
      variant="overlay"
      className="fixed inset-0 z-50 flex items-center justify-center"
      onClick={handleClose}
    >
      <ThemedView
        variant="card"
        className="p-6 w-full h-full sm:h-auto sm:rounded-lg sm:max-w-md sm:mx-4"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex justify-between items-center mb-4">
          <ThemedText
            as="h2"
            variant="primary"
            className="text-xl font-semibold"
          >
            Upload Image
          </ThemedText>
          <button
            onClick={handleClose}
            className="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 text-3xl leading-none"
          >
            ×
          </button>
        </div>

        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <ThemedInput
            type="text"
            placeholder="Title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
          />

          <ThemedInput
            type="text"
            placeholder="Tags (comma separated)"
            value={tags}
            onChange={(e) => setTags(e.target.value)}
          />

          <ThemedInput
            ref={fileInputRef}
            type="file"
            accept="image/*"
            onChange={(e) => setFile(e.target.files?.[0] || null)}
          />

          {error && (
            <ThemedText variant="error" className="text-sm">
              {error}
            </ThemedText>
          )}

          <button
            type="submit"
            disabled={isLoading}
            className="bg-blue-500 text-white py-2 rounded hover:bg-blue-600 disabled:opacity-50"
          >
            {isLoading ? "Uploading..." : "Upload"}
          </button>
        </form>
      </ThemedView>
    </ThemedView>
  );
}
