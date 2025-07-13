import * as React from "react";
import type { Metadata } from "next";

import { backend } from "@/shared/modules/backend/backend.ts";
import { getTranslations } from "@/shared/modules/i18n/get-translations.tsx";
import { Story } from "@/shared/components/userland/story/story.tsx";

export async function generateMetadata(): Promise<Metadata> {
  const { t } = await getTranslations();

  return {
    title: t("Layout", "Articles"),
  };
}

async function IndexPage() {
  const { t, locale } = await getTranslations();

  const stories = await backend.getStoriesByKinds(locale.code, ["article"]);

  return (
    <section className="container px-4 py-8 mx-auto">
      <div className="content">
        <h2>{t("Layout", "Articles")}</h2>

        {stories !== null && stories.length > 0
          ? (
            <div className="divide-y divide-border">
              {stories.map((story) => (
                <Story
                  key={story.id}
                  story={story}
                />
              ))}
            </div>
          )
          : (
            <article>
              {/* TODO(@eser) Add a specific translation key e.g., t("NewsPage", "NoNewsFound") */}
              {t("Layout", "Content not yet available.")}
            </article>
          )}
      </div>
    </section>
  );
}

export { IndexPage as default };
