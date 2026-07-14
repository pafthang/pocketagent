package schedule

import (
	"context"
	"log/slog"
	"time"

	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/models"
	"github.com/pafthang/pocketagent/pkgs/service"
)

// EnqueueFunc persists and publishes a task to the orchestrator.
type EnqueueFunc func(ctx context.Context, w *service.Worker, pb *pbclient.Client, task models.Task) (string, error)

// Scheduler triggers cron schedules and enqueues tasks.
type Scheduler struct {
	Worker   *service.Worker
	Pocket   *pbclient.Client
	Log      *slog.Logger
	Interval time.Duration
	Enqueue  EnqueueFunc
}

func (s *Scheduler) Run(ctx context.Context) {
	if s == nil || s.Worker == nil || s.Pocket == nil || s.Enqueue == nil {
		return
	}
	if s.Interval <= 0 {
		s.Interval = time.Minute
	}

	ticker := time.NewTicker(s.Interval)
	defer ticker.Stop()

	s.tick(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.tick(ctx)
		}
	}
}

func (s *Scheduler) tick(ctx context.Context) {
	now := time.Now().UTC()
	schedules, err := s.Pocket.ListDueSchedules(now)
	if err != nil {
		if s.Log != nil {
			s.Log.Warn("failed to list due schedules", "error", err)
		}
		return
	}

	for _, schedule := range schedules {
		if err := s.runSchedule(ctx, schedule, now); err != nil && s.Log != nil {
			s.Log.Warn("schedule run failed",
				"schedule_id", schedule.ID,
				"name", schedule.Name,
				"error", err,
			)
		}
	}
}

func (s *Scheduler) runSchedule(ctx context.Context, schedule models.Schedule, now time.Time) error {
	task := models.Task{
		SpaceID:        schedule.SpaceID,
		AgentID:        schedule.AgentID,
		Prompt:         schedule.Prompt,
		Workflow:       schedule.Workflow,
		WorkerAgentIDs: schedule.WorkerAgentIDs,
	}

	corrID, err := s.Enqueue(ctx, s.Worker, s.Pocket, task)
	if err != nil {
		return err
	}

	nextRun, err := NextCronRun(schedule.CronExpr, now)
	if err != nil {
		return err
	}

	_, err = s.Pocket.UpdateSchedule(schedule.ID, models.Schedule{
		LastRunAt:  now.Format(time.RFC3339),
		NextRunAt:  nextRun.UTC().Format(time.RFC3339),
		LastTaskID: corrID,
	})
	return err
}
