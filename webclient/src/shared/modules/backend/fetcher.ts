import { getBackendUri } from "@/shared/config.ts";
import type { Result } from "./types.ts";

export async function fetcher<T>(relativePath: string): Promise<T | null> {
  const targetUrl = `${getBackendUri()}${relativePath}`;
  console.log("targetUrl", targetUrl);
  const request = await fetch(targetUrl);

  if (request.status === 404) {
    return null;
  }

  if (request.status >= 500) {
    throw new Error(`Internal server error: ${request.status}`);
  }

  const result = (await request.json()) as Result<T>;

  if (result.error !== undefined && result.error !== null) {
    throw new Error(result.error);
  }

  return result.data;
}
