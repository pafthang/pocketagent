package spaceapis

import (
	"regexp"
	"strings"

	"github.com/pafthang/pocketagent/pkgs/models"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]{0,62}[a-z0-9])?$`)

func isValidRole(role string) bool {
	return role == models.RoleAdmin || role == models.RoleEditor || role == models.RoleViewer
}

func slugify(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return '-'
	}, s)
	s = strings.Trim(s, "-")
	if s == "" {
		return "space"
	}
	return s
}