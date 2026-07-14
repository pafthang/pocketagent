package files

import (
	"github.com/pafthang/pocketagent/internal/files/blob"
	"github.com/pafthang/pocketagent/pkgs/common"
	"github.com/spf13/viper"
)

// Config for the files service (configs/files.yaml).
type Config struct {
	Service              string `mapstructure:"service"`
	LogLevel             string `mapstructure:"log_level"`
	Port                 string `mapstructure:"port"`
	PocketBaseURL        string `mapstructure:"pocketbase_url"`
	PocketBaseAdminEmail string `mapstructure:"pocketbase_admin_email"`
	PocketBaseAdminPass  string `mapstructure:"pocketbase_admin_password"`
	MemoURL              string `mapstructure:"memo_url"`
	MemoServiceToken     string `mapstructure:"memo_service_token"`
	OllamaURL            string `mapstructure:"ollama_url"`
	EmbedModel           string `mapstructure:"embed_model"`
	FilesBackend         string `mapstructure:"files_backend"`
	FilesDataDir         string `mapstructure:"files_data_dir"`
	FilesS3Endpoint      string `mapstructure:"files_s3_endpoint"`
	FilesS3Bucket        string `mapstructure:"files_s3_bucket"`
	FilesS3AccessKey     string `mapstructure:"files_s3_access_key"`
	FilesS3SecretKey     string `mapstructure:"files_s3_secret_key"`
	FilesS3UseSSL        bool   `mapstructure:"files_s3_use_ssl"`
	FilesS3Region        string `mapstructure:"files_s3_region"`
	AuthorizeCacheSecs   int    `mapstructure:"authorize_cache_seconds"`
}

// LoadConfig reads configs/files.yaml.
func LoadConfig() (*Config, error) {
	var cfg Config
	err := common.Load(common.LoaderOptions{
		Service: "files",
		Defaults: func(v *viper.Viper) {
			v.SetDefault("service", "files")
			v.SetDefault("port", "8086")
			v.SetDefault("log_level", "info")
			v.SetDefault("pocketbase_url", "http://127.0.0.1:8090")
			v.SetDefault("memo_url", "http://127.0.0.1:8082")
			v.SetDefault("ollama_url", "http://127.0.0.1:11434")
			v.SetDefault("embed_model", "nomic-embed-text")
			v.SetDefault("files_backend", "local")
			v.SetDefault("files_data_dir", "data/files")
			v.SetDefault("authorize_cache_seconds", 30)
		},
	}, &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.Service == "" {
		cfg.Service = "files"
	}
	cfg.PocketBaseAdminEmail, cfg.PocketBaseAdminPass = common.ResolvePocketBaseAdmin(cfg.PocketBaseAdminEmail, cfg.PocketBaseAdminPass)
	cfg.MemoServiceToken = common.ResolveMemoServiceToken(cfg.MemoServiceToken)
	if err := validateFilesSecrets(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func validateFilesSecrets(cfg *Config) error {
	if err := common.ValidateRequiredSecret("POCKETBASE_ADMIN_EMAIL", cfg.PocketBaseAdminEmail); err != nil {
		return err
	}
	return common.ValidateRequiredSecret("POCKETBASE_ADMIN_PASSWORD", cfg.PocketBaseAdminPass)
}

func (c *Config) ListenAddr() string {
	return ":" + c.Port
}

func (c *Config) StoreConfig() blob.StoreConfig {
	return blob.StoreConfig{
		Backend:     c.FilesBackend,
		DataDir:     c.FilesDataDir,
		S3Endpoint:  c.FilesS3Endpoint,
		S3Bucket:    c.FilesS3Bucket,
		S3AccessKey: c.FilesS3AccessKey,
		S3SecretKey: c.FilesS3SecretKey,
		S3UseSSL:    c.FilesS3UseSSL,
		S3Region:    c.FilesS3Region,
	}
}

func (c *Config) EnvMapWithRoot(root string) map[string]string {
	env := map[string]string{"LOG_LEVEL": c.LogLevel}
	common.SetEnvMap(env,
		"PORT", c.Port,
		"POCKETBASE_URL", c.PocketBaseURL,
		"POCKETBASE_ADMIN_EMAIL", c.PocketBaseAdminEmail,
		"POCKETBASE_ADMIN_PASSWORD", c.PocketBaseAdminPass,
		"MEMO_URL", c.MemoURL,
		"MEMO_SERVICE_TOKEN", c.MemoServiceToken,
		"OLLAMA_URL", c.OllamaURL,
		"EMBED_MODEL", c.EmbedModel,
		"FILES_BACKEND", c.FilesBackend,
		"FILES_DATA_DIR", c.FilesDataDir,
	)
	_ = root
	return env
}
