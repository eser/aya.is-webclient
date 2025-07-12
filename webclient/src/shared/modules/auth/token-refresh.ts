import { getBackendUri } from "@/shared/config.ts";

interface TokenRefreshResponse {
  token: string;
  expiresAt: number;
}

let refreshPromise: Promise<string | null> | null = null;

export async function refreshToken(): Promise<string | null> {
  // If a refresh is already in progress, return the same promise
  if (refreshPromise) {
    return refreshPromise;
  }

  refreshPromise = performTokenRefresh();

  try {
    const result = await refreshPromise;
    return result;
  } finally {
    refreshPromise = null;
  }
}

async function performTokenRefresh(): Promise<string | null> {
  try {
    const currentToken = globalThis.localStorage === undefined
      ? null
      : localStorage.getItem("auth_token");

    if (!currentToken) {
      return null;
    }

    const backendUri = getBackendUri();
    const response = await fetch(`${backendUri}/tr/auth/refresh`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${currentToken}`,
      },
    });

    if (!response.ok) {
      // Token is invalid or expired
      clearAuthData();
      return null;
    }

    const data = await response.json();

    if (data.data?.token) {
      // Store new token
      if (globalThis.localStorage !== undefined) {
        localStorage.setItem("auth_token", data.data.token);
        localStorage.setItem("auth_token_expires", data.data.expiresAt.toString());
      }
      return data.data.token;
    }

    return null;
  } catch (error) {
    console.error("Token refresh failed:", error);
    return null;
  }
}

export function clearAuthData() {
  if (globalThis.localStorage !== undefined) {
    localStorage.removeItem("auth_token");
    localStorage.removeItem("auth_user");
    localStorage.removeItem("auth_token_expires");
  }
}

export function isTokenExpiringSoon(): boolean {
  if (globalThis.localStorage === undefined) {
    return false;
  }

  const expiresAt = localStorage.getItem("auth_token_expires");
  if (!expiresAt) {
    return true; // If no expiration is stored, consider it expiring
  }

  const expirationTime = parseInt(expiresAt) * 1000; // Convert to milliseconds
  const now = Date.now();
  const fiveMinutes = 5 * 60 * 1000;

  // Token expires in less than 5 minutes
  return (expirationTime - now) < fiveMinutes;
}

export function getTokenExpirationTime(): number | null {
  if (globalThis.localStorage === undefined) {
    return null;
  }

  const expiresAt = localStorage.getItem("auth_token_expires");
  if (!expiresAt) {
    return null;
  }

  return parseInt(expiresAt) * 1000; // Convert to milliseconds
}
