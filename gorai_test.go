package gorai_test

import (
	"github.com/go51/gorai"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	env := os.Getenv("GORAI_ENV")
	os.Setenv("GORAI_ENV", "framework")

	code := m.Run()

	os.Setenv("GORAI_ENV", env)

	os.Exit(code)

}

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

func TestLoadConfig(t *testing.T) {

	g := gorai.Load()
	conf := g.Config()

	if conf.Framework.WebServer.Host != "" {
		t.Errorf("設定フィアルの読み込みに失敗しました。")
	}

	if conf.Framework.WebServer.Port != "8080" {
		t.Errorf("設定フィアルの読み込みに失敗しました。")
	}

	if conf.Framework.WebServer.ReadTimeout != 30 {
		t.Errorf("設定フィアルの読み込みに失敗しました。")
	}

	if conf.Framework.WebServer.WriteTimeout != 60 {
		t.Errorf("設定フィアルの読み込みに失敗しました。")
	}
}

func BenchmarkLoad(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = gorai.Load()
	}
}

func TestRun(t *testing.T) {
	//	t.SkipNow()
	g := gorai.Load()
	g.Run()
}
