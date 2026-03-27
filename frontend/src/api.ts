import type { Image } from "./types";

const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080";

export async function fetchImages(tags: string[] = []): Promise<Image[]> {
  const query = tags.map((tag) => `tag=${encodeURIComponent(tag)}`).join("&");

  const response = await fetch(`${API_BASE_URL}/images?${query}`);

  if (!response.ok) {
    await handleErrorResponse(response);
  }

  return response.json() as Promise<Image[]>;
}

export async function uploadImage(
  title: string,
  tags: string[],
  file: File,
): Promise<Image> {
  const formData = new FormData();
  formData.append("title", title);
  formData.append("tags", tags.join(","));
  formData.append("image", file);

  const response = await fetch(`${API_BASE_URL}/uploads`, {
    method: "POST",
    body: formData,
  });

  if (!response.ok) {
    await handleErrorResponse(response);
  }

  return response.json() as Promise<Image>;
}

async function handleErrorResponse(response: Response): Promise<never> {
  let message = `Request failed with status ${response.status}`;
  try {
    const data = await response.json();
    if (data.error) {
      message = data.error;
    }
  } catch {
    // Ignore JSON parse errors
  }
  throw new Error(message);
}

export function imageUrl(filename: string): string {
  return `${API_BASE_URL}/files/${filename}`;
}
