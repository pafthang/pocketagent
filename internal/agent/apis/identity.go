package agentapis

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/internal/agent/identity"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	apimw "github.com/pafthang/pocketagent/pkgs/middle"
	"github.com/pafthang/pocketagent/pkgs/models"
)

type updateIdentityRequest struct {
	IdentityFile     *string `json:"identity_file"`
	SoulFile         *string `json:"soul_file"`
	StyleFile        *string `json:"style_file"`
	InstructionsFile *string `json:"instructions_file"`
	UserFile         *string `json:"user_file"`
}

func getAgentIdentityHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		agent, err := loadAgentInSpace(c, pb, c.Param("id"))
		if err != nil {
			return mapIdentityError(c, err)
		}
		return respondIdentity(c, pb, agent)
	}
}

func putAgentIdentityHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		agent, err := loadAgentInSpace(c, pb, c.Param("id"))
		if err != nil {
			return mapIdentityError(c, err)
		}
		return applyIdentityUpdate(c, pb, agent)
	}
}

func getRuntimeConfigHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		agent, err := loadAgentInSpace(c, pb, c.Param("id"))
		if err != nil {
			return mapIdentityError(c, err)
		}

		files := identity.FromAgent(agent)
		resp := map[string]interface{}{
			"id":            agent.ID,
			"space_id":      agent.SpaceID,
			"name":          agent.Name,
			"model":         agent.Model,
			"tools":         agent.Tools,
			"system_prompt": identity.CompileAgentPrompt(files),
			"identity":      files,
		}
		return c.JSON(http.StatusOK, resp)
	}
}

func respondIdentity(c echo.Context, pb *pbclient.Client, agent models.Agent) error {
	spaceID, err := requireSpaceID(c)
	if err != nil {
		return err
	}

	files := identity.FromAgent(agent)
	user, ok := apimw.UserFromContext(c)
	if ok && user.ID != "" {
		profile, err := pb.GetSpaceProfile(spaceID, user.ID)
		if err != nil {
			return mapIdentityError(c, err)
		}
		files.UserFile = profile.Content
	}
	return c.JSON(http.StatusOK, files)
}

func applyIdentityUpdate(c echo.Context, pb *pbclient.Client, agent models.Agent) error {
	spaceID, err := requireSpaceID(c)
	if err != nil {
		return err
	}

	var req updateIdentityRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if req.IdentityFile == nil && req.SoulFile == nil && req.StyleFile == nil &&
		req.InstructionsFile == nil && req.UserFile == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "no identity fields provided"})
	}

	var updated []string
	if req.IdentityFile != nil || req.SoulFile != nil || req.StyleFile != nil || req.InstructionsFile != nil {
		patched, agentUpdated := identity.ApplyPatch(agent, identity.Patch{
			IdentityFile:     req.IdentityFile,
			SoulFile:         req.SoulFile,
			StyleFile:        req.StyleFile,
			InstructionsFile: req.InstructionsFile,
		})
		var err error
		agent, err = pb.UpdateAgent(agent.ID, patched)
		if err != nil {
			return mapIdentityError(c, err)
		}
		updated = append(updated, agentUpdated...)
	}

	if req.UserFile != nil {
		user, ok := apimw.UserFromContext(c)
		if !ok || user.ID == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "authentication required to update user profile"})
		}
		if _, err := pb.UpsertSpaceProfile(models.SpaceProfile{
			SpaceID: spaceID,
			UserID:  user.ID,
			Content: *req.UserFile,
		}); err != nil {
			return mapIdentityError(c, err)
		}
		updated = append(updated, "user_file")
	}

	return c.JSON(http.StatusOK, models.IdentitySaveResponse{
		OK:      true,
		Updated: updated,
		AgentID: agent.ID,
	})
}

func mapIdentityError(c echo.Context, err error) error {
	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		return httpErr
	}
	var apiErr *pbclient.APIError
	if errors.As(err, &apiErr) {
		status := apiErr.StatusCode
		if status < 400 {
			status = http.StatusInternalServerError
		}
		return c.JSON(status, map[string]string{"error": apiErr.Error()})
	}
	return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
}
