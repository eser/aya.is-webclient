import { fetcher } from "@/shared/modules/backend/fetcher.ts";

import type { StoryEx } from "../stories/types.ts";

export type GetProfileStoryData = StoryEx;

export async function getProfileStory(locale: string, _slug: string, storyslug: string) {
  // const response = await fetcher<GetProfileStoryData>(`/${locale}/profiles/${slug}/stories/${storyslug}`);
  const response = await fetcher<GetProfileStoryData>(`/${locale}/stories/${storyslug}`);

  return response;
}
