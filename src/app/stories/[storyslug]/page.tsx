import * as React from "react";
import type { Metadata, ResolvingMetadata } from "next";
import { notFound } from "next/navigation";

import { mdx } from "@/shared/lib/mdx.tsx";
import { backend } from "@/shared/modules/backend/backend.ts";
import { getTranslations } from "@/shared/modules/i18n/get-translations.tsx";

import { components } from "@/shared/components/userland/userland.ts";
import { StoryMetadata } from "./_components/story-metadata.tsx";
import { StoryFooter } from "@/shared/components/userland/story/story-footer.tsx";
import { StoryShareWrapper } from "./_components/story-share-wrapper.tsx";

export const revalidate = 300;

type IndexPageProps = {
  params: Promise<{
    storyslug: string;
  }>;
};

// TODO(@eser) add more from https://beta.nextjs.org/docs/api-reference/metadata
export async function generateMetadata(props: IndexPageProps, _parent: ResolvingMetadata): Promise<Metadata> {
  const params = await props.params;

  const { locale } = await getTranslations();

  const storyData = await backend.getStory(locale.code, params.storyslug);

  if (storyData === null) {
    notFound();
  }

  return {
    title: storyData.title,
    description: storyData.summary,
  };
}

async function IndexPage(props: IndexPageProps) {
  const params = await props.params;

  const { locale } = await getTranslations();

  const storyData = await backend.getStory(locale.code, params.storyslug);

  if (storyData === null) {
    notFound();
  }

  const contentText = `${storyData.content}`;

  const mdxSource = await mdx(contentText, components);

  // Get current URL for sharing
  const baseUrl =
    process.env.NEXT_PUBLIC_SITE_URL?.replace(/\/$/, "") ?? "https://aya.is";
  const currentUrl = `${baseUrl}/stories/${storyslug}`;

  return (
    <section className="container px-4 py-8 mx-auto">
      <div className="content">
        <h2>{storyData.title}</h2>

        <StoryMetadata story={storyData} />

        <article>{mdxSource?.content}</article>

        <StoryShareWrapper story={storyData} currentUrl={currentUrl} />

        <StoryFooter story={storyData} />
      </div>
    </section>
  );
}

export { IndexPage as default };
