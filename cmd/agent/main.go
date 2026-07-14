package main

import (
	"github.com/pafthang/pocketagent/internal/agent"
	"github.com/pafthang/pocketagent/pkgs/common"
)

func main() {
	common.RunMain("agent", agent.Run)
}
