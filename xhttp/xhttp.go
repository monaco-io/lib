package xhttp

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/monaco-io/lib/typing/xjson"
	"github.com/monaco-io/lib/typing/xopt"
	"github.com/monaco-io/lib/typing/xxml"
	"github.com/monaco-io/lib/typing/xyaml"
)

type requestid string

const requestID requestid = "x-request-id"

type Response[T any] struct {
	Body T   `json:"body" xml:"body" yaml:"body"`
	Code int `json:"code" xml:"code" yaml:"code"`

	*Request `json:"-" xml:"-" yaml:"-"`
}

func Do(ctx context.Context, url string, opts ...xopt.Option[Request]) (*Response[[]byte], error) {
	xrequest, err := build(ctx, url, opts...)
	if err != nil {
		return nil, err
	}

	response, err := xrequest.Client.Do(xrequest.Request)
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

	return &Response[[]byte]{
		Body:    body,
		Code:    response.StatusCode,
		Request: xrequest,
	}, nil
}

func Sugar[T any](ctx context.Context, url string, opts ...xopt.Option[Request]) (*Response[T], error) {
	response, err := Do(ctx, url, opts...)
	if err != nil {
		return nil, err
	}
	var result T
	switch response.Request.decoder {
	case decoderJSON:
		if err := xjson.Unmarshal(response.Body, &result); err != nil {
			return nil, err
		}
	case decoderXML:
		if err := xxml.Unmarshal(response.Body, &result); err != nil {
			return nil, err
		}
	case decoderYAML:
		if err := xyaml.Unmarshal(response.Body, &result); err != nil {
			return nil, err
		}
	case decoderText:
		switch any(result).(type) {
		case string:
			result = any(string(response.Body)).(T)
		case *string:
			result = any(string(response.Body)).(T)
		case []byte:
			result = any(response.Body).(T)
		case *[]byte:
			result = any(response.Body).(T)
		}
	default:
		return nil, fmt.Errorf("Sugar.Decode: unsupported type=%T, body=%s", result, response.Body)
	}
	return &Response[T]{Body: result, Code: response.Code, Request: response.Request}, nil
}
