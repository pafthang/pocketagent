package client

func stringSliceField(record map[string]interface{}, key string) []string {
	raw, ok := record[key]
	if !ok {
		return nil
	}
	switch items := raw.(type) {
	case []string:
		return items
	case []interface{}:
		out := make([]string, 0, len(items))
		for _, item := range items {
			if s, ok := item.(string); ok && s != "" {
				out = append(out, s)
			}
		}
		return out
	default:
		return nil
	}
}