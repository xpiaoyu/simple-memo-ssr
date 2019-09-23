package main

import (
	"crypto/md5"
	"fmt"
	"testing"
	"time"
)

func BenchmarkGetMarkdownAndHtml(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, _, err := getMarkdownAndHtml("article/warehouse.md"); err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
}

func BenchmarkMd5(b *testing.B) {
	data := []byte("test")
	b.Logf("%X", md5.Sum(data))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%X", md5.Sum(data))
	}
	b.StopTimer()
}

func BenchmarkTimeNowNano(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = time.Now().Nanosecond() / 1000
	}
	b.StopTimer()
}
