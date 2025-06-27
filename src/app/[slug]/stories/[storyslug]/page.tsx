import * as React from "react";
import type { Metadata, ResolvingMetadata } from "next";
import { notFound } from "next/navigation";

import { mdx } from "@/shared/lib/mdx.tsx";
import { backend } from "@/shared/modules/backend/backend.ts";
import { getTranslations } from "@/shared/modules/i18n/get-translations.tsx";
import { components } from "@/shared/components/userland/userland.ts";
import { StoryFooter } from "@/shared/components/userland/story/story-footer.tsx";

// export const revalidate = 300;

type IndexPageProps = {
  params: Promise<{
    slug: string;
    storyslug: string;
  }>;
};

// TODO(@eser) add more from https://beta.nextjs.org/docs/api-reference/metadata
export async function generateMetadata(props: IndexPageProps, _parent: ResolvingMetadata): Promise<Metadata> {
  const params = await props.params;

  const { locale } = await getTranslations();

  const profileData = await backend.getProfile(locale.code, params.slug);

  if (profileData === null) {
    notFound();
  }

  const storiesData = await backend.getProfileStories(locale.code, profileData.slug);

  if (storiesData === null) {
    notFound();
  }

  const storyData = storiesData.find((story) => story.slug === params.storyslug);
  if (storyData === undefined) {
    notFound();
  }

  return {
    title: `${profileData.title} - ${storyData.title}`,
    description: storyData.summary,
  };
}

async function IndexPage(props: IndexPageProps) {
  const params = await props.params;

  const { locale } = await getTranslations();

  const profileData = await backend.getProfile(locale.code, params.slug);

  if (profileData === null) {
    notFound();
  }

  const storiesData = await backend.getProfileStories(locale.code, profileData.slug);

  if (storiesData === null) {
    notFound();
  }

  const storyData = storiesData.find((story) => story.slug === params.storyslug);
  if (storyData === undefined) {
    notFound();
  }

  const contentText = `# ${storyData.title}\n\n${storyData.content}`;

  const mdxSource = await mdx(contentText, components);

  return (
    <div className="flex max-w-[980px] flex-col items-start">
      <article className="content">{mdxSource?.content}</article>

      <StoryFooter story={storyData} />
    </div>
  );
}

export { IndexPage as default };
