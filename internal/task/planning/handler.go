package planning

import (
	"context"

	"github.com/pafthang/pocketagent/internal/exec/tools"
	natsclient "github.com/pafthang/pocketagent/internal/nats/client"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/internal/space/projects"
	"github.com/pafthang/pocketagent/pkgs/models"
	"github.com/pafthang/pocketagent/pkgs/ollama"
	"github.com/pafthang/pocketagent/pkgs/service"
)

// Handler returns a NATS consumer for async project planning.
func Handler(w *service.Worker, pb *pbclient.Client, oc *ollama.Client, model string) func(ctx context.Context, cmd models.ProjectPlanCommand) error {
	return func(ctx context.Context, cmd models.ProjectPlanCommand) error {
		planner := &projects.Planner{
			PB:     pb,
			Ollama: oc,
			Model:  model,
			Search: tools.LoadFromEnv(),
			Log:    w.Log,
			Publish: func(ctx context.Context, projectID string, event map[string]any) error {
				return w.PublishJSON(ctx, natsclient.EventSubject(projectID), event)
			},
		}
		return planner.Run(ctx, cmd)
	}
}
