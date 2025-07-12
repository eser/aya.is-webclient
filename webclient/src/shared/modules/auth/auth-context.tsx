"use client";

import * as React from "react";
import { useRouter } from "next/navigation";
import { getBackendUri } from "@/shared/config.ts";
import { getTokenExpirationTime, refreshToken, clearAuthData } from "./token-refresh.ts";

interface User {
  id: string;
  name: string;
  email?: string;
  githubHandle?: string;
}

interface AuthContextType {
  user: User | null;
  token: string | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  authError: string | null;
  login: (locale?: string, redirectUri?: string) => void;
  logout: () => Promise<void>;
  refreshAuth: () => void;
  clearError: () => void;
}

const AuthContext = React.createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const [user, setUser] = React.useState<User | null>(null);
  const [token, setToken] = React.useState<string | null>(null);
  const [isLoading, setIsLoading] = React.useState(true);
  const [authError, setAuthError] = React.useState<string | null>(null);

  // Load auth state from localStorage on mount
  React.useEffect(() => {
    const loadAuthState = () => {
      if (globalThis.localStorage === undefined) {
        setIsLoading(false);
        return;
      }

      try {
        const storedToken = localStorage.getItem("auth_token");
        const storedUser = localStorage.getItem("auth_user");

        if (storedToken && storedUser) {
          setToken(storedToken);
          setUser(JSON.parse(storedUser));
        }
      } catch (error) {
        console.error("Failed to load auth state:", error);
        // Clear invalid data
        clearAuthData();
      } finally {
        setIsLoading(false);
      }
    };

    loadAuthState();
  }, []);

  // Set up automatic token refresh
  React.useEffect(() => {
    if (!token) return;

    const checkAndRefresh = async () => {
      const expirationTime = getTokenExpirationTime();
      if (!expirationTime) return;

      const now = Date.now();
      const timeUntilExpiry = expirationTime - now;

      // Refresh 5 minutes before expiration
      const refreshTime = Math.max(0, timeUntilExpiry - (5 * 60 * 1000));

      if (refreshTime === 0) {
        // Token needs immediate refresh
        const newToken = await refreshToken();
        if (newToken) {
          setToken(newToken);
        } else {
          // Failed to refresh, clear auth state
          clearAuthData();
          setUser(null);
          setToken(null);
          setAuthError("Your session has expired. Please log in again.");
        }
      }

      return refreshTime;
    };

    let timeoutId: NodeJS.Timeout;

    const setupRefreshTimer = async () => {
      const refreshTime = await checkAndRefresh();
      if (refreshTime && refreshTime > 0) {
        timeoutId = setTimeout(async () => {
          const newToken = await refreshToken();
          if (newToken) {
            setToken(newToken);
            // Set up next refresh
            setupRefreshTimer();
          } else {
            // Failed to refresh, clear auth state
            clearAuthData();
            setUser(null);
            setToken(null);
            setAuthError("Your session has expired. Please log in again.");
          }
        }, refreshTime);
      }
    };

    setupRefreshTimer();

    return () => {
      if (timeoutId) {
        clearTimeout(timeoutId);
      }
    };
  }, [token]);

  const login = React.useCallback((locale: string = "tr", redirectUri?: string) => {
    const backendUri = getBackendUri();
    const currentUrl = globalThis.location?.href ?? "";
    const finalRedirectUri = redirectUri || currentUrl;

    // Redirect to backend OAuth endpoint
    const loginUrl = `${backendUri}/${locale}/auth/github/login${finalRedirectUri ? `?redirect_uri=${encodeURIComponent(finalRedirectUri)}` : ""}`;
    if (globalThis.location !== undefined) {
      globalThis.location.href = loginUrl;
    }
  }, []);

  const logout = React.useCallback(async () => {
    try {
      // Call backend logout endpoint
      const backendUri = getBackendUri();
      await fetch(`${backendUri}/tr/auth/logout`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          ...(token ? { Authorization: `Bearer ${token}` } : {}),
        },
      });
    } catch (error) {
      console.error("Logout error:", error);
    } finally {
      // Clear local state regardless of backend response
      clearAuthData();
      setUser(null);
      setToken(null);
      router.push("/");
    }
  }, [token, router]);

  const refreshAuth = React.useCallback(() => {
    // Re-load auth state from localStorage
    const storedToken = localStorage.getItem("auth_token");
    const storedUser = localStorage.getItem("auth_user");

    if (storedToken && storedUser) {
      try {
        setToken(storedToken);
        setUser(JSON.parse(storedUser));
      } catch (error) {
        console.error("Failed to refresh auth:", error);
      }
    }
  }, []);

  const clearError = React.useCallback(() => {
    setAuthError(null);
  }, []);

  const value = React.useMemo(
    () => ({
      user,
      token,
      isLoading,
      isAuthenticated: !!user && !!token,
      authError,
      login,
      logout,
      refreshAuth,
      clearError,
    }),
    [user, token, isLoading, authError, login, logout, refreshAuth, clearError]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = React.useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
