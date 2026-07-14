package port

import (
	"testing"

	"github.com/pafthang/pocketagent/internal/ctrl/catalog"
)

func TestCollectManagedPorts(t *testing.T) {
	ports := CollectManaged([]catalog.Service{
		{Name: "gate", WaitPort: 8080},
		{Name: "exec", HealthPort: 9084},
		{Name: "task", HealthPort: 9085},
		{Name: "memo", WaitPort: 8082, HealthPort: 8082},
	})
	want := map[int]bool{8080: true, 9084: true, 9085: true, 8082: true}
	if len(ports) != len(want) {
		t.Fatalf("CollectManaged() = %v, want %d unique ports", ports, len(want))
	}
	for _, port := range ports {
		if !want[port] {
			t.Fatalf("unexpected port %d in %v", port, ports)
		}
		delete(want, port)
	}
}

func TestListenersOnPort8090(t *testing.T) {
	pids, err := listenersOnPort(8090)
	if err != nil {
		t.Fatalf("listenersOnPort: %v", err)
	}
	if len(pids) == 0 {
		t.Skip("no listener on 8090")
	}
	t.Logf("listeners on 8090: %v", pids)
}