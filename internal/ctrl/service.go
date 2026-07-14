package ctrl

import (
	"flag"
	"fmt"
)

func Run() error {
	configDirFlag := flag.String("config-dir", "configs", "path to configs directory (relative to project root or absolute)")
	flag.Parse()

	deps, err := buildDeps(*configDirFlag)
	if err != nil {
		return err
	}

	deps.Supervisor.Register(deps.Services...)

	fmt.Printf("ctrl: project root %s\n", deps.Root)
	fmt.Printf("ctrl: configs dir %s\n", deps.ConfigDir)

	return deps.Supervisor.Run()
}