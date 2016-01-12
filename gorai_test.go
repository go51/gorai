package gorai_test
import (
	"testing"
	"github.com/go51/gorai"
)

func TestLoad(t *testing.T) {
	g1 := gorai.Load()
	g2 := gorai.Load()

	if g1 == nil {
		t.Error("インスタンス生成 / ロードに失敗しました。")
	}
	if g2 == nil {
		t.Error("インスタンス生成 / ロードに失敗しました。")
	}

	if g1 != g2 {
		t.Error("インスタンス生成 / ロードに失敗しました。")
	}
}

func BenchmarkLoad(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = gorai.Load()
	}
}