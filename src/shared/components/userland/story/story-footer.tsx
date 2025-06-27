import React from "react";
// import Image from "next/image";

import type { StoryEx } from "@/shared/modules/backend/stories/types.ts";
import { getTranslations } from "@/shared/modules/i18n/get-translations.tsx";
import { SiteLink } from "@/shared/components/userland/site-link/site-link.tsx";
import styles from "./story-footer.module.css";

export type StoryFooterProps = {
  story: StoryEx;
};

export async function StoryFooter(props: StoryFooterProps) {
  if (!props.story.author_profile) {
    return null;
  }

  const { t } = await getTranslations();

  return (
    <div className={styles["story-footer"]}>
      {props.story.author_profile.profile_picture_uri && (
        <div className={styles["author-avatar-container"]}>
          <img
            src={props.story.author_profile.profile_picture_uri}
            alt={props.story.author_profile.title ?? "Author image"}
            width={72}
            height={72}
            className={styles["author-avatar-image"]}
          />
        </div>
      )}
      <div className={styles["story-details"]}>
        <div className={styles["author-details"]}>
          <div className={styles["written-by-text"]}>{t("StoryPage", "Written by")}</div>
          <div className={styles["author-name"]}>
            <SiteLink href={`/${props.story.author_profile.slug}`}>
              {props.story.author_profile.title}
            </SiteLink>
          </div>
        </div>
      </div>
    </div>
  );
}
