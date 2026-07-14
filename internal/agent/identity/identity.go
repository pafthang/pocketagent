package identity

import (
	"fmt"
	"strings"

	"github.com/pafthang/pocketagent/pkgs/models"
)

const ConfigKey = "identity"

// FromAgent extracts persona files from agent storage fields.
func FromAgent(agent models.Agent) models.IdentityFiles {
	files := models.IdentityFiles{
		IdentityFile: agent.SystemPrompt,
	}
	if agent.Config == nil {
		return files
	}
	blocks, ok := agent.Config[ConfigKey].(map[string]interface{})
	if !ok {
		return files
	}
	files.SoulFile = stringVal(blocks["soul"])
	files.StyleFile = stringVal(blocks["style"])
	files.InstructionsFile = stringVal(blocks["instructions"])
	return files
}

// Patch carries optional identity file updates.
type Patch struct {
	IdentityFile     *string
	SoulFile         *string
	StyleFile        *string
	InstructionsFile *string
}

// ApplyPatch updates agent storage from a partial identity patch.
func ApplyPatch(agent models.Agent, patch Patch) (models.Agent, []string) {
	var updated []string

	if patch.IdentityFile != nil {
		agent.SystemPrompt = *patch.IdentityFile
		updated = append(updated, "identity_file")
	}

	if patch.SoulFile != nil || patch.StyleFile != nil || patch.InstructionsFile != nil {
		if agent.Config == nil {
			agent.Config = map[string]interface{}{}
		}
		blocks, _ := agent.Config[ConfigKey].(map[string]interface{})
		if blocks == nil {
			blocks = map[string]interface{}{}
		}
		if patch.SoulFile != nil {
			blocks["soul"] = *patch.SoulFile
			updated = append(updated, "soul_file")
		}
		if patch.StyleFile != nil {
			blocks["style"] = *patch.StyleFile
			updated = append(updated, "style_file")
		}
		if patch.InstructionsFile != nil {
			blocks["instructions"] = *patch.InstructionsFile
			updated = append(updated, "instructions_file")
		}
		agent.Config[ConfigKey] = blocks
	}

	return agent, updated
}

// CompileAgentPrompt merges agent-side persona blocks for the system message.
func CompileAgentPrompt(files models.IdentityFiles) string {
	var sections []string
	if s := strings.TrimSpace(files.IdentityFile); s != "" {
		sections = append(sections, "# Identity\n"+s)
	}
	if s := strings.TrimSpace(files.SoulFile); s != "" {
		sections = append(sections, "# Soul\n"+s)
	}
	if s := strings.TrimSpace(files.StyleFile); s != "" {
		sections = append(sections, "# Style\n"+s)
	}
	if s := strings.TrimSpace(files.InstructionsFile); s != "" {
		sections = append(sections, "# Instructions\n"+s)
	}
	return strings.Join(sections, "\n\n")
}

// FormatUserProfile wraps stored user profile content for prompt injection.
func FormatUserProfile(content string) string {
	content = strings.TrimSpace(content)
	if content == "" {
		return ""
	}
	return "# User profile\n" + content
}

func stringVal(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprint(v)
}
