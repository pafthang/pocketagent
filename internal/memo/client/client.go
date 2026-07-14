package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pafthang/pocketagent/internal/memo/auth"
	"github.com/pafthang/pocketagent/internal/memo/svc"
)

const defaultBaseURL = "http://127.0.0.1:8082"

// Client is an HTTP client for the memo vector memory service.
type Client struct {
	BaseURL       string
	ServiceToken  string
	HTTP          *http.Client
	MinSimilarity float32
	ChunkSize     int
	ChunkOverlap  int
}

// New creates a memo service client.
func New(baseURL string, serviceToken ...string) *Client {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	token := ""
	if len(serviceToken) > 0 {
		token = serviceToken[0]
	}
	return &Client{
		BaseURL:       baseURL,
		ServiceToken:  token,
		HTTP:          &http.Client{},
		MinSimilarity: 0.25,
		ChunkSize:     1000,
		ChunkOverlap:  150,
	}
}

// WithRAGOptions configures retrieval quality defaults.
func (c *Client) WithRAGOptions(minSimilarity float32, chunkSize, chunkOverlap int) *Client {
	if minSimilarity > 0 {
		c.MinSimilarity = minSimilarity
	}
	if chunkSize > 0 {
		c.ChunkSize = chunkSize
	}
	if chunkOverlap >= 0 {
		c.ChunkOverlap = chunkOverlap
	}
	return c
}

type searchResult struct {
	ID         string  `json:"ID"`
	Content    string  `json:"Content"`
	Similarity float32 `json:"Similarity"`
}

// Document is a memo search hit.
type Document struct {
	ID         string
	Content    string
	Similarity float32
	Metadata   map[string]string
}

// DocumentRecord is a stored memo document.
type DocumentRecord struct {
	ID       string            `json:"id"`
	Content  string            `json:"content"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// ListResult is a paginated memo document list.
type ListResult struct {
	Documents []DocumentRecord
	Total     int
	Page      int
	PerPage   int
}

// CollectionStats summarizes a space collection.
type CollectionStats struct {
	Collection    string `json:"collection"`
	DocumentCount int    `json:"document_count"`
	ContentBytes  int    `json:"content_bytes"`
}

// Add stores a document in the default collection.
func (c *Client) Add(ctx context.Context, id, content string, embedding []float32) error {
	return c.AddScoped(ctx, "", id, content, embedding, nil)
}

// AddScoped stores a document in a space-specific chromem collection.
func (c *Client) AddScoped(ctx context.Context, spaceID, id, content string, embedding []float32, metadata map[string]string) error {
	body, err := json.Marshal(svc.AddDocumentRequest{
		SpaceID:   spaceID,
		ID:        id,
		Content:   content,
		Embedding: embedding,
		Metadata:  metadata,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/documents", bytes.NewReader(body))
	if err != nil {
		return err
	}
	c.applyAuth(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return readHTTPError(resp)
	}
	return nil
}

// Search returns document contents similar to the query embedding.
func (c *Client) Search(ctx context.Context, queryEmbedding []float32, limit int) ([]string, error) {
	docs, err := c.SearchDocuments(ctx, queryEmbedding, limit)
	if err != nil {
		return nil, err
	}
	return documentContents(docs), nil
}

// SearchDocuments returns memo documents similar to the query embedding.
func (c *Client) SearchDocuments(ctx context.Context, queryEmbedding []float32, limit int) ([]Document, error) {
	return c.SearchDocumentsScoped(ctx, queryEmbedding, "", limit)
}

// SearchDocumentsScoped returns documents from a per-space collection.
func (c *Client) SearchDocumentsScoped(ctx context.Context, queryEmbedding []float32, spaceID string, limit int) ([]Document, error) {
	if limit <= 0 {
		limit = 5
	}

	body, err := json.Marshal(svc.SearchRequest{
		QueryEmbedding: queryEmbedding,
		Limit:          limit,
		SpaceID:        spaceID,
		MinSimilarity:  c.MinSimilarity,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/search", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	c.applyAuth(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, readHTTPError(resp)
	}

	var results []searchResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}

	docs := make([]Document, 0, len(results))
	for _, r := range results {
		docs = append(docs, Document{ID: r.ID, Content: r.Content, Similarity: r.Similarity})
	}
	return docs, nil
}

// ListDocuments returns paginated documents for a space.
func (c *Client) ListDocuments(ctx context.Context, spaceID string, page, perPage int) (ListResult, error) {
	q := url.Values{}
	if spaceID != "" {
		q.Set("space_id", spaceID)
	}
	if page > 0 {
		q.Set("page", strconv.Itoa(page))
	}
	if perPage > 0 {
		q.Set("per_page", strconv.Itoa(perPage))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/documents?"+q.Encode(), nil)
	if err != nil {
		return ListResult{}, err
	}
	c.applyAuth(req)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return ListResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ListResult{}, readHTTPError(resp)
	}

	var payload struct {
		Documents []struct {
			ID       string            `json:"id"`
			Content  string            `json:"content"`
			Metadata map[string]string `json:"metadata"`
		} `json:"documents"`
		Total   int `json:"total"`
		Page    int `json:"page"`
		PerPage int `json:"per_page"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return ListResult{}, err
	}

	out := ListResult{
		Total:   payload.Total,
		Page:    payload.Page,
		PerPage: payload.PerPage,
	}
	for _, doc := range payload.Documents {
		out.Documents = append(out.Documents, DocumentRecord{
			ID:       doc.ID,
			Content:  doc.Content,
			Metadata: doc.Metadata,
		})
	}
	return out, nil
}

