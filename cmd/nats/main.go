package main

import (
	"github.com/pafthang/pocketagent/internal/nats"
	"github.com/pafthang/pocketagent/pkgs/common"
)

func main() {
	common.RunMain("nats", nats.Run)
}
