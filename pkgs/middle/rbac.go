package middle

import (
	"time"

	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	mwrbac "github.com/pafthang/pocketagent/pkgs/middle/rbac"
)

type PocketRBAC = mwrbac.PocketRBAC

func NewPocketRBAC(pb *pbclient.Client, authorizeCacheTTL time.Duration) *PocketRBAC {
	return mwrbac.New(pb, authorizeCacheTTL)
}