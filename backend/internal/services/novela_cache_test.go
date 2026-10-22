package service

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/lanxre/kyokusulib/internal/models/dto"
	"github.com/redis/go-redis/v9"
)

func TestNovelaCacheHitAndInvalidate(t *testing.T) {
	ctx := context.Background()

	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skipf("redis is not available: %v", err)
	}
	t.Cleanup(func() { _ = rdb.Close() })

	s := newAuthTestServiceWithRedis(t, rdb)

	filters := dto.NovelaFilters{
		Search: fmt.Sprintf("cache-test-%d", time.Now().UnixNano()),
		Limit:  20,
	}

	if _, _, err := s.GetNovelas(ctx, 0, filters); err != nil {
		t.Fatalf("first list request: %v", err)
	}
	if hits, misses := s.CacheStats(); hits != 0 || misses != 1 {
		t.Fatalf("after first list request: hits=%d misses=%d, want 0/1", hits, misses)
	}

	if _, _, err := s.GetNovelas(ctx, 0, filters); err != nil {
		t.Fatalf("second list request: %v", err)
	}
	if hits, misses := s.CacheStats(); hits != 1 || misses != 1 {
		t.Fatalf("after second list request: hits=%d misses=%d, want 1/1", hits, misses)
	}

	s.InvalidateCache(ctx)
	if _, _, err := s.GetNovelas(ctx, 0, filters); err != nil {
		t.Fatalf("list request after invalidation: %v", err)
	}
	if hits, misses := s.CacheStats(); hits != 1 || misses != 2 {
		t.Fatalf("after invalidation: hits=%d misses=%d, want 1/2", hits, misses)
	}

	// Page cache: start from a fresh generation so a previous test run cannot
	// satisfy the first request.
	s.InvalidateCache(ctx)
	beforeHits, beforeMisses := s.CacheStats()

	if _, err := s.GetNovelaById(ctx, 7, 0, false); err != nil {
		t.Fatalf("first page request: %v", err)
	}
	if _, err := s.GetNovelaById(ctx, 7, 0, false); err != nil {
		t.Fatalf("second page request: %v", err)
	}
	hits, misses := s.CacheStats()
	if hits-beforeHits != 1 || misses-beforeMisses != 1 {
		t.Fatalf(
			"page cache: hits delta=%d misses delta=%d, want 1/1",
			hits-beforeHits, misses-beforeMisses,
		)
	}

	s.InvalidateCache(ctx)
	if _, err := s.GetNovelaById(ctx, 7, 0, false); err != nil {
		t.Fatalf("page request after invalidation: %v", err)
	}
	hits, misses = s.CacheStats()
	if hits-beforeHits != 1 || misses-beforeMisses != 2 {
		t.Fatalf(
			"page cache after invalidation: hits delta=%d misses delta=%d, want 1/2",
			hits-beforeHits, misses-beforeMisses,
		)
	}
}
