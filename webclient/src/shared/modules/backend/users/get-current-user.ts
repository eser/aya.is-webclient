import { fetcher } from "@/shared/modules/backend/fetcher.ts";
import type { Profile } from "@/shared/modules/backend/profiles/types.ts";
import type { User } from "./types.ts";

export type CurrentUserData = User & {
  individual_profile?: Profile;
};

export async function getCurrentUser(locale: string): Promise<CurrentUserData | null> {
  const response = await fetcher<CurrentUserData>(`/${locale}/users/me`);

  return response;
}
