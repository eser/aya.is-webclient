package storage

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
	"github.com/eser/aya.is/services/pkg/lib/vars"
)

func (r *Repository) GetProfileIDBySlug(ctx context.Context, slug string) (string, error) {
	var result string

	err := r.cache.Execute(
		ctx,
		"profile_id_by_slug:"+slug,
		&result,
		func(ctx context.Context) (any, error) {
			row, err := r.queries.GetProfileIDBySlug(ctx, GetProfileIDBySlugParams{Slug: slug})
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, nil //nolint:nilnil
				}

				return nil, err
			}

			return row, nil
		},
	)

	return result, err //nolint:wrapcheck
}

func (r *Repository) CheckProfileSlugExists(ctx context.Context, slug string) (bool, error) {
	var result bool

	err := r.cache.Execute(
		ctx,
		"profile_slug_exists:"+slug,
		&result,
		func(ctx context.Context) (any, error) {
			exists, err := r.queries.CheckProfileSlugExists(
				ctx,
				CheckProfileSlugExistsParams{Slug: slug},
			)
			if err != nil {
				return nil, err
			}

			return exists, nil
		},
	)

	return result, err //nolint:wrapcheck
}

func (r *Repository) GetProfileIDByCustomDomain(
	ctx context.Context,
	domain string,
) (*string, error) {
	var result *string

	err := r.cache.Execute(
		ctx,
		"profile_id_by_custom_domain:"+domain,
		&result,
		func(ctx context.Context) (any, error) {
			row, err := r.queries.GetProfileIDByCustomDomain(
				ctx,
				GetProfileIDByCustomDomainParams{
					CustomDomain: sql.NullString{String: domain, Valid: true},
				},
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, nil //nolint:nilnil
				}

				return nil, err
			}

			return &row, nil
		},
	)

	return result, err //nolint:wrapcheck
}

func (r *Repository) GetProfileByID(
	ctx context.Context,
	localeCode string,
	id string,
) (*profiles.Profile, error) {
	row, err := r.queries.GetProfileByID(ctx, GetProfileByIDParams{LocaleCode: localeCode, ID: id})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	result := &profiles.Profile{
		ID:                row.Profile.ID,
		Slug:              row.Profile.Slug,
		Kind:              row.Profile.Kind,
		CustomDomain:      vars.ToStringPtr(row.Profile.CustomDomain),
		ProfilePictureURI: vars.ToStringPtr(row.Profile.ProfilePictureURI),
		Pronouns:          vars.ToStringPtr(row.Profile.Pronouns),
		Title:             row.ProfileTx.Title,
		Description:       row.ProfileTx.Description,
		Properties:        vars.ToObject(row.Profile.Properties),
		CreatedAt:         row.Profile.CreatedAt,
		UpdatedAt:         vars.ToTimePtr(row.Profile.UpdatedAt),
		DeletedAt:         vars.ToTimePtr(row.Profile.DeletedAt),
	}

	return result, nil
}

func (r *Repository) ListProfiles(
	ctx context.Context,
	localeCode string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*profiles.Profile], error) {
	var wrappedResponse cursors.Cursored[[]*profiles.Profile]

	rows, err := r.queries.ListProfiles(
		ctx,
		ListProfilesParams{
			LocaleCode: localeCode,
			FilterKind: vars.MapValueToNullString(cursor.Filters, "kind"),
		},
	)
	if err != nil {
		return wrappedResponse, err
	}

	result := make([]*profiles.Profile, len(rows))
	for i, row := range rows {
		result[i] = &profiles.Profile{
			ID:                row.Profile.ID,
			Slug:              row.Profile.Slug,
			Kind:              row.Profile.Kind,
			CustomDomain:      vars.ToStringPtr(row.Profile.CustomDomain),
			ProfilePictureURI: vars.ToStringPtr(row.Profile.ProfilePictureURI),
			Pronouns:          vars.ToStringPtr(row.Profile.Pronouns),
			Title:             row.ProfileTx.Title,
			Description:       row.ProfileTx.Description,
			Properties:        vars.ToObject(row.Profile.Properties),
			CreatedAt:         row.Profile.CreatedAt,
			UpdatedAt:         vars.ToTimePtr(row.Profile.UpdatedAt),
			DeletedAt:         vars.ToTimePtr(row.Profile.DeletedAt),
		}
	}

	wrappedResponse.Data = result

	if len(result) == cursor.Limit && len(result) > 0 {
		wrappedResponse.CursorPtr = &result[len(result)-1].ID
	}

	return wrappedResponse, nil
}

