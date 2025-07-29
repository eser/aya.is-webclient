import { fetcher } from "../fetcher.ts";

export type UpdateProfilePageParams = {
  slug: string;
  order: number;
  cover_picture_uri?: string;
  published_at?: string;
};

export async function updateProfilePage(
  locale: string,
  profileSlug: string,
  pageId: string,
  params: UpdateProfilePageParams,
) {
  const url = `/${locale}/profiles/${profileSlug}/_pages/${pageId}`;

  const response = await fetcher(url, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(params),
  });

  return response;
}
