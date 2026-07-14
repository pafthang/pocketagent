package port

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/pafthang/pocketagent/internal/ctrl/catalog"
)

const freeTimeout = 5 * time.Second

// CollectManaged returns unique TCP ports used for readiness checks.
func CollectManaged(services []catalog.Service) []int {
	seen := make(map[int]struct{})
	ports := make([]int, 0)
	for _, svc := range services {
		for _, port := range ServicePorts(&svc) {
			if _, ok := seen[port]; ok {
				continue
			}
			seen[port] = struct{}{}
			ports = append(ports, port)
		}
	}
	return ports
}

// ServicePorts returns wait and health ports for a service definition.
func ServicePorts(def *catalog.Service) []int {
	if def == nil {
		return nil
	}
	ports := make([]int, 0, 2)
	if def.WaitPort > 0 {
		ports = append(ports, def.WaitPort)
	}
	if def.HealthPort > 0 && def.HealthPort != def.WaitPort {
		ports = append(ports, def.HealthPort)
	}
	return ports
}

// EnsureFree terminates stale listeners and waits until the port is available.
func EnsureFree(port int) error {
	if !isListening("127.0.0.1", port) && !isListening("0.0.0.0", port) {
		return nil
	}

	pids, err := listenersOnPort(port)
	if err != nil {
		return fmt.Errorf("find listeners on :%d: %w", port, err)
	}
	if len(pids) == 0 {
		return fmt.Errorf("port :%d is in use but listener PID could not be determined", port)
	}

	self := os.Getpid()
	for _, pid := range pids {
		if pid == self {
			continue
		}
		fmt.Printf("ctrl: freeing port :%d (pid %d)\n", port, pid)
		if err := terminatePID(pid); err != nil {
			return fmt.Errorf("stop pid %d on :%d: %w", pid, port, err)
		}
	}

	deadline := time.Now().Add(freeTimeout)
	for time.Now().Before(deadline) {
		if !isListening("127.0.0.1", port) && !isListening("0.0.0.0", port) {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("port :%d still in use after stopping %v", port, pids)
}

// FreeAll attempts to clear every managed port before startup.
func FreeAll(ports []int) error {
	var first error
	for _, port := range ports {
		if err := EnsureFree(port); err != nil && first == nil {
			first = err
		}
	}
	return first
}

func isListening(host string, port int) bool {
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", addr, 300*time.Millisecond)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

func terminatePID(pid int) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	selfPGID, _ := syscall.Getpgid(os.Getpid())

	if pgid, err := syscall.Getpgid(pid); err == nil && pgid != selfPGID {
		_ = syscall.Kill(-pgid, syscall.SIGTERM)
	} else {
		_ = proc.Signal(syscall.SIGTERM)
	}

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if err := proc.Signal(syscall.Signal(0)); err != nil {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	if pgid, err := syscall.Getpgid(pid); err == nil && pgid != selfPGID {
		_ = syscall.Kill(-pgid, syscall.SIGKILL)
	} else {
		_ = proc.Signal(syscall.SIGKILL)
	}
	return nil
}

func listenersOnPort(port int) ([]int, error) {
	return platformListenersOnPort(port)
}