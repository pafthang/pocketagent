package common

import (
	"net/http"

	"github.com/pafthang/pocketagent/pkgs/common/egress"
)

type EgressConfig = egress.Config

func LoadEgressConfig() EgressConfig      { return egress.LoadConfig() }
func ValidateEgressHost(host string) error { return egress.ValidateHost(host) }
func ValidateEgressURL(rawURL string) error { return egress.ValidateURL(rawURL) }
func EgressHTTPClient() *http.Client       { return egress.HTTPClient() }