package gobar_test

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/talwat/gobar"
)

func BenchmarkRender(b *testing.B) {
	bar := gobar.NewBar(0, 10000, "benchmark", "done!")

	for i := 0; i < int(bar.Total); i++ {
		bar.Increment(1)
	}
}

func TestBasic(t *testing.T) {
	t.Parallel()

	bar := gobar.NewBar(0, 10, "basic", "done!")

	for i := 0; i < int(bar.Total); i++ {
		time.Sleep(10 * time.Millisecond)
		bar.Increment(1)
	}
}

func TestIO(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequest("GET", "https://dl.google.com/go/go1.14.2.src.tar.gz", nil)
	resp, _ := http.DefaultClient.Do(req)

	bar := gobar.NewBar(0, resp.ContentLength, "io", "done!")

	defer resp.Body.Close()

	io.Copy(io.MultiWriter(bar), resp.Body)
}
