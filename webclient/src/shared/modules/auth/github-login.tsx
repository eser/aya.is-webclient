"use client";

import * as React from "react";

import { Button } from "@/shared/components/ui/button.tsx";

import styles from "./github-login.module.css";

export function GitHubLogin() {
  return (
    <Button className={styles.button} asChild>
      <div>
        <span>GitHub ile Giriş</span>
      </div>
    </Button>
  );
}
