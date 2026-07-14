package main

import (
	"github.com/pafthang/pocketagent/internal/ctrl"
	"github.com/pafthang/pocketagent/pkgs/common"
)

func main() {
	common.RunMain("ctrl", ctrl.Run)
}
