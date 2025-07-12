import { getBackendUri } from "@/shared/config.ts";
import { isTokenExpiringSoon, refreshToken } from "@/shared/modules/auth/token-refresh.ts";
import type { Result } from "./types.ts";

export async function fetcher<T>(relativePath: string): Promise<T | null> {
  const targetUrl = `${getBackendUri()}${relativePath}`;
  console.log("targetUrl", targetUrl);

  // Get auth token from localStorage (only available on client)
  let authToken: string | null = null;
  if (globalThis.localStorage !== undefined) {
    authToken = localStorage.getItem("auth_token");
  }

  // Check if token needs refresh
  if (authToken && isTokenExpiringSoon()) {
    const newToken = await refreshToken();
    if (newToken) {
      authToken = newToken;
    }
  }

  const headers: HeadersInit = {
    "Content-Type": "application/json",
  };

  if (authToken) {
    headers["Authorization"] = `Bearer ${authToken}`;
  }

  const request = await fetch(targetUrl, {
    headers,
  });

  // If unauthorized, try to refresh token and retry once
  if (request.status === 401 && authToken) {
    const newToken = await refreshToken();
    if (newToken) {
      headers["Authorization"] = `Bearer ${newToken}`;
      const retryRequest = await fetch(targetUrl, { headers });

      if (retryRequest.status === 404) {
        return null;
      }

      if (retryRequest.status >= 500) {
        throw new Error(`Internal server error: ${retryRequest.status}`);
      }

      const retryResult = (await retryRequest.json()) as Result<T>;

      if (retryResult.error !== undefined && retryResult.error !== null) {
        throw new Error(retryResult.error);
      }

      return retryResult.data;
    }
  }

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
