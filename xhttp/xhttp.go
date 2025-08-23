package xhttp

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/monaco-io/lib/typing/xjson"
	"github.com/monaco-io/lib/typing/xopt"
	"github.com/monaco-io/lib/typing/xstr"
	"github.com/monaco-io/lib/typing/xxml"
	"github.com/monaco-io/lib/typing/xyaml"
)

type requestid string

const requestID requestid = "x-request-id"

func NativeDo(ctx context.Context, url string, opts ...xopt.Option[Request]) (*http.Response, error) {
	xrequest, err := build(ctx, url, opts...)
	if err != nil {
		return nil, err
	}

	// Use the custom client if provided, otherwise use the default client
	client := xrequest.Client
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Do(xrequest.Request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type Response[T any] struct {
	Body T
	Code int

	contentType string
	*http.Request
}

func Do(ctx context.Context, url string, opts ...xopt.Option[Request]) (*Response[[]byte], error) {
	response, err := NativeDo(ctx, url, opts...)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, errors.New("response is nil")
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.Body != nil {
		defer response.Body.Close()
	}
	contentType := response.Header.Get(ContentType)
	if idx := strings.Index(contentType, xstr.SEMICOLON); idx != -1 {
		contentType = contentType[:idx]
	} else if contentType == "" {
		if response.Request != nil {
			reqContentType := response.Request.Header.Get(ContentType)
			if idx := strings.Index(reqContentType, xstr.SEMICOLON); idx != -1 {
				contentType = reqContentType[:idx]
			}
		}
	}
	return &Response[[]byte]{
		Body:        body,
		Code:        response.StatusCode,
		contentType: contentType,
		Request:     response.Request,
	}, nil
}

func Sugar[T any](ctx context.Context, url string, opts ...xopt.Option[Request]) (*Response[T], error) {
	response, err := Do(ctx, url, opts...)
	if err != nil {
		return nil, err
	}
	var result T
	switch response.contentType {
	case ContentTypeJSON:
		if err := xjson.Unmarshal(response.Body, &result); err != nil {
			return nil, err
		}
	case ContentTypeXML:
		if err := xxml.Unmarshal(response.Body, &result); err != nil {
			return nil, err
		}
	case ContentTypeYAML:
		if err := xyaml.Unmarshal(response.Body, &result); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unsupported content type")
	}
	return &Response[T]{Body: result, Code: response.Code}, nil
}
