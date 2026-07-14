package spaceapis

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/internal/space/audit"
	"github.com/pafthang/pocketagent/internal/space/auth"
	"github.com/pafthang/pocketagent/pkgs/httpx"
	"github.com/pafthang/pocketagent/pkgs/models"
)

func requestVerificationHandler(d Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := auth.UserFromContext(c.Request().Context())
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}
		if user.Verified {
			return c.JSON(http.StatusOK, map[string]string{"status": "already_verified"})
		}

		token, err := auth.GenerateToken()
		if err != nil {
			return httpx.RespondError(c, httpx.WrapInternal("generate token", err))
		}
		expiresAt := time.Now().UTC().Add(time.Duration(d.Cfg.VerificationTTLHours) * time.Hour).Format(time.RFC3339)
		if err := d.PB.CreateEmailVerification(user.ID, user.Email, auth.HashToken(token), expiresAt); err != nil {
			return httpx.MapPocketError(c, err)
		}
		logVerificationLink(d.Log, d.Cfg.PublicBaseURL, token)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":     "verification_sent",
			"expires_at": expiresAt,
		})
	}
}

func verifyEmailHandler(d Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req struct {
			Token string `json:"token"`
		}
		if err := c.Bind(&req); err != nil {
			return err
		}
		if req.Token == "" {
			return httpx.RespondError(c, httpx.ErrBadRequest("token is required"))
		}

		record, err := d.PB.FindEmailVerificationByTokenHash(auth.HashToken(req.Token))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if verificationExpired(stringFieldFromRecord(record, "expires_at")) {
			return httpx.RespondError(c, httpx.ErrBadRequest("verification token expired"))
		}

		userID := stringFieldFromRecord(record, "user_id")
		email := stringFieldFromRecord(record, "email")
		if err := d.PB.SetUserVerified(userID, true); err != nil {
			return httpx.MapPocketError(c, err)
		}
		if err := d.PB.MarkEmailVerificationDone(stringFieldFromRecord(record, "id")); err != nil {
			return httpx.MapPocketError(c, err)
		}

		d.Audit.Record(c, audit.SystemAuditLog(d.PB, audit.AuditEmailVerified, userID, email, nil))

		return c.JSON(http.StatusOK, map[string]string{"status": "verified"})
	}
}

func issueEmailVerification(d Deps, user models.AuthUser) error {
	if user.Verified {
		return nil
	}
	token, err := auth.GenerateToken()
	if err != nil {
		return err
	}
	expiresAt := time.Now().UTC().Add(time.Duration(d.Cfg.VerificationTTLHours) * time.Hour).Format(time.RFC3339)
	if err := d.PB.CreateEmailVerification(user.ID, user.Email, auth.HashToken(token), expiresAt); err != nil {
		return err
	}
	logVerificationLink(d.Log, d.Cfg.PublicBaseURL, token)
	return nil
}

func verificationExpired(expiresAt string) bool {
	if expiresAt == "" {
		return true
	}
	t, err := time.Parse(time.RFC3339, expiresAt)
	if err != nil {
		return true
	}
	return time.Now().UTC().After(t)
}

func stringFieldFromRecord(record map[string]interface{}, key string) string {
	if v, ok := record[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
