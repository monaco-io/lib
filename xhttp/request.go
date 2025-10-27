package xhttp

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/monaco-io/lib/typing/xjson"
	"github.com/monaco-io/lib/typing/xopt"
	"github.com/monaco-io/lib/typing/xxml"
	"github.com/monaco-io/lib/typing/xyaml"
)

// Interceptor 定义拦截器接口
type Interceptor interface {
	// 在请求发送前处理
	Before(req *http.Request) error
	// 在响应返回后处理
	After(resp *http.Response, req *http.Request) error
}

// InterceptorFunc 适配函数式拦截器
type InterceptorFunc struct {
	BeforeFunc func(req *http.Request) error
	AfterFunc  func(resp *http.Response, req *http.Request) error
}

var _ Interceptor = (*InterceptorFunc)(nil)

func (f *InterceptorFunc) Before(req *http.Request) error {
	if f != nil && f.BeforeFunc != nil {
		return f.BeforeFunc(req)
	}
	return nil
}
func (f *InterceptorFunc) After(resp *http.Response, req *http.Request) error {
	if f != nil && f.AfterFunc != nil {
		return f.AfterFunc(resp, req)
	}
	return nil
}

var interceptors []Interceptor

// RegisterInterceptor 注册拦截器
func RegisterInterceptors(i ...Interceptor) {
	interceptors = append(interceptors, i...)
}

const (
	ContentType           = "Content-Type"
	ContentTypeHTML       = "text/html"
	ContentTypeJSON       = "application/json"
	ContentTypeXML        = "application/xml"
	ContentTypeYAML       = "application/x-yaml"
	ContentTypeText       = "text/plain"
	ContentTypeURLEncoded = "application/x-www-form-urlencoded"
)

type Request struct {
	*http.Request
	*http.Client

	decoder
}

type decoder string

const defaultDecoder = decoderJSON

const (
	decoderJSON decoder = "json"
	decoderXML  decoder = "xml"
	decoderYAML decoder = "yaml"
	decoderText decoder = "text"
)

func Method(method string) xopt.Option[Request] {
	return func(request *Request) {
		request.Method = method
	}
}

func DecoderJSON() xopt.Option[Request] {
	return func(request *Request) {
		request.decoder = decoderJSON
	}
}

func DecoderXML() xopt.Option[Request] {
	return func(request *Request) {
		request.decoder = decoderXML
	}
}

func DecoderYAML() xopt.Option[Request] {
	return func(request *Request) {
		request.decoder = decoderYAML
	}
}

func DecoderText() xopt.Option[Request] {
	return func(request *Request) {
		request.decoder = decoderText
	}
}

func Client(client *http.Client) xopt.Option[Request] {
	return func(request *Request) {
		request.Client = client
	}
}

func NativeBody(contentType string, body io.Reader) xopt.Option[Request] {
	return func(request *Request) {
		request.Header.Set(ContentType, contentType)
		request.Body = io.NopCloser(body)
	}
}

func URLRawQuery(query url.Values) xopt.Option[Request] {
	return func(request *Request) {
		request.URL.RawQuery = query.Encode()
	}
}

func BodyJSON[T any](body T) xopt.Option[Request] {
	return func(request *Request) {
		NativeBody(ContentTypeJSON, xjson.MarshalReaderX(body))(request)
	}
}

func BodyYAML[T any](body T) xopt.Option[Request] {
	return func(request *Request) {
		NativeBody(ContentTypeYAML, xyaml.MarshalReaderX(body))(request)
	}
}

func BodyXML[T any](body T) xopt.Option[Request] {
	return func(request *Request) {
		NativeBody(ContentTypeXML, xxml.MarshalReaderX(body))(request)
	}
}

func BodyText(body string) xopt.Option[Request] {
	return func(request *Request) {
		NativeBody(ContentTypeText, strings.NewReader(body))(request)
	}
}

func BodyWWWFormURLEncoded(body url.Values) xopt.Option[Request] {
	return func(request *Request) {
		NativeBody(ContentTypeURLEncoded, strings.NewReader(body.Encode()))(request)
	}
}

func BodyMultipartForm(form *multipart.Form) xopt.Option[Request] {
	return func(request *Request) {
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)

		// Write form values
		for key, values := range form.Value {
			for _, value := range values {
				_ = writer.WriteField(key, value)
			}
		}

		// Write form files
		for key, files := range form.File {
			for _, fileHeader := range files {
				file, err := fileHeader.Open()
				if err != nil {
					continue
				}
				defer file.Close()

				part, err := writer.CreateFormFile(key, fileHeader.Filename)
				if err != nil {
					continue
				}
				_, _ = io.Copy(part, file)
			}
		}

		writer.Close()

		contentType := writer.FormDataContentType()
		NativeBody(contentType, &buf)(request)
	}
}

func BasicAuth(username, password string) xopt.Option[Request] {
	return func(request *Request) {
		request.SetBasicAuth(username, password)
	}
}

func Header(key, value string, replace ...bool) xopt.Option[Request] {
	return func(request *Request) {
		if xopt.Boolean(replace...) {
			request.Header.Set(key, value)
		} else {
			request.Header.Add(key, value)
		}
	}
}

func Headers(headers http.Header, replace ...bool) xopt.Option[Request] {
	return func(request *Request) {
		for key, values := range headers {
			for _, value := range values {
				Header(key, value, replace...)(request)
			}
		}
	}
}

func Transport(transport http.RoundTripper) xopt.Option[Request] {
	return func(request *Request) {
		request.Transport = transport
	}
}

func Timeout(timeout time.Duration) xopt.Option[Request] {
	return func(request *Request) {
		request.Timeout = timeout
	}
}

func Jar(jar http.CookieJar) xopt.Option[Request] {
	return func(request *Request) {
		request.Jar = jar
	}
}

func CheckRedirect(f func(req *http.Request, via []*http.Request) error) xopt.Option[Request] {
	return func(request *Request) {
		request.CheckRedirect = f
	}
}

func build(ctx context.Context, url string, opts ...xopt.Option[Request]) (*Request, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet, url, nil,
	)
	if err != nil {
		return nil, err
	}
	xrequest := &Request{
		Request: request,
		Client:  http.DefaultClient,
		decoder: defaultDecoder,
	}
	xopt.Apply(opts, xrequest)
	// 执行所有拦截器的Before方法
	for _, i := range interceptors {
		if err := i.Before(xrequest.Request); err != nil {
			return nil, err
		}
	}
	return xrequest, nil
}

// doWithInterceptors 执行请求并调用拦截器链
func (r *Request) doWithInterceptors(req *http.Request) (*http.Response, error) {
	resp, err := r.Do(req)
	for _, i := range interceptors {
		_ = i.After(resp, r.Request)
	}
	return resp, err
}
