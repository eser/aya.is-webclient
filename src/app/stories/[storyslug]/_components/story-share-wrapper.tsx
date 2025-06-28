import React, { Suspense } from "react";

import type { StoryEx } from "@/shared/modules/backend/stories/types.ts";
import { StoryShare } from "./story-share.tsx";

export type StoryShareWrapperProps = {
  story: StoryEx;
  currentUrl: string;
};

export function StoryShareWrapper(props: StoryShareWrapperProps) {
  return (
    <Suspense fallback={<div>Loading share options...</div>}>
      <StoryShare story={props.story} currentUrl={props.currentUrl} />
    </Suspense>
  );
}
