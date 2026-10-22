package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/lanxre/kyokusulib/internal/models/dto"
	"github.com/lanxre/kyokusulib/internal/utils/conv"
	"github.com/lanxre/kyokusulib/internal/utils/cursor"
)

var (
	errNegativeLimit     = errors.New("limit must be non-negative")
	errLimitTooLarge     = errors.New("limit must not exceed 100")
	errNegativeOffset    = errors.New("offset must be non-negative")
	errOffsetUnsupported = errors.New("offset pagination is not supported; use cursor")
	errInvalidCursor     = errors.New("invalid cursor")
)

// parseFeedPagination parses the common limit/cursor query parameters. Offset
// is accepted only when zero so older clients that still send offset=0 keep
// working, while deep pages must use the keyset cursor.
func parseFeedPagination(r *http.Request, defaultLimit int) (int, *cursor.Cursor, error) {
	query := r.URL.Query()

	limit := conv.StringToInt(query.Get("limit"))
	if limit < 0 {
		return 0, nil, errNegativeLimit
	}
	if limit > 100 {
		return 0, nil, errLimitTooLarge
	}
	if limit == 0 {
		limit = defaultLimit
	}

	offset := conv.StringToInt(query.Get("offset"))
	if offset < 0 {
		return 0, nil, errNegativeOffset
	}
	if offset > 0 {
		return 0, nil, errOffsetUnsupported
	}

	raw := query.Get("cursor")
	if raw == "" {
		return limit, nil, nil
	}

	after, err := cursor.Decode(raw)
	if err != nil {
		return 0, nil, errInvalidCursor
	}
	return limit, &after, nil
}

func setCursorHeaders[T any](w http.ResponseWriter, page dto.CursorPage[T]) {
	w.Header().Set("X-Has-More", strconv.FormatBool(page.HasMore))
	if page.NextCursor != "" {
		w.Header().Set("X-Next-Cursor", page.NextCursor)
	}
}
