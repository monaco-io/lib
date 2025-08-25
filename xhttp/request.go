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
	"github.com/monaco-io/lib/typing/xstr"
	"github.com/monaco-io/lib/typing/xxml"
	"github.com/monaco-io/lib/typing/xyaml"
)

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
}

func Method(method string) xopt.Option[Request] {
	return func(request *Request) {
		request.Method = method
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
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Value(requestID).(string); !ok {
		ctx = context.WithValue(ctx, requestID, xstr.UUIDX())
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	xrequest := &Request{
		Request: request,
		Client:  &http.Client{},
	}
	xopt.Apply(opts, xrequest)
	return xrequest, nil
}
