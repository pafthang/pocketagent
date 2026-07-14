package main

import (
	"github.com/pafthang/pocketagent/internal/memo"
	"github.com/pafthang/pocketagent/pkgs/common"
)

func main() {
	common.RunMain("memo", memo.Run)
}
