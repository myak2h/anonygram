import { imageUrl } from "../api";
import type { Image } from "../types";
import { ThemedView, ThemedText } from "./themed";

interface ImageCardProps {
  image: Image;
  onTagClick?: (tag: string) => void;
}

export default function ImageCard({ image, onTagClick }: ImageCardProps) {
  return (
    <ThemedView variant="card" className="rounded p-4">
      <img
        src={imageUrl(image.url)}
        alt={image.title}
        className="h-80 mb-2 mx-auto"
      />
      <ThemedText as="h3" variant="secondary" className="text-md font-semibold">
        {image.title}
      </ThemedText>
      <div className="flex items-end">
        <div className="flex-1 ">
          <div className="flex gap-2">
            {image.tags.map((tag) => (
              <ThemedText
                key={tag}
                as="span"
                variant="link"
                onClick={() => onTagClick?.(tag)}
                className="text-sm hover:underline cursor-pointer"
              >
                #{tag}
              </ThemedText>
            ))}
          </div>
        </div>
        <ThemedText variant="muted" as="p" className="text-sm">
          {new Date(image.createdAt).toLocaleString()}
        </ThemedText>
      </div>
    </ThemedView>
  );
}
