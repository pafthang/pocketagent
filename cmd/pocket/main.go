package main

import (
	"github.com/pafthang/pocketagent/internal/pocket"
	"github.com/pafthang/pocketagent/pkgs/common"
)

func main() {
	common.RunMain("pocket", pocket.Run)
}
