package bootstrap

import pbclient "github.com/pafthang/pocketagent/internal/pocket/client"

// ExpectedCollections lists PocketBase collections created by Run.
func ExpectedCollections() []string {
	return []string{
		pbclient.AgentsCollection,
		pbclient.UsersCollection,
		pbclient.SpacesCollection,
		pbclient.SpaceMembersCollection,
		pbclient.TeamsCollection,
		pbclient.TeamMembersCollection,
		pbclient.SpaceInvitesCollection,
		pbclient.AuditLogsCollection,
		pbclient.EmailVerificationsCollection,
		pbclient.TasksCollection,
		pbclient.SchedulesCollection,
		pbclient.MCPServersCollection,
		pbclient.SkillsCollection,
		pbclient.SpaceProfilesCollection,
		pbclient.TaskEventsCollection,
		pbclient.ProjectsCollection,
		pbclient.ProjectItemsCollection,
		pbclient.FilesCollection,
	}
}