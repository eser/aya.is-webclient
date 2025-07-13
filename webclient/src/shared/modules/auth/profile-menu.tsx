"use client";

import * as React from "react";
import { Button } from "@/shared/components/ui/button.tsx";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/shared/components/ui/dropdown-menu.tsx";
import { useAuth } from "@/shared/modules/auth/auth-context.tsx";
import { useRouter } from "next/navigation";
import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";

interface ProfileMenuProps {
  className?: string;
}

export function ProfileMenu({ className }: ProfileMenuProps) {
  const { user, logout, isAuthenticated } = useAuth();
  const router = useRouter();
  const { t } = useTranslations();

  if (!isAuthenticated || !user) {
    return null;
  }

  const handleProfileClick = () => {
    router.push(`/@${user.githubHandle || user.id}`);
  };

  const handleLogout = async () => {
    await logout();
  };

  const avatarUrl = user.githubHandle
    ? `https://github.com/${user.githubHandle}.png?size=32`
    : `https://api.dicebear.com/7.x/initials/svg?seed=${user.name}`;

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" className={`relative h-8 w-8 rounded-full ${className}`}>
          <img
            className="h-8 w-8 rounded-full"
            src={avatarUrl}
            alt={user.name}
          />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-56" align="end" forceMount>
        <DropdownMenuLabel className="font-normal">
          <div className="flex flex-col space-y-1">
            <p className="text-sm font-medium leading-none">{user.name}</p>
            {user.email && (
              <p className="text-xs leading-none text-muted-foreground">
                {user.email}
              </p>
            )}
          </div>
        </DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={handleProfileClick}>
          {t("Auth", "profile")}
        </DropdownMenuItem>
        <DropdownMenuItem onClick={handleLogout}>
          {t("Auth", "logout")}
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
