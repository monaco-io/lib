package xhttp

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http/httptrace"
	"time"

	"github.com/monaco-io/lib/typing/xstr"
)

// TraceResult 存储追踪结果
type TraceResult struct {
	TotalDuration   time.Duration `json:"total_duration"`
	DNSDuration     time.Duration `json:"dns_duration"`
	ConnectDuration time.Duration `json:"connect_duration"`
	TlsDuration     time.Duration `json:"tls_duration"`
	FirstByteDelay  time.Duration `json:"first_byte_delay"`
	StatusCode      int           `json:"status_code"`
	Error           string        `json:"error,omitempty"`

	startTime      time.Time `json:"-"`
	dnsStart       time.Time `json:"-"`
	connectStart   time.Time `json:"-"`
	tlsStart       time.Time `json:"-"`
	firstByteStart time.Time `json:"-"`
}
type contextTraceResultKey struct{}

func withTraceContext(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Value(requestID).(string); !ok {
		ctx = context.WithValue(ctx, requestID, xstr.UUIDX())
	}
	var tr TraceResult
	ctx = context.WithValue(ctx, contextTraceResultKey{}, &tr)
	tr.startTime = time.Now()

	// 创建追踪配置
	trace := httptrace.ClientTrace{
		// DNS解析开始
		DNSStart: func(info httptrace.DNSStartInfo) {
			tr.dnsStart = time.Now()
		},

		// DNS解析完成
		DNSDone: func(info httptrace.DNSDoneInfo) {
			tr.DNSDuration = time.Since(tr.dnsStart)
		},

		// 连接开始
		ConnectStart: func(network, addr string) {
			tr.connectStart = time.Now()
		},

		// 连接完成
		ConnectDone: func(network, addr string, err error) {
			if err != nil {
				log.Printf("连接失败: %v", err)
				tr.Error = fmt.Sprintf("连接失败: %v", err)
				return
			}
			tr.ConnectDuration = time.Since(tr.connectStart)
		},

		// TLS握手开始
		TLSHandshakeStart: func() {
			tr.tlsStart = time.Now()
		},

		// TLS握手完成
		TLSHandshakeDone: func(state tls.ConnectionState, err error) {
			if err != nil {
				return
			}
			tr.TlsDuration = time.Since(tr.tlsStart)
		},

		// 准备发送请求
		WroteRequest: func(info httptrace.WroteRequestInfo) {
			tr.firstByteStart = time.Now()
		},

		// 收到第一个响应字节
		GotFirstResponseByte: func() {
			tr.FirstByteDelay = time.Since(tr.firstByteStart)
		},
	}

	// 创建带追踪的上下文
	return httptrace.WithClientTrace(ctx, &trace)
}
