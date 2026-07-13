package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// WebSearchTool performs real web search using a public API or scraping (placeholder)
func WebSearch(query string) (string, error) {
	// Example using a public search API (in real project use Serper, Tavily, etc.)
	resp, err := http.Get("https://api.duckduckgo.com/?q=" + query + "&format=json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	return fmt.Sprintf("Search results for '%s': %v", query, result["Abstract"]), nil
}
