package memoapis

// IngestMemoryRequest is the public gate API body for adding memory.
type IngestMemoryRequest struct {
	ID       string            `json:"id"`
	Content  string            `json:"content"`
	Tags     []string          `json:"tags"`
	Metadata map[string]string `json:"metadata"`
}

// SearchMemoryRequest is the public gate API body for semantic search.
type SearchMemoryRequest struct {
	Query         string  `json:"query"`
	Limit         int     `json:"limit"`
	MinSimilarity float32 `json:"min_similarity"`
}

type memoryDocument struct {
	ID        string            `json:"id"`
	Content   string            `json:"content"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	CreatedAt string            `json:"created_at,omitempty"`
	Tags      []string          `json:"tags,omitempty"`
}

type memorySettings struct {
	MemoryBackend        string `json:"memory_backend"`
	MemoryUseInference   bool   `json:"memory_use_inference"`
	Mem0LLMProvider      string `json:"mem0_llm_provider"`
	Mem0LLMModel         string `json:"mem0_llm_model"`
	Mem0EmbedderProvider string `json:"mem0_embedder_provider"`
	Mem0EmbedderModel    string `json:"mem0_embedder_model"`
	Mem0VectorStore      string `json:"mem0_vector_store"`
	Mem0OllamaBaseURL    string `json:"mem0_ollama_base_url"`
	Mem0AutoLearn        bool   `json:"mem0_auto_learn"`
}