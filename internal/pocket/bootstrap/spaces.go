package bootstrap

import (
	"database/sql"
	"errors"

	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/models"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

func ensureSpaceCollections(app core.App) error {
	if err := ensureUsersCollection(app); err != nil {
		return err
	}
	if err := ensureCollectionSpaces(app); err != nil {
		return err
	}
	if err := ensureCollectionSpaceMembers(app); err != nil {
		return err
	}
	if err := ensureCollectionTeams(app); err != nil {
		return err
	}
	if err := ensureCollectionTeamMembers(app); err != nil {
		return err
	}
	if err := ensureInviteAuditCollections(app); err != nil {
		return err
	}
	if _, err := getOrCreateAdminSpace(app); err != nil {
		return err
	}
	return nil
}

func ensureSpaceBootstrap(app core.App, superuserEmail, superuserPassword string) error {
	if err := ensureSpaceCollections(app); err != nil {
		return err
	}
	return ensureBootstrapSuperuser(app, superuserEmail, superuserPassword)
}

func ensureUsersCollection(app core.App) error {
	col, err := app.FindCollectionByNameOrId(pbclient.UsersCollection)
	if errors.Is(err, sql.ErrNoRows) {
		col = core.NewAuthCollection(pbclient.UsersCollection)
		applyUsersRules(col)
		return app.Save(col)
	}
	if err != nil {
		return err
	}

	before := usersRulesSnapshot(col)
	applyUsersRules(col)
	if usersRulesSnapshot(col) == before {
		return nil
	}
	return app.Save(col)
}

func applyUsersRules(col *core.Collection) {
	// Registration only via backend (superuser); users manage own record.
	col.ListRule = nil
	col.ViewRule = types.Pointer("@request.auth.id = id")
	col.CreateRule = nil
	col.UpdateRule = types.Pointer("@request.auth.id = id")
	col.DeleteRule = types.Pointer("@request.auth.id = id")
}

type usersRules struct {
	list, view, create, update, delete string
}

func usersRulesSnapshot(col *core.Collection) usersRules {
	return usersRules{
		list:   ruleString(col.ListRule),
		view:   ruleString(col.ViewRule),
		create: ruleString(col.CreateRule),
		update: ruleString(col.UpdateRule),
		delete: ruleString(col.DeleteRule),
	}
}

func ruleString(rule *string) string {
	if rule == nil {
		return "<nil>"
	}
	return *rule
}

func ensureCollectionSpaces(app core.App) error {
	return ensureLockedCollection(app, pbclient.SpacesCollection, func(col *core.Collection) {
		addFieldIfMissing(col, &core.TextField{Name: "name", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "slug", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "description"})
		addFieldIfMissing(col, &core.BoolField{Name: "is_system"})
	})
}

func ensureCollectionSpaceMembers(app core.App) error {
	return ensureLockedCollection(app, pbclient.SpaceMembersCollection, func(col *core.Collection) {
		addFieldIfMissing(col, &core.TextField{Name: "space_id", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "user_id", Required: true})
		addFieldIfMissing(col, &core.SelectField{
			Name:     "role",
			Required: true,
			Values:   []string{models.RoleAdmin, models.RoleEditor, models.RoleViewer},
		})
	})
}

func ensureCollectionTeams(app core.App) error {
	return ensureLockedCollection(app, pbclient.TeamsCollection, func(col *core.Collection) {
		addFieldIfMissing(col, &core.TextField{Name: "space_id", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "name", Required: true})
		addFieldIfMissing(col, &core.TextField{Name: "description"})
	})
}

func ensureCollectionTeamMembers(app core.App) error {
	return ensureLockedCollection(app, pbclient.TeamMembersCollection, func(col *core.Collection) {
		addFieldIfMissing(col, &core.TextField{Name: "team_id", Required: true})
		addFieldIfMissing(col, &core.SelectField{
			Name:     "member_type",
			Required: true,
			Values:   []string{models.MemberTypeUser, models.MemberTypeAgent},
		})
		addFieldIfMissing(col, &core.TextField{Name: "member_id", Required: true})
	})
}
