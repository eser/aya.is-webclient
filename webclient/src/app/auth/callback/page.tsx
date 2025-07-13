"use client";

import * as React from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { getBackendUri } from "@/shared/config.ts";

export default function AuthCallbackPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [status, setStatus] = React.useState<"loading" | "success" | "error">("loading");
  const [errorMessage, setErrorMessage] = React.useState<string>("");

  React.useEffect(() => {
    const handleCallback = async () => {
      try {
        const code = searchParams.get("code");
        const state = searchParams.get("state");

        if (!code) {
          throw new Error("Authorization code not found");
        }

        // Exchange code for token by calling backend callback endpoint
        const backendUri = getBackendUri();
        const response = await fetch(
          `${backendUri}/tr/auth/github/callback?code=${code}&state=${state || ""}`,
          {
            method: "GET",
            headers: {
              "Content-Type": "application/json",
            },
          }
        );

        if (!response.ok) {
          throw new Error("Failed to authenticate");
        }

        const data = await response.json();

        if (!data.data?.token) {
          throw new Error("No token received from server");
        }

        // Store JWT token in localStorage for now
        // TODO(@eser) Use secure httpOnly cookie in production
        if (globalThis.localStorage !== undefined) {
          localStorage.setItem("auth_token", data.data.token);
          if (data.data.user) {
            localStorage.setItem("auth_user", JSON.stringify(data.data.user));
          }
          // Store token expiration (24 hours from now)
          const expiresAt = Math.floor(Date.now() / 1000) + (24 * 60 * 60);
          localStorage.setItem("auth_token_expires", expiresAt.toString());
        }

        setStatus("success");

        // Redirect to home or intended destination
        setTimeout(() => {
          router.push("/");
        }, 1000);
      } catch (error) {
        console.error("Auth callback error:", error);
        setStatus("error");
        setErrorMessage(error instanceof Error ? error.message : "Authentication failed");
      }
    };

    handleCallback();
  }, [searchParams, router]);

  return (
    <div style={{
      display: "flex",
      flexDirection: "column",
      alignItems: "center",
      justifyContent: "center",
      minHeight: "100vh",
      gap: "1rem"
    }}>
      {status === "loading" && (
        <>
          <div>Kimlik doğrulanıyor...</div>
          <div>Lütfen bekleyin</div>
        </>
      )}

      {status === "success" && (
        <>
          <div>✓ Başarıyla giriş yapıldı!</div>
          <div>Ana sayfaya yönlendiriliyorsunuz...</div>
        </>
      )}

      {status === "error" && (
        <>
          <div>✗ Kimlik doğrulama başarısız</div>
          <div>{errorMessage}</div>
          <button
            onClick={() => router.push("/")}
            style={{
              marginTop: "1rem",
              padding: "0.5rem 1rem",
              cursor: "pointer",
            }}
          >
            Ana Sayfaya Dön
          </button>
        </>
      )}
    </div>
  );
}
