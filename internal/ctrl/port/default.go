//go:build !linux

package port

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func platformListenersOnPort(port int) ([]int, error) {
	path, err := exec.LookPath("lsof")
	if err != nil {
		return nil, fmt.Errorf("lsof not available: %w", err)
	}

	cmd := exec.Command(path, "-nP", "-iTCP:"+strconv.Itoa(port), "-sTCP:LISTEN", "-t")
	out, err := cmd.Output()
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok && exit.ExitCode() == 1 {
			return nil, nil
		}
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	pids := make([]int, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}
		pid, err := strconv.Atoi(line)
		if err != nil {
			continue
		}
		pids = append(pids, pid)
	}
	return pids, nil
}