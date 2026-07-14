package main

import (
	"github.com/pafthang/pocketagent/internal/gate"
	"github.com/pafthang/pocketagent/pkgs/common"
)

func main() {
	common.RunMain("gate", gate.Run)
}
