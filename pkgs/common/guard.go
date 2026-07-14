package common

import (
	"log/slog"

	"github.com/pafthang/pocketagent/pkgs/common/guard"
)

var ErrPromptRejected = guard.ErrPromptRejected

type PromptGuardMode = guard.Mode

const (
	PromptGuardBlock = guard.ModeBlock
	PromptGuardWarn  = guard.ModeWarn
)

type PromptGuardConfig = guard.Config

func LoadPromptGuardConfig() PromptGuardConfig { return guard.LoadConfig() }
func CheckPrompt(cfg PromptGuardConfig, text string) error { return guard.Check(cfg, text) }
func DetectInjectionPatterns(text string) []string         { return guard.DetectInjectionPatterns(text) }
func GuardPrompt(log *slog.Logger, cfg PromptGuardConfig, text string) error {
	return guard.Prompt(log, cfg, text)
}