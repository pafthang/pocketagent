package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pafthang/pocketagent/pkgs/common"
)

// SearchWeb runs a web search using the configured provider (duckduckgo by default).
func SearchWeb(cfg Config, query string) (string, error) {
	return searchWeb(cfg, query)
}

func searchWeb(cfg Config, query string) (string, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return "", fmt.Errorf("query is required")
	}

	fn := resolveSearchProvider(cfg)
	return fn(query)
}

func resolveSearchProvider(cfg Config) func(string) (string, error) {
	switch strings.ToLower(cfg.SearchProvider) {
	case "serper":
		return serperSearch(cfg.SerperAPIKey)
	case "tavily":
		return tavilySearch(cfg.TavilyAPIKey)
	case "duckduckgo", "ddg":
		return duckDuckGoSearch
	default:
		if cfg.SerperAPIKey != "" {
			return serperSearch(cfg.SerperAPIKey)
		}
		if cfg.TavilyAPIKey != "" {
			return tavilySearch(cfg.TavilyAPIKey)
		}
		return duckDuckGoSearch
	}
}

func serperSearch(apiKey string) func(string) (string, error) {
	return func(query string) (string, error) {
		if apiKey == "" {
			return "", fmt.Errorf("SERPER_API_KEY is not configured")
		}

		body, _ := json.Marshal(map[string]string{"q": query})
		req, err := http.NewRequest(http.MethodPost, "https://google.serper.dev/search", bytes.NewReader(body))
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-KEY", apiKey)

		resp, err := searchHTTPClient().Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("serper: HTTP %d: %s", resp.StatusCode, truncate(string(respBody), 300))
		}

		var parsed struct {
			Organic []struct {
				Title   string `json:"title"`
				Link    string `json:"link"`
				Snippet string `json:"snippet"`
			} `json:"organic"`
			AnswerBox struct {
				Answer string `json:"answer"`
				Title  string `json:"title"`
			} `json:"answerBox"`
		}
		if err := json.Unmarshal(respBody, &parsed); err != nil {
			return "", err
		}

		var b strings.Builder
		b.WriteString(fmt.Sprintf("Serper results for %q:\n", query))
		if parsed.AnswerBox.Answer != "" {
			b.WriteString("- Answer: ")
			b.WriteString(parsed.AnswerBox.Answer)
			b.WriteString("\n")
		}
		limit := 5
		for i, item := range parsed.Organic {
			if i >= limit {
				break
			}
			b.WriteString(fmt.Sprintf("%d. %s\n   %s\n   %s\n", i+1, item.Title, item.Link, item.Snippet))
		}
		if parsed.AnswerBox.Answer == "" && len(parsed.Organic) == 0 {
			b.WriteString("(no results)")
		}
		return b.String(), nil
	}
}

func tavilySearch(apiKey string) func(string) (string, error) {
	return func(query string) (string, error) {
		if apiKey == "" {
			return "", fmt.Errorf("TAVILY_API_KEY is not configured")
		}

		body, _ := json.Marshal(map[string]interface{}{
			"api_key":             apiKey,
			"query":               query,
			"search_depth":        "basic",
			"include_answer":      true,
			"max_results":         5,
			"include_raw_content": false,
		})

		req, err := http.NewRequest(http.MethodPost, "https://api.tavily.com/search", bytes.NewReader(body))
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := searchHTTPClient().Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("tavily: HTTP %d: %s", resp.StatusCode, truncate(string(respBody), 300))
		}

		var parsed struct {
			Answer  string `json:"answer"`
			Results []struct {
				Title   string `json:"title"`
				URL     string `json:"url"`
				Content string `json:"content"`
			} `json:"results"`
		}
		if err := json.Unmarshal(respBody, &parsed); err != nil {
			return "", err
		}

		var b strings.Builder
		b.WriteString(fmt.Sprintf("Tavily results for %q:\n", query))
		if parsed.Answer != "" {
			b.WriteString("- Answer: ")
			b.WriteString(parsed.Answer)
			b.WriteString("\n")
		}
		for i, item := range parsed.Results {
			b.WriteString(fmt.Sprintf("%d. %s\n   %s\n   %s\n", i+1, item.Title, item.URL, truncate(item.Content, 240)))
		}
		if parsed.Answer == "" && len(parsed.Results) == 0 {
			b.WriteString("(no results)")
		}
		return b.String(), nil
	}
}

func duckDuckGoSearch(query string) (string, error) {
	endpoint := "https://api.duckduckgo.com/?q=" + url.QueryEscape(query) + "&format=json&no_redirect=1&no_html=1"
	resp, err := searchHTTPClient().Get(endpoint)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Abstract      string `json:"Abstract"`
		AbstractText  string `json:"AbstractText"`
		Heading       string `json:"Heading"`
		RelatedTopics []struct {
			Text string `json:"Text"`
		} `json:"RelatedTopics"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	summary := result.AbstractText
	if summary == "" {
		summary = result.Abstract
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("DuckDuckGo results for %q:\n", query))
	if summary != "" {
		b.WriteString("- ")
		b.WriteString(summary)
		b.WriteString("\n")
	}
	for i, topic := range result.RelatedTopics {
		if i >= 3 || topic.Text == "" {
			continue
		}
		b.WriteString(fmt.Sprintf("- %s\n", topic.Text))
	}
	if summary == "" && len(result.RelatedTopics) == 0 {
		b.WriteString("(no results; configure SERPER_API_KEY or TAVILY_API_KEY for richer search)")
	}
	return b.String(), nil
}

func searchHTTPClient() *http.Client {
	return common.EgressHTTPClient()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
