import { getCustomDomain } from "./profiles/get-custom-domain.ts";
import { getProfilesByKinds } from "./profiles/get-profiles-by-kinds.ts";
import { getProfile } from "./profiles/get-profile.ts";
import { getProfilePage } from "./profiles/get-profile-page.ts";
import { getProfileContributions } from "./profiles/get-profile-contributions.ts";
import { getProfileMembers } from "./profiles/get-profile-members.ts";
import { getProfileStories } from "./profiles/get-profile-stories.ts";
import { getProfileStory } from "./profiles/get-profile-story.ts";
import { checkProfileSlug } from "./profiles/check-profile-slug.ts";
import { createProfile } from "./profiles/create-profile.ts";
import { getStory } from "./stories/get-story.ts";
import { getStoriesByKinds } from "./stories/get-stories-by-kinds.ts";
import { getSpotlight } from "./site/get-spotlight.ts";
import { handleAuthCallback } from "./auth/handle-callback.ts";
import { getCurrentSession } from "./sessions/get-current-session.ts";

export const backend = {
  // Profiles
  getCustomDomain,
  getProfilesByKinds,
  getProfile,
  getProfilePage,
  getProfileContributions,
  getProfileMembers,
  getProfileStories,
  getProfileStory,
  checkProfileSlug,
  createProfile,

  // Stories
  getStory,
  getStoriesByKinds,

  // Site
  getSpotlight,

  // Auth
  handleAuthCallback,

  // Sessions
  getCurrentSession,
};
