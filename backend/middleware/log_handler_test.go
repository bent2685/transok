package middleware

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLogHandlerCreatesLogDir(t *testing.T) {
	base := t.TempDir()

	LogHandler(base)

	if _, err := os.Stat(filepath.Join(base, "logs", "gin.log")); err != nil {
		t.Fatalf("日志目录不存在时应自动创建并打开 gin.log: %v", err)
	}
}
