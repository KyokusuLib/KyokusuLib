package cursor

import (
	"testing"
	"time"
)

func TestEncodeDecodeRoundTrip(t *testing.T) {
	original := Cursor{
		Time: time.Date(2026, 10, 22, 12, 34, 56, 123456000, time.UTC),
		ID:   42,
	}

	decoded, err := Decode(Encode(original.Time, original.ID))
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	if !decoded.Time.Equal(original.Time) {
		t.Errorf("Time = %v, want %v", decoded.Time, original.Time)
	}
	if decoded.ID != original.ID {
		t.Errorf("ID = %d, want %d", decoded.ID, original.ID)
	}
}

func TestDecodeRejectsInvalidCursors(t *testing.T) {
	tests := []string{
		"",
		"not-base64",
		"bm90LWEtY3Vyc29y",
		"MTIzOjA",
		"MTIzOi0x",
		"MTIzOmFiYw",
	}

	for _, value := range tests {
		t.Run(value, func(t *testing.T) {
			if _, err := Decode(value); err == nil {
				t.Fatal("Decode() error = nil, want ErrInvalid")
			}
		})
	}
}
