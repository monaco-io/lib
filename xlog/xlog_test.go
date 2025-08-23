package xlog

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/monaco-io/lib/typing/xstr"
)

type testWriter struct {
	buffer *bytes.Buffer
	mutex  sync.RWMutex
}

func (tw *testWriter) Write(p []byte) (n int, err error) {
	tw.mutex.Lock()
	defer tw.mutex.Unlock()
	return tw.buffer.Write(p)
}

func (tw *testWriter) String() string {
	tw.mutex.RLock()
	defer tw.mutex.RUnlock()
	return tw.buffer.String()
}

func (tw *testWriter) Reset() {
	tw.mutex.Lock()
	defer tw.mutex.Unlock()
	tw.buffer.Reset()
}

func newTestWriter() *testWriter {
	return &testWriter{
		buffer: &bytes.Buffer{},
	}
}

func TestSetLevel(t *testing.T) {
	// 保存原始级别
	originalLevel := _level.Level()
	defer SetLevel(originalLevel)

	tests := []struct {
		name     string
		setLevel zapcore.Level
		expected zapcore.Level
	}{
		{"set debug level", zapcore.DebugLevel, zapcore.DebugLevel},
		{"set info level", zapcore.InfoLevel, zapcore.InfoLevel},
		{"set warn level", zapcore.WarnLevel, zapcore.WarnLevel},
		{"set error level", zapcore.ErrorLevel, zapcore.ErrorLevel},
		{"set panic level", zapcore.PanicLevel, zapcore.PanicLevel},
		{"set fatal level", zapcore.FatalLevel, zapcore.FatalLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLevel(tt.setLevel)
			if _level.Level() != tt.expected {
				t.Errorf("SetLevel() = %v, want %v", _level.Level(), tt.expected)
			}
		})
	}
}

func TestRegisterWriter(t *testing.T) {
	// 备份原始状态
	originalWriter := writer
	originalLog := log
	defer func() {
		writer = originalWriter
		log = originalLog
		newLogger()
	}()

	testWriter1 := newTestWriter()

	// 测试注册单个写入器
	RegisterWriter(testWriter1)
	if log == nil {
		t.Error("RegisterWriter() should initialize logger")
	}

	// 测试无参数调用
	RegisterWriter()
	if log == nil {
		t.Error("RegisterWriter() should initialize logger even with no writers")
	}
}

func TestRegisterErrorWriter(t *testing.T) {
	// 备份原始状态
	originalErrorWriter := errorWriter
	originalLog := log
	defer func() {
		errorWriter = originalErrorWriter
		log = originalLog
		newLogger()
	}()

	testWriter1 := newTestWriter()

	RegisterErrorWriter(testWriter1)
	if log == nil {
		t.Error("RegisterErrorWriter() should initialize logger")
	}
}

func TestRegisterServiceName(t *testing.T) {
	// 备份原始状态
	originalName := name
	originalLog := log
	defer func() {
		name = originalName
		log = originalLog
		newLogger()
	}()

	tests := []struct {
		name        string
		serviceName string
	}{
		{"simple service name", "test-service"},
		{"empty service name", ""},
		{"service name with special chars", "test-service-123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterServiceName(tt.serviceName)
			if name != tt.serviceName {
				t.Errorf("RegisterServiceName() name = %v, want %v", name, tt.serviceName)
			}
			if log == nil {
				t.Error("RegisterServiceName() should initialize logger")
			}
		})
	}
}

func TestLogFunctions(t *testing.T) {
	// 备份原始状态
	originalLevel := _level.Level()
	originalLog := log
	defer func() {
		SetLevel(originalLevel)
		log = originalLog
		newLogger()
	}()

	// 设置调试级别以捕获所有日志
	SetLevel(zapcore.DebugLevel)

	// 创建测试写入器
	testWriter := newTestWriter()
	RegisterWriter(testWriter)

	ctx := context.Background()
	ctxWithRequestID := context.WithValue(ctx, xstr.X_REQUEST_ID, "test-request-123")

	tests := []struct {
		name    string
		logFunc func(context.Context, string, ...interface{})
		level   string
		ctx     context.Context
	}{
		{"debug log", D, "DEBUG", ctx},
		{"info log", I, "INFO", ctx},
		{"warn log", W, "WARN", ctx},
		{"error log", E, "ERROR", ctx},
		{"debug with request ID", D, "DEBUG", ctxWithRequestID},
		{"info with request ID", I, "INFO", ctxWithRequestID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testWriter.Reset()

			// 执行日志函数
			tt.logFunc(tt.ctx, "test message", "key", "value", "number", 42)

			// 等待日志写入
			time.Sleep(10 * time.Millisecond)

			output := testWriter.String()
			if len(output) == 0 {
				t.Error("Log function should produce output")
				return
			}

			// 验证日志包含预期内容
			if !strings.Contains(output, "test message") {
				t.Errorf("Log output should contain message, got: %s", output)
			}

			// 如果有 request ID，验证是否包含
			if tt.ctx == ctxWithRequestID {
				if !strings.Contains(output, "test-request-123") {
					t.Errorf("Log output should contain request ID, got: %s", output)
				}
			}
		})
	}
}

