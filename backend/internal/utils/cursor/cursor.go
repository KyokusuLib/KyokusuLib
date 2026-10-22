package cursor

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var ErrInvalid = errors.New("invalid cursor")

// Cursor is an opaque position in a feed ordered by (time, id) descending.
type Cursor struct {
	Time time.Time
	ID   int64
}

func Encode(t time.Time, id int64) string {
	raw := fmt.Sprintf("%d:%d", t.UTC().UnixNano(), id)
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

func Decode(value string) (Cursor, error) {
	raw, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return Cursor{}, ErrInvalid
	}

	parts := strings.SplitN(string(raw), ":", 2)
	if len(parts) != 2 {
		return Cursor{}, ErrInvalid
	}

	nanos, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return Cursor{}, ErrInvalid
	}

	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || id <= 0 {
		return Cursor{}, ErrInvalid
	}

	return Cursor{Time: time.Unix(0, nanos).UTC(), ID: id}, nil
}
