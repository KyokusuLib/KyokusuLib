package dto

// CursorPage is an internal pagination result. The HTTP layer serializes Items
// and exposes NextCursor/HasMore via response headers, keeping the existing
// array response shape for compatibility.
type CursorPage[T any] struct {
	Items      []T
	NextCursor string
	HasMore    bool
}
