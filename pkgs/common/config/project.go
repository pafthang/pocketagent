package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// FindProjectRoot walks up from cwd until go.mod is found.
func FindProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}

// FindConfigsDir returns the directory with per-service YAML configs.
func FindConfigsDir() (string, error) {
	if dir := os.Getenv("CONFIG_DIR"); dir != "" {
		return dir, nil
	}
	root, err := FindProjectRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "configs"), nil
}

// InitRuntimeDirs sets CONFIG_DIR and POCKETAGENT_ROOT for the current process.
func InitRuntimeDirs(root, configOverride string) (configDir string, err error) {
	if root == "" {
		root, err = FindProjectRoot()
		if err != nil {
			return "", err
		}
	}

	configDir = configOverride
	if configDir == "" {
		configDir = filepath.Join(root, "configs")
	} else if !filepath.IsAbs(configDir) {
		configDir = filepath.Join(root, configDir)
	}

	_ = os.Setenv("POCKETAGENT_ROOT", root)
	_ = os.Setenv("CONFIG_DIR", configDir)

	return configDir, nil
}

// ConfigFilePath returns the YAML path for a service config.
func ConfigFilePath(service string) (string, error) {
	dir, err := FindConfigsDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, service+".yaml"), nil
}