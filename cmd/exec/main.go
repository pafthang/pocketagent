package main

import (
	"github.com/pafthang/pocketagent/internal/exec"
	"github.com/pafthang/pocketagent/pkgs/common"
)

func main() {
	common.RunMain("exec", exec.Run)
}