func (r *Repository) ListProfilePagesByProfileID(
	ctx context.Context,
	localeCode string,
	profileID string,
) ([]*profiles.ProfilePageBrief, error) {
	rows, err := r.queries.ListProfilePagesByProfileID(
		ctx,
		ListProfilePagesByProfileIDParams{LocaleCode: localeCode, ProfileID: profileID},
	)
	if err != nil {
		return nil, err
	}

	profilePages := make([]*profiles.ProfilePageBrief, len(rows))
	for i, row := range rows {
		profilePages[i] = &profiles.ProfilePageBrief{
			ID:              row.ID,
			Slug:            row.Slug,
			CoverPictureURI: vars.ToStringPtr(row.CoverPictureURI),
			Title:           row.Title,
			Summary:         row.Summary,
		}
	}

	return profilePages, nil
}

func (r *Repository) GetProfilePageByProfileIDAndSlug(
	ctx context.Context,
	localeCode string,
	profileID string,
	pageSlug string,
) (*profiles.ProfilePage, error) {
	row, err := r.queries.GetProfilePageByProfileIDAndSlug(
		ctx,
		GetProfilePageByProfileIDAndSlugParams{
			LocaleCode: localeCode,
			ProfileID:  profileID,
			PageSlug:   pageSlug,
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	result := &profiles.ProfilePage{
		ID:              row.ID,
		Slug:            row.Slug,
		CoverPictureURI: vars.ToStringPtr(row.CoverPictureURI),
		Title:           row.Title,
		Summary:         row.Summary,
		Content:         row.Content,
		PublishedAt:     vars.ToTimePtr(row.PublishedAt),
	}

	return result, nil
}

func (r *Repository) ListProfileLinksByProfileID(
	ctx context.Context,
	_localeCode string,
	profileID string,
) ([]*profiles.ProfileLinkBrief, error) {
	rows, err := r.queries.ListProfileLinksByProfileID(
		ctx,
		ListProfileLinksByProfileIDParams{ProfileID: profileID},
	)
	if err != nil {
		return nil, err
	}

	profileLinks := make([]*profiles.ProfileLinkBrief, len(rows))
	for i, row := range rows {
		profileLinks[i] = &profiles.ProfileLinkBrief{
			ID:         row.ID,
			Kind:       row.Kind,
			IsVerified: row.IsVerified,
			PublicID:   row.PublicID.String,
			URI:        row.URI.String,
			Title:      row.Title,
		}
	}

	return profileLinks, nil
}

func (r *Repository) ListProfileContributions( //nolint:funlen
	ctx context.Context,
	localeCode string,
	profileID string,
	kinds []string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*profiles.ProfileMembership], error) {
	var wrappedResponse cursors.Cursored[[]*profiles.ProfileMembership]

	rows, err := r.queries.ListProfileMemberships(
		ctx,
		ListProfileMembershipsParams{
			LocaleCode:              localeCode,
			FilterProfileID:         sql.NullString{String: "", Valid: false},
			FilterProfileKind:       sql.NullString{String: strings.Join(kinds, ","), Valid: true},
			FilterMemberProfileID:   sql.NullString{String: profileID, Valid: true},
			FilterMemberProfileKind: sql.NullString{String: "", Valid: false},
		},
	)
	if err != nil {
		return wrappedResponse, err
	}

	profileMemberships := make([]*profiles.ProfileMembership, len(rows))
	for i, row := range rows { //nolint:dupl
		profileMemberships[i] = &profiles.ProfileMembership{
			ID:         row.ProfileMembership.ID,
			Kind:       row.ProfileMembership.Kind,
			StartedAt:  vars.ToTimePtr(row.ProfileMembership.StartedAt),
			FinishedAt: vars.ToTimePtr(row.ProfileMembership.FinishedAt),
			Properties: vars.ToObject(row.ProfileMembership.Properties),
			Profile: &profiles.Profile{
				ID:                row.Profile.ID,
				Slug:              row.Profile.Slug,
				Kind:              row.Profile.Kind,
				CustomDomain:      vars.ToStringPtr(row.Profile.CustomDomain),
				ProfilePictureURI: vars.ToStringPtr(row.Profile.ProfilePictureURI),
				Pronouns:          vars.ToStringPtr(row.Profile.Pronouns),
				Title:             row.ProfileTx.Title,
				Description:       row.ProfileTx.Description,
				Properties:        vars.ToObject(row.Profile.Properties),
				CreatedAt:         row.Profile.CreatedAt,
				UpdatedAt:         vars.ToTimePtr(row.Profile.UpdatedAt),
				DeletedAt:         vars.ToTimePtr(row.Profile.DeletedAt),
			},
			MemberProfile: &profiles.Profile{
				ID:                row.Profile_2.ID,
				Slug:              row.Profile_2.Slug,
				Kind:              row.Profile_2.Kind,
				CustomDomain:      vars.ToStringPtr(row.Profile_2.CustomDomain),
				ProfilePictureURI: vars.ToStringPtr(row.Profile_2.ProfilePictureURI),
				Pronouns:          vars.ToStringPtr(row.Profile_2.Pronouns),
				Title:             row.ProfileTx_2.Title,
				Description:       row.ProfileTx_2.Description,
				Properties:        vars.ToObject(row.Profile_2.Properties),
				CreatedAt:         row.Profile_2.CreatedAt,
				UpdatedAt:         vars.ToTimePtr(row.Profile_2.UpdatedAt),
				DeletedAt:         vars.ToTimePtr(row.Profile_2.DeletedAt),
			},
		}
	}

	wrappedResponse.Data = profileMemberships

	if len(profileMemberships) == cursor.Limit {
		wrappedResponse.CursorPtr = &profileMemberships[len(profileMemberships)-1].ID
	}

	return wrappedResponse, nil
}

//nolint:funlen,dupl
func (r *Repository) ListProfileMembers(
	ctx context.Context,
	localeCode string,
	profileID string,
	kinds []string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*profiles.ProfileMembership], error) {
	var wrappedResponse cursors.Cursored[[]*profiles.ProfileMembership]

	rows, err := r.queries.ListProfileMemberships(
		ctx,
		ListProfileMembershipsParams{
			LocaleCode:              localeCode,
			FilterProfileID:         sql.NullString{String: profileID, Valid: true},
			FilterProfileKind:       sql.NullString{String: "", Valid: false},
			FilterMemberProfileID:   sql.NullString{String: "", Valid: false},
			FilterMemberProfileKind: sql.NullString{String: strings.Join(kinds, ","), Valid: true},
		},
	)
	if err != nil {
		return wrappedResponse, err
	}

	profileMemberships := make([]*profiles.ProfileMembership, len(rows))
	for i, row := range rows {
		profileMemberships[i] = &profiles.ProfileMembership{
			ID:         row.ProfileMembership.ID,
			Kind:       row.ProfileMembership.Kind,
			StartedAt:  vars.ToTimePtr(row.ProfileMembership.StartedAt),
			FinishedAt: vars.ToTimePtr(row.ProfileMembership.FinishedAt),
			Properties: vars.ToObject(row.ProfileMembership.Properties),
			Profile: &profiles.Profile{
				ID:                row.Profile.ID,
				Slug:              row.Profile.Slug,
				Kind:              row.Profile.Kind,
				CustomDomain:      vars.ToStringPtr(row.Profile.CustomDomain),
				ProfilePictureURI: vars.ToStringPtr(row.Profile.ProfilePictureURI),
				Pronouns:          vars.ToStringPtr(row.Profile.Pronouns),
				Title:             row.ProfileTx.Title,
				Description:       row.ProfileTx.Description,
				Properties:        vars.ToObject(row.Profile.Properties),
				CreatedAt:         row.Profile.CreatedAt,
				UpdatedAt:         vars.ToTimePtr(row.Profile.UpdatedAt),
				DeletedAt:         vars.ToTimePtr(row.Profile.DeletedAt),
			},
			MemberProfile: &profiles.Profile{
				ID:                row.Profile_2.ID,
				Slug:              row.Profile_2.Slug,
				Kind:              row.Profile_2.Kind,
				CustomDomain:      vars.ToStringPtr(row.Profile_2.CustomDomain),
				ProfilePictureURI: vars.ToStringPtr(row.Profile_2.ProfilePictureURI),
				Pronouns:          vars.ToStringPtr(row.Profile_2.Pronouns),
				Title:             row.ProfileTx_2.Title,
				Description:       row.ProfileTx_2.Description,
				Properties:        vars.ToObject(row.Profile_2.Properties),
				CreatedAt:         row.Profile_2.CreatedAt,
				UpdatedAt:         vars.ToTimePtr(row.Profile_2.UpdatedAt),
				DeletedAt:         vars.ToTimePtr(row.Profile_2.DeletedAt),
			},
		}
	}

	wrappedResponse.Data = profileMemberships

	if len(profileMemberships) == cursor.Limit {
		wrappedResponse.CursorPtr = &profileMemberships[len(profileMemberships)-1].ID
	}

	return wrappedResponse, nil
}

