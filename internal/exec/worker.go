package exec

import (
	natsclient "github.com/pafthang/pocketagent/internal/nats/client"
	"github.com/pafthang/pocketagent/internal/exec/handler"
	"github.com/pafthang/pocketagent/internal/space/activity"
	"github.com/pafthang/pocketagent/pkgs/service"
)

func wireWorker(w *service.Worker, deps *WorkerDeps) error {
	w.EventRecorder = activity.Recorder(deps.Pocket)
	_, err := service.SubscribeJSON(w.Consumer, natsclient.SubjectTasks, handler.Task(w, handler.Deps{
		Executor: deps.Executor,
		Pocket:   deps.Pocket,
		ToolCfg:  deps.ToolCfg,
	}))
	return err
}