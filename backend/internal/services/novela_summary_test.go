package service

import (
	"context"
	"testing"
)

func TestGetNovelaSummaryOmitsChapters(t *testing.T) {
	s := newAuthTestService(t)
	ctx := context.Background()

	const novelaID = 7
	summary, err := s.GetNovelaById(ctx, novelaID, 0, false)
	if err != nil {
		t.Fatalf("get summary: %v", err)
	}

	if len(summary.Volumes) == 0 {
		t.Fatal("expected volumes")
	}
	if summary.FirstChapterID == nil || *summary.FirstChapterID == "" {
		t.Fatal("expected first_chapter_id")
	}

	total := 0
	for _, volume := range summary.Volumes {
		if len(volume.Chapters) != 0 {
			t.Fatalf("volume %s: expected no chapters in summary, got %d", volume.ID, len(volume.Chapters))
		}
		total += volume.ChaptersCount
	}
	if total == 0 {
		t.Fatal("expected chapters_count > 0")
	}

	full, err := s.GetNovelaById(ctx, novelaID, 0, true)
	if err != nil {
		t.Fatalf("get full: %v", err)
	}

	fullTotal := 0
	for _, volume := range full.Volumes {
		fullTotal += len(volume.Chapters)
		if volume.ChaptersCount != len(volume.Chapters) {
			t.Fatalf(
				"volume %s: chapters_count=%d, chapters=%d",
				volume.ID, volume.ChaptersCount, len(volume.Chapters),
			)
		}
	}
	if fullTotal != total {
		t.Fatalf("summary total=%d, full total=%d", total, fullTotal)
	}
	if len(full.Volumes) == 0 || len(full.Volumes[0].Chapters) == 0 {
		t.Fatal("expected chapters in full response")
	}
	if full.FirstChapterID == nil || *full.FirstChapterID != full.Volumes[0].Chapters[0].ID {
		t.Fatalf("first_chapter_id does not match first chapter")
	}
}

func TestGetVolumeChapters(t *testing.T) {
	s := newAuthTestService(t)
	ctx := context.Background()

	full, err := s.GetNovelaById(ctx, 7, 0, false)
	if err != nil {
		t.Fatalf("get summary: %v", err)
	}

	var volumeID string
	var expected int
	for _, volume := range full.Volumes {
		if volume.ChaptersCount > 0 {
			volumeID = volume.ID
			expected = volume.ChaptersCount
			break
		}
	}
	if volumeID == "" {
		t.Skip("no volume with chapters in test database")
	}

	chapters, err := s.GetVolumeChapters(ctx, volumeID, 0)
	if err != nil {
		t.Fatalf("get volume chapters: %v", err)
	}
	if len(chapters) != expected {
		t.Fatalf("expected %d chapters, got %d", expected, len(chapters))
	}
	for i, chapter := range chapters {
		if chapter.ID == "" {
			t.Fatalf("chapter %d: empty id", i)
		}
		if i > 0 && chapters[i-1].Number > chapter.Number {
			t.Fatalf("chapters are not sorted by number")
		}
	}
}
