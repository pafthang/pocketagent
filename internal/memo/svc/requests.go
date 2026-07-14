package svc

// AddDocumentRequest is the memo service internal ingest payload.
type AddDocumentRequest struct {
	SpaceID   string            `json:"space_id,omitempty"`
	ID        string            `json:"id"`
	Content   string            `json:"content"`
	Embedding []float32         `json:"embedding"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// SearchRequest is the memo service internal search payload.
type SearchRequest struct {
	QueryEmbedding []float32 `json:"query_embedding"`
	Limit          int       `json:"limit"`
	SpaceID        string    `json:"space_id,omitempty"`
	MinSimilarity  float32   `json:"min_similarity,omitempty"`
}