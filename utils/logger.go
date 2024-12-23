package utils

import (
	"log/slog"
	"os"
	"sync"
)

// Logger はカスタムロガーの構造体です。
type Logger struct {
	*slog.Logger // ここがフィールド名です
	mu           sync.Mutex
}

// NewLogger は新しいロガーを作成します。
func NewLogger(logLevel slog.Level) (*Logger, error) {
	var writer slog.Handler
	writer = slog.NewJSONHandler(os.Stdout, nil)

	// ロガーの初期化
	logger := slog.New(writer)

	// Logger構造体のインスタンスを返す
	return &Logger{Logger: logger}, nil // フィールド名を修正
}

// Info は情報レベルのログを記録します。
func (l *Logger) Info(msg string, traceID string, fields ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Logger.Info(msg, append([]interface{}{"trace_id", traceID}, fields...)...)
}

// Error はエラーレベルのログを記録します。
func (l *Logger) Error(msg string, traceID string, err error, fields ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Logger.Error(msg, append([]interface{}{"trace_id", traceID, "error", err.Error()}, fields...)...)
}

// Debug はデバッグレベルのログを記録します。
func (l *Logger) Debug(msg string, traceID string, fields ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Logger.Debug(msg, append([]interface{}{"trace_id", traceID}, fields...)...)
}
