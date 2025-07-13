-- name: GetProfileIDBySlug :one
SELECT id
FROM "profile"
WHERE slug = sqlc.arg(slug)
  AND deleted_at IS NULL
LIMIT 1;

-- name: GetProfileIDByCustomDomain :one
SELECT id
FROM "profile"
WHERE custom_domain = sqlc.arg(custom_domain)
  AND deleted_at IS NULL
LIMIT 1;

-- name: GetProfileByID :one
SELECT sqlc.embed(p), sqlc.embed(pt)
FROM "profile" p
  INNER JOIN "profile_tx" pt ON pt.profile_id = p.id
  AND pt.locale_code = sqlc.arg(locale_code)
WHERE p.id = sqlc.arg(id)
  AND p.deleted_at IS NULL
LIMIT 1;

-- name: ListProfiles :many
SELECT sqlc.embed(p), sqlc.embed(pt)
FROM "profile" p
  INNER JOIN "profile_tx" pt ON pt.profile_id = p.id
  AND pt.locale_code = sqlc.arg(locale_code)
WHERE (sqlc.narg(filter_kind)::TEXT IS NULL OR p.kind = ANY(string_to_array(sqlc.narg(filter_kind)::TEXT, ',')))
  AND p.deleted_at IS NULL;

-- name: CreateProfile :exec
INSERT INTO "profile" (id, slug)
VALUES (sqlc.arg(id), sqlc.arg(slug));

-- name: UpdateProfile :execrows
UPDATE "profile"
SET slug = sqlc.arg(slug)
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: RemoveProfile :execrows
UPDATE "profile"
SET deleted_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: ListProfileLinksForKind :many
SELECT pl.*
FROM "profile_link" pl
  INNER JOIN "profile" p ON p.id = pl.profile_id
  AND p.deleted_at IS NULL
WHERE pl.kind = sqlc.arg(kind)
  AND pl.deleted_at IS NULL
ORDER BY pl."order";

-- name: ListProfilePagesByProfileID :many
SELECT pp.*, ppt.*
FROM "profile_page" pp
  INNER JOIN "profile_page_tx" ppt ON ppt.profile_page_id = pp.id
  AND ppt.locale_code = sqlc.arg(locale_code)
WHERE pp.profile_id = sqlc.arg(profile_id)
  AND pp.deleted_at IS NULL
ORDER BY pp."order";

-- name: GetProfilePageByProfileIDAndSlug :one
SELECT pp.*, ppt.*
FROM "profile_page" pp
  INNER JOIN "profile_page_tx" ppt ON ppt.profile_page_id = pp.id
  AND ppt.locale_code = sqlc.arg(locale_code)
WHERE pp.profile_id = sqlc.arg(profile_id) AND pp.slug = sqlc.arg(page_slug) AND pp.deleted_at IS NULL
ORDER BY pp."order";

-- name: ListProfileLinksByProfileID :many
SELECT *
FROM "profile_link"
WHERE profile_id = sqlc.arg(profile_id)
  AND is_hidden = FALSE
  AND deleted_at IS NULL
ORDER BY "order";

-- name: ListProfileMemberships :many
SELECT
  sqlc.embed(pm),
  sqlc.embed(p1),
  sqlc.embed(p1t),
  sqlc.embed(p2),
  sqlc.embed(p2t)
FROM
	"profile_membership" pm
  INNER JOIN "profile" p1 ON p1.id = pm.profile_id
    AND (sqlc.narg(filter_profile_kind)::TEXT IS NULL OR p1.kind = ANY(string_to_array(sqlc.narg(filter_profile_kind)::TEXT, ',')))
    AND p1.deleted_at IS NULL
  INNER JOIN "profile_tx" p1t ON p1t.profile_id = p1.id
	  AND p1t.locale_code = sqlc.arg(locale_code)
  INNER JOIN "profile" p2 ON p2.id = pm.member_profile_id
    AND (sqlc.narg(filter_member_profile_kind)::TEXT IS NULL OR p2.kind = ANY(string_to_array(sqlc.narg(filter_member_profile_kind)::TEXT, ',')))
    AND p2.deleted_at IS NULL
  INNER JOIN "profile_tx" p2t ON p2t.profile_id = p2.id
	  AND p2t.locale_code = sqlc.arg(locale_code)
WHERE pm.deleted_at IS NULL
    AND (sqlc.narg(filter_profile_id)::TEXT IS NULL OR pm.profile_id = sqlc.narg(filter_profile_id)::TEXT)
    AND (sqlc.narg(filter_member_profile_id)::TEXT IS NULL OR pm.member_profile_id = sqlc.narg(filter_member_profile_id)::TEXT);

-- name: GetProfileMembershipsByMemberProfileID :many
SELECT
  pm.id as membership_id,
  pm.kind as membership_kind,
  pm.started_at,
  pm.finished_at,
  pm.properties as membership_properties,
  pm.created_at as membership_created_at,
  sqlc.embed(p),
  sqlc.embed(pt)
FROM
  "profile_membership" pm
  INNER JOIN "profile" p ON p.id = pm.profile_id
    AND p.deleted_at IS NULL
  INNER JOIN "profile_tx" pt ON pt.profile_id = p.id
    AND pt.locale_code = sqlc.arg(locale_code)
WHERE
  pm.deleted_at IS NULL
  AND pm.member_profile_id = sqlc.arg(member_profile_id)
  AND (pm.finished_at IS NULL OR pm.finished_at > NOW())
ORDER BY pm.created_at DESC;
