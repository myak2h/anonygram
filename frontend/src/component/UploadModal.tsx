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
  const [tagInput, setTagInput] = useState("");
  const [tags, setTags] = useState<string[]>([]);
  const [file, setFile] = useState<File | null>(null);
  const [preview, setPreview] = useState<string | null>(null);
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  if (!isOpen) return null;

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0] || null;
    setFile(selectedFile);
    if (selectedFile) {
      const url = URL.createObjectURL(selectedFile);
      setPreview(url);
    } else {
      setPreview(null);
    }
  };

  const handleAddTag = () => {
    const tag = tagInput.trim();
    if (tag && !tags.includes(tag)) {
      setTags((prev) => [...prev, tag]);
      setTagInput("");
    }
  };

  const handleTagKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      e.preventDefault();
      handleAddTag();
    }
  };

  const handleRemoveTag = (tag: string) => {
    setTags((prev) => prev.filter((t) => t !== tag));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (!file) {
      setError("Please select an image");
      return;
    }

    setIsLoading(true);
    try {
      await onUpload(title, tags, file);
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
    setTagInput("");
    setTags([]);
    setFile(null);
    setPreview(null);
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
        className="p-6 w-full h-full sm:h-auto sm:rounded-lg sm:max-w-md sm:mx-4 overflow-y-auto"
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
          <div
            className="border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg p-4 text-center cursor-pointer hover:border-blue-500 dark:hover:border-blue-400 transition-colors"
            onClick={() => fileInputRef.current?.click()}
          >
            {preview ? (
              <img
                src={preview}
                alt="Preview"
                className="max-h-48 mx-auto rounded"
              />
            ) : (
              <div className="py-8">
                <ThemedText variant="muted" className="text-4xl mb-2 block">
                  📷
                </ThemedText>
                <ThemedText variant="muted">
                  Click to select an image
                </ThemedText>
              </div>
            )}
            <input
              ref={fileInputRef}
              type="file"
              accept="image/*"
              onChange={handleFileChange}
              className="hidden"
            />
          </div>
          <ThemedInput
            type="text"
            placeholder="Title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
          />
          <div className="flex flex-col gap-2">
            <div className="flex items-center border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-700 overflow-hidden">
              <input
                type="text"
                placeholder="Add tags..."
                value={tagInput}
                onChange={(e) => setTagInput(e.target.value)}
                onKeyDown={handleTagKeyDown}
                className="flex-1 text-sm py-2 px-3 bg-transparent text-gray-900 dark:text-gray-100 placeholder-gray-400 dark:placeholder-gray-500 outline-none"
              />
              <button
                type="button"
                onClick={handleAddTag}
                className="bg-blue-500 text-white px-3 py-2 text-sm hover:bg-blue-600 h-full"
              >
                +
              </button>
            </div>

            {tags.length > 0 && (
              <div className="flex flex-wrap gap-2">
                {tags.map((tag) => (
                  <span
                    key={tag}
                    className="inline-flex items-center gap-1 bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300 px-2 py-1 rounded text-sm"
                  >
                    {tag}
                    <button
                      type="button"
                      onClick={() => handleRemoveTag(tag)}
                      className="text-blue-500 dark:text-blue-400 hover:text-blue-700 dark:hover:text-blue-200 font-bold"
                    >
                      ×
                    </button>
                  </span>
                ))}
              </div>
            )}
          </div>

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