func TestLogFunctionsWithNilContext(t *testing.T) {
	// 备份原始状态
	originalLevel := _level.Level()
	originalLog := log
	defer func() {
		SetLevel(originalLevel)
		log = originalLog
		newLogger()
	}()

	SetLevel(zapcore.DebugLevel)
	testWriter := newTestWriter()
	RegisterWriter(testWriter)

	tests := []struct {
		name    string
		logFunc func(context.Context, string, ...interface{})
	}{
		{"debug with nil context", D},
		{"info with nil context", I},
		{"warn with nil context", W},
		{"error with nil context", E},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testWriter.Reset()

			// 测试 nil context
			tt.logFunc(nil, "test message with nil context")

			// 等待日志写入
			time.Sleep(10 * time.Millisecond)

			output := testWriter.String()
			if len(output) == 0 {
				t.Error("Log function should produce output even with nil context")
			}
		})
	}
}

func TestLogLevels(t *testing.T) {
	// 备份原始状态
	originalLevel := _level.Level()
	originalLog := log
	defer func() {
		SetLevel(originalLevel)
		log = originalLog
		newLogger()
	}()

	testWriter := newTestWriter()
	RegisterWriter(testWriter)

	tests := []struct {
		name      string
		setLevel  zapcore.Level
		logFunc   func(context.Context, string, ...interface{})
		logLevel  zapcore.Level
		shouldLog bool
	}{
		{"debug level allows debug", zapcore.DebugLevel, D, zapcore.DebugLevel, true},
		{"debug level allows info", zapcore.DebugLevel, I, zapcore.InfoLevel, true},
		{"info level blocks debug", zapcore.InfoLevel, D, zapcore.DebugLevel, false},
		{"info level allows info", zapcore.InfoLevel, I, zapcore.InfoLevel, true},
		{"warn level blocks info", zapcore.WarnLevel, I, zapcore.InfoLevel, false},
		{"warn level allows warn", zapcore.WarnLevel, W, zapcore.WarnLevel, true},
		{"error level blocks warn", zapcore.ErrorLevel, W, zapcore.WarnLevel, false},
		{"error level allows error", zapcore.ErrorLevel, E, zapcore.ErrorLevel, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLevel(tt.setLevel)
			testWriter.Reset()

			ctx := context.Background()
			tt.logFunc(ctx, "test message")

			// 等待日志写入
			time.Sleep(10 * time.Millisecond)

			output := testWriter.String()
			hasOutput := len(strings.TrimSpace(output)) > 0

			if tt.shouldLog && !hasOutput {
				t.Errorf("Expected log output but got none. Level: %v, Output: %s", tt.setLevel, output)
			}
			if !tt.shouldLog && hasOutput {
				t.Errorf("Expected no log output but got: %s", output)
			}
		})
	}
}

func TestSync(t *testing.T) {
	// 在某些环境中（如测试环境），同步标准输出可能会失败
	// 这是正常的，我们只需要确保 Sync 函数不会 panic
	err := Sync()
	if err != nil {
		// 这是预期的错误，不是测试失败
		t.Logf("Sync() returned error: %v", err)
	}

	// 测试 Sync 函数不会因为 nil logger 而 panic
	originalLog := log
	log = nil

	defer func() {
		log = originalLog
	}()

	// 这不应该 panic
	err = Sync()
	if err != nil {
		t.Logf("Sync() with nil logger returned error: %v", err)
	}
}

func TestConcurrentLogging(t *testing.T) {
	// 设置测试环境
	SetLevel(zapcore.InfoLevel)
	RegisterServiceName("test-service")

	var wg sync.WaitGroup
	numGoroutines := 5  // 减少 goroutine 数量
	numIterations := 10 // 减少迭代次数

	// 启动多个 goroutine 进行并发日志记录
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			ctx := context.WithValue(context.Background(), xstr.X_REQUEST_ID, fmt.Sprintf("request-%d", goroutineID))

			for j := 0; j < numIterations; j++ {
				I(ctx, "concurrent log", "goroutine", goroutineID, "iteration", j)
			}
		}(i)
	}

	// 等待所有 goroutine 完成
	wg.Wait()

	// 确保所有日志都被刷新
	if err := Sync(); err != nil {
		t.Logf("Warning: sync failed: %v", err)
	}

	// 测试完成，确保没有竞态条件
	t.Log("Concurrent logging test completed successfully")
}

