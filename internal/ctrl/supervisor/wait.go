package supervisor

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

func waitForHealth(ctx context.Context, url string, timeout time.Duration, proc *managedProc) error {
	deadline := time.Now().Add(timeout)
	client := &http.Client{Timeout: 2 * time.Second}

	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if procExited(proc) {
			return fmt.Errorf("process exited before %s was ready", url)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		resp, err := client.Do(req)
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				if procExited(proc) {
					return fmt.Errorf("process exited after %s became reachable (port conflict?)", url)
				}
				return nil
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for %s", url)
}

func waitForPort(ctx context.Context, addr string, timeout time.Duration, proc *managedProc) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if procExited(proc) {
			return fmt.Errorf("process exited before %s was ready", addr)
		}

		conn, err := net.DialTimeout("tcp", addr, 500*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			if procExited(proc) {
				return fmt.Errorf("process exited after %s became reachable (port conflict?)", addr)
			}
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for %s", addr)
}