import { useEffect, useState, useCallback } from "react";
import type { Image } from "../types";
import { fetchImages, uploadImage } from "../api";

export default function useImages() {
  const [images, setImages] = useState<Image[]>([]);
  const [filters, setFilters] = useState<string[]>([]);

  const loadImages = useCallback(() => {
    fetchImages(filters)
      .then(setImages)
      .catch((error) => {
        console.error("Error fetching images:", error);
      });
  }, [filters]);

  useEffect(() => {
    loadImages();
  }, [loadImages]);

  const addImage = useCallback((image: Image) => {
    setImages((prev) => {
      if (prev.some((img) => img.id === image.id)) return prev;
      return [image, ...prev];
    });
  }, []);

  const postImage = async (title: string, tags: string[], file: File) => {
    const newImage = await uploadImage(title, tags, file);
    addImage(newImage);
  };

  const addFilter = (tag: string) => {
    if (!filters.includes(tag)) {
      setFilters((prev) => [...prev, tag]);
    }
  };

  const removeFilter = (tag: string) => {
    setFilters((prev) => prev.filter((t) => t !== tag));
  };

  return { images, filters, postImage, addImage, addFilter, removeFilter };
}
