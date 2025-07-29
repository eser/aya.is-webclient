import { fetcher } from "../fetcher.ts";

export async function listProfilePages(locale: string, profileSlug: string) {
  const url = `/${locale}/profiles/${profileSlug}/_pages`;

  const response = await fetcher(url, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });

  return response;
}