// GetDocument returns a single document by ID.
func (c *Client) GetDocument(ctx context.Context, spaceID, id string) (DocumentRecord, error) {
	q := url.Values{}
	if spaceID != "" {
		q.Set("space_id", spaceID)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/documents/"+url.PathEscape(id)+"?"+q.Encode(), nil)
	if err != nil {
		return DocumentRecord{}, err
	}
	c.applyAuth(req)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return DocumentRecord{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return DocumentRecord{}, fmt.Errorf("memo: document not found")
	}
	if resp.StatusCode != http.StatusOK {
		return DocumentRecord{}, readHTTPError(resp)
	}

	var doc DocumentRecord
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return DocumentRecord{}, err
	}
	return doc, nil
}

// DeleteDocument removes a document by ID.
func (c *Client) DeleteDocument(ctx context.Context, spaceID, id string) error {
	q := url.Values{}
	if spaceID != "" {
		q.Set("space_id", spaceID)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.BaseURL+"/documents/"+url.PathEscape(id)+"?"+q.Encode(), nil)
	if err != nil {
		return err
	}
	c.applyAuth(req)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return readHTTPError(resp)
	}
	return nil
}

// Stats returns collection statistics for a space.
func (c *Client) Stats(ctx context.Context, spaceID string) (CollectionStats, error) {
	q := url.Values{}
	if spaceID != "" {
		q.Set("space_id", spaceID)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/stats?"+q.Encode(), nil)
	if err != nil {
		return CollectionStats{}, err
	}
	c.applyAuth(req)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return CollectionStats{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CollectionStats{}, readHTTPError(resp)
	}

	var stats CollectionStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return CollectionStats{}, err
	}
	return stats, nil
}

func documentContents(docs []Document) []string {
	contents := make([]string, 0, len(docs))
	for _, doc := range docs {
		if doc.Content != "" {
			contents = append(contents, doc.Content)
		}
	}
	return contents
}

func (c *Client) applyAuth(req *http.Request) {
	if c != nil && strings.TrimSpace(c.ServiceToken) != "" {
		req.Header.Set(auth.HeaderServiceToken, c.ServiceToken)
	}
}

func readHTTPError(resp *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if len(body) == 0 {
		return fmt.Errorf("memo: HTTP %d", resp.StatusCode)
	}
	return fmt.Errorf("memo: HTTP %d: %s", resp.StatusCode, string(body))
}
