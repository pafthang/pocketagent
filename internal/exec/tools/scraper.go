package tools

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pafthang/pocketagent/pkgs/common"
)

func scrapePage(rawURL string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return "", fmt.Errorf("url is required")
	}
	if err := common.ValidateEgressURL(rawURL); err != nil {
		return "", err
	}

	resp, err := common.EgressHTTPClient().Get(rawURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("scrape_page: HTTP %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	text := strings.TrimSpace(doc.Find("body").Text())
	if len(text) > 4000 {
		text = text[:4000] + "..."
	}

	return fmt.Sprintf("Content from %s:\n%s", rawURL, text), nil
}
