package common

import (
	cfg "github.com/pafthang/pocketagent/pkgs/common/config"
	"github.com/spf13/viper"
)

type LoaderOptions = cfg.LoaderOptions
type CtrlServiceDef = cfg.CtrlServiceDef
type CtrlConfig = cfg.CtrlConfig

func Load(opts LoaderOptions, dest any) error { return cfg.Load(opts, dest) }
func LoadCtrlConfig() (*CtrlConfig, error)  { return cfg.LoadCtrlConfig() }
func FindProjectRoot() (string, error)      { return cfg.FindProjectRoot() }
func FindConfigsDir() (string, error)       { return cfg.FindConfigsDir() }
func InitRuntimeDirs(root, configOverride string) (string, error) {
	return cfg.InitRuntimeDirs(root, configOverride)
}
func ConfigFilePath(service string) (string, error) { return cfg.ConfigFilePath(service) }
func PortString(v *viper.Viper, key string) string  { return cfg.PortString(v, key) }
func SetEnvMap(env map[string]string, pairs ...string) {
	cfg.SetEnvMap(env, pairs...)
}