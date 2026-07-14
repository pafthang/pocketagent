//go:build linux

package port

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func platformListenersOnPort(port int) ([]int, error) {
	inodes, err := socketInodesForPort(port)
	if err != nil {
		return nil, err
	}
	if len(inodes) == 0 {
		return nil, nil
	}

	pids := make(map[int]struct{})
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		fdDir := filepath.Join("/proc", entry.Name(), "fd")
		fds, err := os.ReadDir(fdDir)
		if err != nil {
			continue
		}

		for _, fd := range fds {
			link, err := os.Readlink(filepath.Join(fdDir, fd.Name()))
			if err != nil {
				continue
			}
			if !strings.HasPrefix(link, "socket:[") {
				continue
			}
			inode := strings.TrimSuffix(strings.TrimPrefix(link, "socket:["), "]")
			if _, ok := inodes[inode]; ok {
				pids[pid] = struct{}{}
			}
		}
	}

	out := make([]int, 0, len(pids))
	for pid := range pids {
		out = append(out, pid)
	}
	return out, nil
}

func socketInodesForPort(port int) (map[string]struct{}, error) {
	inodes := make(map[string]struct{})

	type table struct {
		path    string
		portHex string
	}
	tables := []table{
		{"/proc/net/tcp", portHexLE(port)},
		{"/proc/net/tcp6", portHexBE(port)},
	}

	for _, tbl := range tables {
		found, err := socketInodesFromTable(tbl.path, tbl.portHex)
		if err != nil {
			return nil, err
		}
		for inode := range found {
			inodes[inode] = struct{}{}
		}
	}

	return inodes, nil
}

func socketInodesFromTable(path, portHex string) (map[string]struct{}, error) {
	inodes := make(map[string]struct{})

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return inodes, nil
		}
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		return inodes, scanner.Err()
	}

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 10 {
			continue
		}
		if fields[3] != "0A" {
			continue
		}

		local := fields[1]
		colon := strings.LastIndex(local, ":")
		if colon < 0 {
			continue
		}
		if !strings.EqualFold(local[colon+1:], portHex) {
			continue
		}
		inodes[fields[9]] = struct{}{}
	}
	return inodes, scanner.Err()
}

func portHexLE(port int) string {
	lo := byte(port & 0xff)
	hi := byte(port >> 8)
	return fmt.Sprintf("%02X%02X", lo, hi)
}

func portHexBE(port int) string {
	return fmt.Sprintf("%04X", port)
}