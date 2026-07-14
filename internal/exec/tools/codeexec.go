package tools

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func codeExec(cfg Config, args map[string]interface{}) (string, error) {
	if !cfg.CodeExecEnabled {
		return "", fmt.Errorf("code_exec is disabled (set CODE_EXEC_ENABLED=true to enable)")
	}

	code := ArgString(args, "code", "input")
	if code == "" {
		return "", fmt.Errorf("code is required")
	}

	lang := strings.ToLower(ArgString(args, "language", "lang"))
	if lang == "" {
		lang = "python"
	}

	interpreter, langArgs, err := resolveInterpreter(lang)
	if err != nil {
		return "", err
	}

	workDir, err := os.MkdirTemp("", "pocketagent-code-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(workDir)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.CodeExecTimeout)
	defer cancel()

	cmdArgs := append(append([]string{}, langArgs...), code)
	cmd := exec.CommandContext(ctx, interpreter, cmdArgs...)
	cmd.Dir = workDir
	cmd.Env = []string{
		"PATH=/usr/local/bin:/usr/bin:/bin",
		"HOME=" + workDir,
		"TMPDIR=" + workDir,
		"PYTHONNOUSERSITE=1",
		"PYTHONDONTWRITEBYTECODE=1",
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	runErr := cmd.Run()
	out := trimOutput(stdout.String(), cfg.CodeExecMaxOut)
	errOut := trimOutput(stderr.String(), cfg.CodeExecMaxOut)

	var b strings.Builder
	b.WriteString(fmt.Sprintf("code_exec (%s):\n", lang))
	if out != "" {
		b.WriteString("stdout:\n")
		b.WriteString(out)
		b.WriteString("\n")
	}
	if errOut != "" {
		b.WriteString("stderr:\n")
		b.WriteString(errOut)
		b.WriteString("\n")
	}
	if runErr != nil {
		b.WriteString("error: ")
		b.WriteString(runErr.Error())
		return b.String(), nil
	}
	if out == "" && errOut == "" {
		b.WriteString("(no output)")
	}
	return b.String(), nil
}

func resolveInterpreter(lang string) (string, []string, error) {
	switch lang {
	case "python", "python3", "py":
		if path, err := exec.LookPath("python3"); err == nil {
			return path, []string{"-c"}, nil
		}
		return "", nil, fmt.Errorf("python3 not found in PATH")
	default:
		return "", nil, fmt.Errorf("unsupported language %q (allowed: python)", lang)
	}
}

func trimOutput(s string, max int) string {
	s = strings.TrimSpace(s)
	if max <= 0 {
		max = 8192
	}
	if len(s) <= max {
		return s
	}
	return s[:max] + "\n...(truncated)"
}

// CodeExecWorkRoot returns the temp directory root used for sandbox runs.
func CodeExecWorkRoot() string {
	return filepath.Clean(os.TempDir())
}