package posts

import (
	"errors"
	"strings"
)

func findSubstring(s, substr string) (int, int, error) {
	idx := strings.Index(s, substr)
	if idx == -1 {
		return 0, 0, errors.New("substring not found")
	}
	return idx, idx + len(substr), nil
}
