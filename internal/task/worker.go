package task

import (
	"context"

	natsclient "github.com/pafthang/pocketagent/internal/nats/client"
	"github.com/pafthang/pocketagent/internal/space/activity"
	"github.com/pafthang/pocketagent/internal/task/orchestrator"
	"github.com/pafthang/pocketagent/internal/task/planning"
	"github.com/pafthang/pocketagent/internal/task/schedule"
	"github.com/pafthang/pocketagent/pkgs/service"
)

func registerConsumers(w *service.Worker, d *WorkerDeps) error {
	if _, err := service.SubscribeJSON(w.Consumer, natsclient.SubjectOrchestrator, orchestrator.Handler(orchestrator.Deps{
		Worker:     w,
		Store:      d.Store,
		Decomposer: d.Decomposer,
		Memory:     d.Memory,
		Embedder:   d.Ollama,
		TimeoutSec: d.Config.TimeoutSec,
	})); err != nil {
		return err
	}
	if _, err := service.SubscribeJSON(w.Consumer, natsclient.SubjectProjectsPlan, planning.Handler(
		w, d.Pocket, d.Ollama, d.Config.LLMModel,
	)); err != nil {
		return err
	}
	return nil
}

func startScheduler(w *service.Worker, d *WorkerDeps) {
	scheduler := &schedule.Scheduler{
		Worker:   w,
		Pocket:   d.Pocket,
		Log:      w.Log,
		Interval: d.Config.SchedulerInterval(),
		Enqueue:  EnqueueTask,
	}
	go scheduler.Run(context.Background())
}

func wireWorker(w *service.Worker, cfg *Config) (*WorkerDeps, error) {
	deps, err := buildDeps(w, cfg)
	if err != nil {
		return nil, err
	}
	w.EventRecorder = activity.Recorder(deps.Pocket)
	if err := registerConsumers(w, deps); err != nil {
		return nil, err
	}
	startScheduler(w, deps)
	return deps, nil
}
