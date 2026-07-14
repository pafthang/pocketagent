package svc

import (
	"github.com/philippgille/chromem-go"
)

type searchHit struct {
	ID         string  `json:"ID"`
	Content    string  `json:"Content"`
	Similarity float32 `json:"Similarity"`
}

func defaultPerPage(perPage int) int {
	if perPage <= 0 {
		return 50
	}
	return perPage
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func toSearchResponse(results []chromem.Result) []searchHit {
	out := make([]searchHit, 0, len(results))
	for _, r := range results {
		out = append(out, searchHit{
			ID:         r.ID,
			Content:    r.Content,
			Similarity: r.Similarity,
		})
	}
	return out
}