package common

import sec "github.com/pafthang/pocketagent/pkgs/common/secrets"

const (
	DevSuperuserEmail    = sec.DevSuperuserEmail
	DevSuperuserPassword = sec.DevSuperuserPassword
	DevMemoServiceToken  = sec.DevMemoServiceToken
)

func IsProduction() bool { return sec.IsProduction() }
func ValidateRequiredSecret(name, value string) error {
	return sec.ValidateRequiredSecret(name, value)
}
func ResolvePocketBaseAdmin(email, password string) (string, string) {
	return sec.ResolvePocketBaseAdmin(email, password)
}
func ResolvePocketBaseSuperuser(email, password string) (string, string) {
	return sec.ResolvePocketBaseSuperuser(email, password)
}
func ResolveMemoServiceToken(token string) string { return sec.ResolveMemoServiceToken(token) }