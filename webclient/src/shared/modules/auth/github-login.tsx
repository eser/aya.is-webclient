"use client";

import * as React from "react";

import { Button } from "@/shared/components/ui/button.tsx";
import { useAuth } from "@/shared/modules/auth/auth-context.tsx";

import styles from "./github-login.module.css";

interface GitHubLoginProps {
  locale?: string;
  redirectUri?: string;
}

export function GitHubLogin({ locale = "tr", redirectUri }: GitHubLoginProps) {
  const { login, logout, isAuthenticated, user } = useAuth();

  const handleClick = () => {
    if (isAuthenticated) {
      logout();
    } else {
      login(locale, redirectUri);
    }
  };

  return (
    <Button className={styles.button} onClick={handleClick}>
      <span>{isAuthenticated ? `Çıkış (${user?.name || user?.githubHandle || "User"})` : "GitHub ile Giriş"}</span>
    </Button>
  );
}
