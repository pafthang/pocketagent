package chunk

import "strings"

// Text splits long text into overlapping chunks for better retrieval.
func Text(text string, maxChars, overlap int) []string {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}
	if maxChars <= 0 {
		maxChars = 1000
	}
	if overlap < 0 {
		overlap = 0
	}
	if len(text) <= maxChars {
		return []string{text}
	}

	var chunks []string
	start := 0
	for start < len(text) {
		end := start + maxChars
		if end > len(text) {
			end = len(text)
		} else {
			if cut := strings.LastIndex(text[start:end], "\n\n"); cut > maxChars/2 {
				end = start + cut
			} else if cut := strings.LastIndexAny(text[start:end], ".!?"); cut > maxChars/2 {
				end = start + cut + 1
			}
		}

		part := strings.TrimSpace(text[start:end])
		if part != "" {
			chunks = append(chunks, part)
		}
		if end >= len(text) {
			break
		}
		next := end - overlap
		if next <= start {
			next = end
		}
		start = next
	}
	return chunks
}