func (r *Repository) GetProfileMembershipsByMemberProfileID(
	ctx context.Context,
	localeCode string,
	memberProfileID string,
) ([]*profiles.ProfileMembership, error) {
	rows, err := r.queries.GetProfileMembershipsByMemberProfileID(
		ctx,
		GetProfileMembershipsByMemberProfileIDParams{
			LocaleCode:      localeCode,
			MemberProfileID: memberProfileID,
		},
	)
	if err != nil {
		return nil, err
	}

	memberships := make([]*profiles.ProfileMembership, len(rows))
	for i, row := range rows {
		memberships[i] = &profiles.ProfileMembership{
			ID:         row.MembershipID,
			Kind:       row.MembershipKind,
			StartedAt:  vars.ToTimePtr(row.StartedAt),
			FinishedAt: vars.ToTimePtr(row.FinishedAt),
			Properties: vars.ToObject(row.MembershipProperties),
			Profile: &profiles.Profile{
				ID:                row.Profile.ID,
				Slug:              row.Profile.Slug,
				Kind:              row.Profile.Kind,
				CustomDomain:      vars.ToStringPtr(row.Profile.CustomDomain),
				ProfilePictureURI: vars.ToStringPtr(row.Profile.ProfilePictureURI),
				Pronouns:          vars.ToStringPtr(row.Profile.Pronouns),
				Title:             row.ProfileTx.Title,
				Description:       row.ProfileTx.Description,
				Properties:        vars.ToObject(row.Profile.Properties),
				CreatedAt:         row.Profile.CreatedAt,
				UpdatedAt:         vars.ToTimePtr(row.Profile.UpdatedAt),
				DeletedAt:         vars.ToTimePtr(row.Profile.DeletedAt),
			},
			// MemberProfile is not needed for this use case since we're filtering by member profile ID
			MemberProfile: nil,
		}
	}

	return memberships, nil
}

func (r *Repository) CreateProfile(
	ctx context.Context,
	id string,
	slug string,
	kind string,
	customDomain *string,
	profilePictureURI *string,
	pronouns *string,
	properties map[string]any,
) error {
	params := CreateProfileParams{
		ID:                id,
		Slug:              slug,
		Kind:              kind,
		CustomDomain:      vars.ToSQLNullString(customDomain),
		ProfilePictureURI: vars.ToSQLNullString(profilePictureURI),
		Pronouns:          vars.ToSQLNullString(pronouns),
		Properties:        vars.ToSQLNullRawMessage(properties),
	}

	return r.queries.CreateProfile(ctx, params)
}

func (r *Repository) CreateProfileTx(
	ctx context.Context,
	profileID string,
	localeCode string,
	title string,
	description string,
	properties map[string]any,
) error {
	params := CreateProfileTxParams{
		ProfileID:   profileID,
		LocaleCode:  localeCode,
		Title:       title,
		Description: description,
		Properties:  vars.ToSQLNullRawMessage(properties),
	}

	return r.queries.CreateProfileTx(ctx, params)
}
