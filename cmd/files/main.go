package main

import (
	"github.com/pafthang/pocketagent/internal/files"
	"github.com/pafthang/pocketagent/pkgs/common"
)

func main() {
	common.RunMain("files", files.Run)
}