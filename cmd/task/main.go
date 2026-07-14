package main

import (
	"github.com/pafthang/pocketagent/internal/task"
	"github.com/pafthang/pocketagent/pkgs/common"
)

func main() {
	common.RunMain("task", task.Run)
}
