package bootstrap

import (
	"github.com/pocketbase/pocketbase/core"
)

// Config carries seed credentials for schema bootstrap.
type Config struct {
	SuperuserEmail    string
	SuperuserPassword string
}

// Step runs one bootstrap phase against a PocketBase app.
type Step func(app core.App) error

// RegisterAll returns ordered bootstrap steps for collection/schema setup.
func RegisterAll(cfg Config) []Step {
	return []Step{
		ensureAgentsCollection,
		func(app core.App) error {
			return ensureSpaceBootstrap(app, cfg.SuperuserEmail, cfg.SuperuserPassword)
		},
		ensureTasksCollection,
		ensureSchedulesCollection,
		ensureMCPServersCollection,
		ensureSkillsCollection,
		ensureSpaceProfilesCollection,
		ensureTaskEventsCollection,
		ensureProjectsCollection,
		ensureProjectItemsCollection,
		ensureFilesCollection,
	}
}

// Run executes all registered bootstrap steps.
func Run(app core.App, cfg Config) error {
	for _, step := range RegisterAll(cfg) {
		if err := step(app); err != nil {
			return err
		}
	}
	return nil
}