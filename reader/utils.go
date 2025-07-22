package reader

import "strings"

// ParseComponentInfo parses a ComponentInfo string (e.g. "ComponentInfo{pkg/act}")
// and returns its package and activity names. Empty strings are returned when
// any part is missing or the input does not conform to the expected format.
func ParseComponentInfo(component string) (packageName, activityName string) {
	const prefix = "ComponentInfo{"
	const suffix = "}"

	// Basic validation and trimming.
	if !strings.HasPrefix(component, prefix) || !strings.HasSuffix(component, suffix) {
		return "", ""
	}

	// Extract the inside of the braces.
	inner := component[len(prefix) : len(component)-len(suffix)]

	// Split into package and activity based on the first '/' character.
	// According to examples, the slash is always present (even if followed by nothing),
	// but we handle the absence gracefully.
	parts := strings.SplitN(inner, "/", 2)

	switch len(parts) {
	case 0:
		return "", ""
	case 1:
		return strings.TrimSpace(parts[0]), ""
	default:
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}
}

// ParseCommentText trims whitespace around an XML comment text and returns
// the clean application name string. It provides a single location to handle
// any future sanitization logic for comment-based metadata.
func ParseCommentText(comment string) string {
	return strings.TrimSpace(comment)
}
