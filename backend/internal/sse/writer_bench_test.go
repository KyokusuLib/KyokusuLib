package sse

import (
	"net/http/httptest"
	"testing"
)

func BenchmarkSSEWriterSendEvent(b *testing.B) {
	data := []byte(`{"id":7,"messageId":42,"text":"hello world","createdAt":"2026-10-22T12:00:00Z"}`)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		sw, ok := NewSSEWriter(rec)
		if !ok {
			b.Fatal("response writer does not support flushing")
		}
		if err := sw.SendEvent("tg_post", 7, data); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSSEWriterHeartbeat(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		sw, ok := NewSSEWriter(rec)
		if !ok {
			b.Fatal("response writer does not support flushing")
		}
		if err := sw.SendHeartbeat(); err != nil {
			b.Fatal(err)
		}
	}
}
