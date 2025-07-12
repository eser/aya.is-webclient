"use client";

import * as React from "react";
import { useAuth } from "./auth-context.tsx";

export function AuthErrorBanner() {
  const { authError, clearError } = useAuth();

  if (!authError) {
    return null;
  }

  return (
    <div style={{
      position: "fixed",
      top: 0,
      left: 0,
      right: 0,
      backgroundColor: "#fee",
      color: "#c00",
      padding: "1rem",
      textAlign: "center",
      borderBottom: "1px solid #fcc",
      zIndex: 1000,
      display: "flex",
      justifyContent: "center",
      alignItems: "center",
      gap: "1rem",
    }}>
      <span>{authError}</span>
      <button
        onClick={clearError}
        style={{
          background: "none",
          border: "1px solid #c00",
          color: "#c00",
          padding: "0.25rem 0.5rem",
          cursor: "pointer",
          borderRadius: "0.25rem",
        }}
      >
        Dismiss
      </button>
    </div>
  );
}