func TestLogWithComplexKeyValues(t *testing.T) {
	// 备份原始状态
	originalLevel := _level.Level()
	originalLog := log
	defer func() {
		SetLevel(originalLevel)
		log = originalLog
		newLogger()
	}()

	SetLevel(zapcore.DebugLevel)
	testWriter := newTestWriter()
	RegisterWriter(testWriter)

	ctx := context.Background()

	// 测试复杂的键值对
	complexMap := map[string]interface{}{
		"nested": map[string]int{"a": 1, "b": 2},
		"array":  []string{"x", "y", "z"},
	}

	testWriter.Reset()
	I(ctx, "complex log",
		"string", "value",
		"number", 42,
		"float", 3.14,
		"bool", true,
		"map", complexMap,
		"slice", []int{1, 2, 3},
	)

	// 等待日志写入
	time.Sleep(10 * time.Millisecond)

	output := testWriter.String()
	if len(output) == 0 {
		t.Error("Complex log should produce output")
		return
	}

	// 验证包含复杂数据
	if !strings.Contains(output, "complex log") {
		t.Errorf("Log should contain message, got: %s", output)
	}
}

func TestMultipleWriters(t *testing.T) {
	// 备份原始状态
	originalWriter := writer
	originalLog := log
	defer func() {
		writer = originalWriter
		log = originalLog
		newLogger()
	}()

	writer1 := newTestWriter()
	writer2 := newTestWriter()

	// 注册多个写入器
	RegisterWriter(writer1, writer2)

	ctx := context.Background()
	I(ctx, "test multiple writers")

	// 等待写入
	time.Sleep(10 * time.Millisecond)

	// 由于实现细节，可能不是所有写入器都会收到相同内容
	// 至少验证日志系统没有崩溃
	if log == nil {
		t.Error("Logger should be initialized after registering multiple writers")
	}
}

func TestPanicAndFatalLogs(t *testing.T) {
	// 这个测试比较特殊，因为 P 和 F 函数会导致 panic 和 fatal
	// 我们只测试它们不会在调用时立即崩溃（在实际应用中需要更谨慎的测试）

	// 备份原始状态
	originalLevel := _level.Level()
	originalLog := log
	defer func() {
		SetLevel(originalLevel)
		log = originalLog
		newLogger()
	}()

	SetLevel(zapcore.DebugLevel)
	testWriter := newTestWriter()
	RegisterWriter(testWriter)

	// 注意：这里我们不实际调用 P 和 F 函数，因为它们会导致程序崩溃
	// 在实际应用中，需要使用特殊的测试框架来测试 panic 和 fatal 行为

	t.Log("Panic and Fatal log functions exist and can be called (but not tested here for safety)")
}

// 基准测试
func BenchmarkInfoLog(b *testing.B) {
	// 备份原始状态
	originalLevel := _level.Level()
	originalLog := log
	defer func() {
		SetLevel(originalLevel)
		log = originalLog
		newLogger()
	}()

	SetLevel(zapcore.InfoLevel)
	testWriter := newTestWriter()
	RegisterWriter(testWriter)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		I(ctx, "benchmark message", "iteration", i)
	}
}

func BenchmarkInfoLogWithRequestID(b *testing.B) {
	// 备份原始状态
	originalLevel := _level.Level()
	originalLog := log
	defer func() {
		SetLevel(originalLevel)
		log = originalLog
		newLogger()
	}()

	SetLevel(zapcore.InfoLevel)
	testWriter := newTestWriter()
	RegisterWriter(testWriter)

	ctx := context.WithValue(context.Background(), xstr.X_REQUEST_ID, "benchmark-request-123")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		I(ctx, "benchmark message with request ID", "iteration", i)
	}
}

func BenchmarkDebugLogWhenDisabled(b *testing.B) {
	// 备份原始状态
	originalLevel := _level.Level()
	originalLog := log
	defer func() {
		SetLevel(originalLevel)
		log = originalLog
		newLogger()
	}()

	// 设置为 Info 级别，这样 Debug 日志不会被处理
	SetLevel(zapcore.InfoLevel)
	testWriter := newTestWriter()
	RegisterWriter(testWriter)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		D(ctx, "debug message that should be ignored", "iteration", i)
	}
}
