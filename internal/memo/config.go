package memo

import (
	"path/filepath"
	"strings"

	"github.com/pafthang/pocketagent/pkgs/common"
	"github.com/spf13/viper"
)

// Config for the memo service (configs/memo.yaml).
// Listen address: canonical field is `port` (same as other services).
// Optional override: MEMO_LISTEN env / yaml `listen` for a full bind address (e.g. 0.0.0.0:8082).
type Config struct {
	Service          string  `mapstructure:"service"`
	LogLevel         string  `mapstructure:"log_level"`
	Port             string  `mapstructure:"port"`
	Listen           string  `mapstructure:"listen"` // optional override; prefer port
	Collection       string  `mapstructure:"collection"`
	DataDir          string  `mapstructure:"data_dir"`
	PersistCompress  bool    `mapstructure:"persist_compress"`
	RAGMinSimilarity float32 `mapstructure:"rag_min_similarity"`
	RAGChunkSize     int     `mapstructure:"rag_chunk_size"`
	RAGChunkOverlap  int     `mapstructure:"rag_chunk_overlap"`
	ServiceToken     string  `mapstructure:"service_token"`
}

// LoadConfig reads configs/memo.yaml.
func LoadConfig() (*Config, error) {
	var cfg Config
	err := common.Load(common.LoaderOptions{
		Service: "memo",
		Defaults: func(v *viper.Viper) {
			v.SetDefault("service", "memo")
			v.SetDefault("port", "8082")
			v.SetDefault("log_level", "info")
			v.SetDefault("collection", "memory")
			v.SetDefault("data_dir", "data/memo")
			v.SetDefault("persist_compress", false)
			v.SetDefault("rag_min_similarity", 0.25)
			v.SetDefault("rag_chunk_size", 1000)
			v.SetDefault("rag_chunk_overlap", 150)
		},
	}, &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.Service == "" {
		cfg.Service = "memo"
	}
	if cfg.RAGMinSimilarity <= 0 {
		cfg.RAGMinSimilarity = 0.25
	}
	if cfg.RAGChunkSize <= 0 {
		cfg.RAGChunkSize = 1000
	}
	if cfg.RAGChunkOverlap < 0 {
		cfg.RAGChunkOverlap = 150
	}
	cfg.normalizeListen()
	cfg.ServiceToken = common.ResolveMemoServiceToken(cfg.ServiceToken)
	if err := validateMemoSecrets(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func validateMemoSecrets(cfg *Config) error {
	if common.IsProduction() {
		return common.ValidateRequiredSecret("MEMO_SERVICE_TOKEN", cfg.ServiceToken)
	}
	return nil
}

func (c *Config) normalizeListen() {
	if port := strings.TrimSpace(c.Port); port == "" {
		c.Port = "8082"
	}

	listen := strings.TrimSpace(c.Listen)
	if listen != "" {
		return
	}

	port := strings.TrimPrefix(strings.TrimSpace(c.Port), ":")
	c.Listen = ":" + port
}

func (c *Config) ListenAddr() string {
	if c.Listen != "" {
		return c.Listen
	}
	return ":" + strings.TrimPrefix(strings.TrimSpace(c.Port), ":")
}

func (c *Config) EnvMapWithRoot(root string) map[string]string {
	env := map[string]string{"LOG_LEVEL": c.LogLevel}
	dataDir := c.DataDir
	if dataDir != "" && !filepath.IsAbs(dataDir) {
		dataDir = filepath.Join(root, dataDir)
	}
	common.SetEnvMap(env,
		"PORT", c.Port,
		"MEMO_DATA_DIR", dataDir,
		"MEMO_SERVICE_TOKEN", c.ServiceToken,
	)
	return env
}
