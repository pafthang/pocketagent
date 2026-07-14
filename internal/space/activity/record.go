package activity

import (
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/models"
)

// Recorder returns a best-effort async task event persister.
func Recorder(pb *pbclient.Client) func(models.TaskEvent) {
	return func(event models.TaskEvent) {
		Record(pb, event)
	}
}

// Record persists a task event when space_id is set.
func Record(pb *pbclient.Client, event models.TaskEvent) {
	if pb == nil || event.SpaceID == "" || event.TaskID == "" {
		return
	}
	go func() {
		_, _ = pb.CreateTaskEventRecord(event)
	}()
}
