package files

import (
	"bytes"
	"image/color"
	"image/png"
	"io"
	"testing"

	"github.com/disintegration/imaging"
)

func benchmarkSourcePNG(b *testing.B, width, height int) []byte {
	b.Helper()
	img := imaging.New(width, height, color.NRGBA{R: 120, G: 80, B: 200, A: 255})
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		b.Fatal(err)
	}
	return buf.Bytes()
}

func benchmarkResize(b *testing.B, width, height int, folderWidth, folderHeight int) {
	data := benchmarkSourcePNG(b, width, height)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		decoded, err := imaging.Decode(bytes.NewReader(data))
		if err != nil {
			b.Fatal(err)
		}
		resized := imaging.Fill(decoded, folderWidth, folderHeight, imaging.Center, imaging.Lanczos)
		if err := imaging.Encode(io.Discard, resized, imaging.JPEG, imaging.JPEGQuality(80)); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkImagePosterPipeline(b *testing.B) {
	benchmarkResize(b, 1200, 1600, 600, 900)
}

func BenchmarkImageChapterPipeline(b *testing.B) {
	benchmarkResize(b, 1600, 2400, 1200, 0)
}
