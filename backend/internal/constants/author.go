package constants

import "regexp"

var (
	AuthorNameCleanRegex = regexp.MustCompile(`[^\w\s-]`)
	AuthorNameSpaceRegex = regexp.MustCompile(`\s+`)
)