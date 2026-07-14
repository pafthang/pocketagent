package bootstrap

import (
	"fmt"
	"strings"

	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/models"
	"github.com/pocketbase/pocketbase/core"
)

func ensureBootstrapSuperuser(app core.App, email, password string) error {
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)
	if email == "" || password == "" {
		app.Logger().Info("bootstrap superuser skipped: set POCKETBASE_SUPERUSER_EMAIL and POCKETBASE_SUPERUSER_PASSWORD")
		return nil
	}

	adminSpace, err := getOrCreateAdminSpace(app)
	if err != nil {
		return fmt.Errorf("admin space: %w", err)
	}

	user, err := upsertAuthUser(app, pbclient.UsersCollection, email, password)
	if err != nil {
		return fmt.Errorf("bootstrap user: %w", err)
	}

	if _, err := upsertAuthUser(app, core.CollectionNameSuperusers, email, password); err != nil {
		return fmt.Errorf("bootstrap superuser: %w", err)
	}

	if err := ensureAdminSpaceMember(app, adminSpace.Id, user.Id); err != nil {
		return fmt.Errorf("admin space membership: %w", err)
	}

	app.Logger().Info(
		"bootstrap superuser ready",
		"email", email,
		"space", models.SystemSpaceSlug,
		"role", models.RoleAdmin,
	)
	return nil
}

func upsertAuthUser(app core.App, collectionName, email, password string) (*core.Record, error) {
	col, err := app.FindCollectionByNameOrId(collectionName)
	if err != nil {
		return nil, err
	}

	record, err := app.FindAuthRecordByEmail(col, email)
	if err != nil {
		record = core.NewRecord(col)
	}

	record.SetEmail(email)
	record.SetPassword(password)

	if err := app.Save(record); err != nil {
		return nil, err
	}
	return record, nil
}

func ensureAdminSpaceMember(app core.App, spaceID, userID string) error {
	col, err := app.FindCollectionByNameOrId(pbclient.SpaceMembersCollection)
	if err != nil {
		return err
	}

	records, err := app.FindRecordsByFilter(col.Id,
		"space_id = {:space_id} && user_id = {:user_id}",
		"", 1, 0,
		map[string]any{"space_id": spaceID, "user_id": userID},
	)
	if err != nil {
		return err
	}
	if len(records) > 0 {
		existing := records[0]
		if existing.GetString("role") != models.RoleAdmin {
			existing.Set("role", models.RoleAdmin)
			return app.Save(existing)
		}
		return nil
	}

	member := core.NewRecord(col)
	member.Set("space_id", spaceID)
	member.Set("user_id", userID)
	member.Set("role", models.RoleAdmin)
	return app.Save(member)
}

func getOrCreateAdminSpace(app core.App) (*core.Record, error) {
	col, err := app.FindCollectionByNameOrId(pbclient.SpacesCollection)
	if err != nil {
		return nil, err
	}

	records, err := app.FindRecordsByFilter(col.Id, "slug = {:slug}", "", 1, 0, map[string]any{
		"slug": models.SystemSpaceSlug,
	})
	if err != nil {
		return nil, err
	}
	if len(records) > 0 {
		return records[0], nil
	}

	record := core.NewRecord(col)
	record.Set("name", "Admin")
	record.Set("slug", models.SystemSpaceSlug)
	record.Set("description", "System space for super-admins")
	record.Set("is_system", true)
	if err := app.Save(record); err != nil {
		return nil, err
	}
	return record, nil
}
