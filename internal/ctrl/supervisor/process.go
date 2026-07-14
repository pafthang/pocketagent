package supervisor

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/pafthang/pocketagent/internal/ctrl/catalog"
)

type managedProc struct {
	name string
	cmd  *exec.Cmd
}

func (s *Supervisor) start(def *catalog.Service) (*managedProc, error) {
	cmd := exec.Command("go", "run", def.Package)
	cmd.Dir = s.root
	cmd.Env = s.buildEnv(def)
	cmd.Stdout = &prefixWriter{prefix: def.Name, w: os.Stdout}
	cmd.Stderr = &prefixWriter{prefix: def.Name, w: os.Stderr}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &managedProc{name: def.Name, cmd: cmd}, nil
}

func (s *Supervisor) buildEnv(def *catalog.Service) []string {
	base := map[string]string{
		"CONFIG_DIR":       s.configDir,
		"POCKETAGENT_ROOT": s.root,
	}
	for k, v := range def.Env {
		base[k] = v
	}

	merged := os.Environ()
	for k, v := range base {
		merged = overwriteEnv(merged, k, v)
	}
	return merged
}

func (s *Supervisor) shutdownStarted() {
	var wg sync.WaitGroup

	for i := len(s.order) - 1; i >= 0; i-- {
		name := s.order[i]
		proc, ok := s.procs[name]
		if !ok || proc.cmd.Process == nil {
			continue
		}

		wg.Add(1)
		go func(p *managedProc) {
			defer wg.Done()
			stopProcess(p.cmd, s.stopTimeout)
		}(proc)
	}

	wg.Wait()
}

func stopProcess(cmd *exec.Cmd, timeout time.Duration) {
	if cmd.Process == nil {
		return
	}

	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err == nil {
		_ = syscall.Kill(-pgid, syscall.SIGTERM)
	} else {
		_ = cmd.Process.Signal(syscall.SIGTERM)
	}

	waitDone := make(chan struct{})
	go func() {
		_ = cmd.Wait()
		close(waitDone)
	}()

	select {
	case <-waitDone:
	case <-time.After(timeout):
		if pgid, err := syscall.Getpgid(cmd.Process.Pid); err == nil {
			_ = syscall.Kill(-pgid, syscall.SIGKILL)
		} else {
			_ = cmd.Process.Kill()
		}
		<-waitDone
	}
}

func procExited(proc *managedProc) bool {
	if proc == nil || proc.cmd == nil || proc.cmd.Process == nil {
		return true
	}
	err := proc.cmd.Process.Signal(syscall.Signal(0))
	return err != nil
}

func overwriteEnv(env []string, key, value string) []string {
	prefix := key + "="
	for i, entry := range env {
		if len(entry) >= len(prefix) && entry[:len(prefix)] == prefix {
			env[i] = prefix + value
			return env
		}
	}
	return append(env, prefix+value)
}

type prefixWriter struct {
	prefix string
	w      io.Writer
}

func (p *prefixWriter) Write(b []byte) (int, error) {
	line := fmt.Sprintf("[%s] %s", p.prefix, string(b))
	_, err := p.w.Write([]byte(line))
	return len(b), err
}