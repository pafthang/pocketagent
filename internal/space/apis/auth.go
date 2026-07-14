package spaceapis

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/internal/space/audit"
	"github.com/pafthang/pocketagent/internal/space/auth"
	"github.com/pafthang/pocketagent/pkgs/httpx"
	"github.com/pafthang/pocketagent/pkgs/models"
)

func registerHandler(d Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.Bind(&req); err != nil {
			return err
		}
		if req.Email == "" || req.Password == "" {
			return httpx.RespondError(c, httpx.ErrBadRequest("email and password are required"))
		}
		if len(req.Password) < 8 {
			return httpx.RespondError(c, httpx.ErrBadRequest("password must be at least 8 characters"))
		}

		user, err := d.PB.RegisterUser(req.Email, req.Password)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		_ = issueEmailVerification(d, user)
		d.Audit.Record(c, audit.SystemAuditLog(d.PB, audit.AuditAuthRegister, user.ID, user.Email, map[string]interface{}{"email": user.Email}))

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"user":                        user,
			"email_verification_required": !user.Verified,
		})
	}
}

func loginHandler(d Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.Bind(&req); err != nil {
			return err
		}
		if req.Email == "" || req.Password == "" {
			return httpx.RespondError(c, httpx.ErrBadRequest("email and password are required"))
		}

		session, err := d.PB.AuthWithPassword(req.Email, req.Password)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if d.Cfg.RequireEmailVerification && !session.User.Verified {
			return httpx.RespondError(c, httpx.ErrForbidden("email verification required"))
		}
		d.Audit.Record(c, audit.SystemAuditLog(d.PB, audit.AuditAuthLogin, session.User.ID, session.User.Email, nil))
		return c.JSON(http.StatusOK, session)
	}
}

func refreshHandler(d Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := auth.ExtractToken(c.Request().Header.Get("Authorization"))
		if token == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "authorization required"})
		}

		session, err := d.PB.AuthRefresh(token)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if d.Cfg.RequireEmailVerification && !session.User.Verified {
			return httpx.RespondError(c, httpx.ErrForbidden("email verification required"))
		}
		return c.JSON(http.StatusOK, session)
	}
}

func authorizeHandler(d Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := auth.UserFromContext(c.Request().Context())
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}

		var req models.AuthorizeRequest
		if err := c.Bind(&req); err != nil {
			return err
		}
		if req.SpaceID == "" || req.Action == "" {
			return httpx.RespondError(c, httpx.ErrBadRequest("space_id and action are required"))
		}

		resp, err := d.Auth.Authorize(user.ID, req.SpaceID, req.Action)
		if err != nil {
			return httpx.RespondError(c, httpx.WrapInternal("authorize", err))
		}
		return c.JSON(http.StatusOK, resp)
	}
}