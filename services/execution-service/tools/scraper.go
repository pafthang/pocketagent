package tools

import (
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ScrapePage parses web page content
func ScrapePage(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// Extract main content
	text := doc.Find("body").Text()
	text = strings.TrimSpace(text)

	if len(text) > 1000 {
		text = text[:1000] + "..."
	}

	return fmt.Sprintf("Content from %s:\n%s", url, text), nil
}
