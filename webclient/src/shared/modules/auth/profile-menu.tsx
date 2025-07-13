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
import { useTranslations } from "@/shared/modules/i18n/use-translations.tsx";

interface ProfileMenuProps {
  className?: string;
}

export function ProfileMenu({ className }: ProfileMenuProps) {
  const { user, logout, isAuthenticated } = useAuth();
  const { t } = useTranslations();

  if (!isAuthenticated || !user) {
    return null;
  }

  const handleLogout = async () => {
    await logout();
  };

  const fallbackAvatarUrl = `https://api.dicebear.com/7.x/initials/svg?seed=${encodeURIComponent(user.name)}`;
  const avatarUrl = user.githubHandle
    ? `https://github.com/${user.githubHandle}.png?size=32`
    : fallbackAvatarUrl;

  const [imageError, setImageError] = React.useState(false);

  const handleImageError = () => {
    setImageError(true);
  };


  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" className={`relative h-8 w-8 rounded-full p-0 ${className}`}>
          <img
            className="h-8 w-8 rounded-full object-cover"
            src={imageError ? fallbackAvatarUrl : avatarUrl}
            alt={user.name}
            onError={handleImageError}
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
        <DropdownMenuItem onClick={handleLogout}>
          {t("Auth", "logout")}
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
