package config

import (
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

func bindCommonEnv(v *viper.Viper) {
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	_ = v.BindEnv("log_level", "LOG_LEVEL")
	_ = v.BindEnv("port", "PORT")
	_ = v.BindEnv("nats.url", "NATS_URL")
	_ = v.BindEnv("nats_url", "NATS_URL")
	_ = v.BindEnv("ollama.url", "OLLAMA_URL")
	_ = v.BindEnv("ollama_url", "OLLAMA_URL")
	_ = v.BindEnv("pocketbase.url", "POCKETBASE_URL")
	_ = v.BindEnv("pocketbase_url", "POCKETBASE_URL")
	_ = v.BindEnv("memo.url", "MEMO_URL")
	_ = v.BindEnv("memo_url", "MEMO_URL")
	_ = v.BindEnv("memo.service_token", "MEMO_SERVICE_TOKEN")
	_ = v.BindEnv("memo_service_token", "MEMO_SERVICE_TOKEN")
	_ = v.BindEnv("service_token", "MEMO_SERVICE_TOKEN")
	_ = v.BindEnv("space.url", "SPACE_URL")
	_ = v.BindEnv("space_url", "SPACE_URL")
	_ = v.BindEnv("agent.url", "AGENT_URL")
	_ = v.BindEnv("agent_url", "AGENT_URL")
	_ = v.BindEnv("store_dir", "NATS_STORE_DIR")
	_ = v.BindEnv("http_port", "HTTP_PORT")
	_ = v.BindEnv("otel_exporter_otlp_endpoint", "OTEL_EXPORTER_OTLP_ENDPOINT")
	_ = v.BindEnv("otel_service_name", "OTEL_SERVICE_NAME")
	_ = v.BindEnv("data_dir", "POCKETBASE_DATA_DIR")
	_ = v.BindEnv("listen", "MEMO_LISTEN")
	_ = v.BindEnv("collection", "MEMO_COLLECTION")
	_ = v.BindEnv("data_dir", "MEMO_DATA_DIR")
	_ = v.BindEnv("persist_compress", "MEMO_PERSIST_COMPRESS")
	_ = v.BindEnv("timeout_sec", "TASK_TIMEOUT_SEC")
	_ = v.BindEnv("max_subtasks", "TASK_MAX_SUBTASKS")
	_ = v.BindEnv("health_port", "HEALTH_PORT")
	_ = v.BindEnv("embed_model", "EMBED_MODEL")
	_ = v.BindEnv("llm_model", "LLM_MODEL")
	_ = v.BindEnv("pocketbase_admin_email", "POCKETBASE_ADMIN_EMAIL")
	_ = v.BindEnv("pocketbase_admin_password", "POCKETBASE_ADMIN_PASSWORD")
	_ = v.BindEnv("superuser_email", "POCKETBASE_SUPERUSER_EMAIL")
	_ = v.BindEnv("superuser_password", "POCKETBASE_SUPERUSER_PASSWORD")
	_ = v.BindEnv("rate_limit_enabled", "RATE_LIMIT_ENABLED")
	_ = v.BindEnv("rate_limit_per_minute", "RATE_LIMIT_PER_MINUTE")
	_ = v.BindEnv("rate_limit_burst", "RATE_LIMIT_BURST")
	_ = v.BindEnv("auth_rate_limit_per_minute", "AUTH_RATE_LIMIT_PER_MINUTE")
	_ = v.BindEnv("auth_rate_limit_burst", "AUTH_RATE_LIMIT_BURST")
	_ = v.BindEnv("public_base_url", "PUBLIC_BASE_URL")
	_ = v.BindEnv("require_email_verification", "REQUIRE_EMAIL_VERIFICATION")
	_ = v.BindEnv("invite_ttl_hours", "INVITE_TTL_HOURS")
	_ = v.BindEnv("verification_ttl_hours", "VERIFICATION_TTL_HOURS")
	_ = v.BindEnv("stream_llm_tokens", "STREAM_LLM_TOKENS")
	_ = v.BindEnv("search_provider", "SEARCH_PROVIDER")
	_ = v.BindEnv("serper_api_key", "SERPER_API_KEY")
	_ = v.BindEnv("tavily_api_key", "TAVILY_API_KEY")
	_ = v.BindEnv("code_exec_enabled", "CODE_EXEC_ENABLED")
	_ = v.BindEnv("code_exec_timeout_sec", "CODE_EXEC_TIMEOUT_SEC")
	_ = v.BindEnv("mcp_servers", "MCP_SERVERS")
	_ = v.BindEnv("prompt_guard_enabled", "PROMPT_GUARD_ENABLED")
	_ = v.BindEnv("prompt_guard_mode", "PROMPT_GUARD_MODE")
	_ = v.BindEnv("prompt_max_length", "PROMPT_MAX_LENGTH")
	_ = v.BindEnv("egress_allowlist", "EGRESS_ALLOWLIST")
	_ = v.BindEnv("egress_allowlist_enabled", "EGRESS_ALLOWLIST_ENABLED")
	_ = v.BindEnv("rag_min_similarity", "MEMO_RAG_MIN_SIMILARITY")
	_ = v.BindEnv("rag_chunk_size", "MEMO_RAG_CHUNK_SIZE")
	_ = v.BindEnv("rag_chunk_overlap", "MEMO_RAG_CHUNK_OVERLAP")
	_ = v.BindEnv("files_data_dir", "FILES_DATA_DIR")
	_ = v.BindEnv("files_backend", "FILES_BACKEND")
}

// PortString normalizes port from viper (string or int).
func PortString(v *viper.Viper, key string) string {
	if s := v.GetString(key); s != "" {
		return s
	}
	if p := v.GetInt(key); p > 0 {
		return strconv.Itoa(p)
	}
	return ""
}

// SetEnvMap adds non-empty key/value pairs.
func SetEnvMap(env map[string]string, pairs ...string) {
	for i := 0; i+1 < len(pairs); i += 2 {
		if pairs[i+1] != "" {
			env[pairs[i]] = pairs[i+1]
		}
	}
}