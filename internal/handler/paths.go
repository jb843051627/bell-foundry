package handler

import (
	"net/url"
	"strings"
)

func pathParts(path string) []string {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return nil
	}
	raw := strings.Split(trimmed, "/")
	parts := make([]string, 0, len(raw))
	for _, item := range raw {
		if item == "" {
			continue
		}
		decoded, err := url.PathUnescape(item)
		if err != nil {
			parts = append(parts, item)
			continue
		}
		parts = append(parts, decoded)
	}
	return parts
}

func idAction(path string, root string) (string, string, bool) {
	parts := pathParts(path)
	if len(parts) < 3 || parts[0] != "api" || parts[1] != root {
		return "", "", false
	}
	if len(parts) == 3 {
		return parts[2], "", true
	}
	return parts[2], parts[3], true
}
