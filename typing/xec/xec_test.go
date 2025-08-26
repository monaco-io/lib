package xec

import (
	"errors"
	"strings"
	"testing"
)

func TestError_New(t *testing.T) {
	type fields struct {
		Message []string
		Code    int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"case 1", fields{[]string{"not found"}, 404}, "[404] [not found]"},
		{"case 2", fields{[]string{"internal error"}, 500}, "[500] [internal error]"},
		{"case 3", fields{[]string{"internal1", "internal2"}, 500}, "[500] [internal1 internal2]"},
		{"case 4", fields{[]string{}, 404}, "[404] []"},
		{"case 5", fields{nil, 404}, "[404] []"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New(tt.fields.Code, tt.fields.Message...)
			if got := e.Error(); got != tt.want {
				t.Errorf("Error.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Is(t *testing.T) {
	type fields struct {
		Message []string
		Code    int
	}
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "nil error",
			fields: fields{[]string{"not found"}, 404},
			args:   args{nil},
			want:   false,
		},
		{
			name:   "same error code",
			fields: fields{[]string{"not found"}, 404},
			args:   args{New(404, "resource not found")},
			want:   true,
		},
		{
			name:   "different error code",
			fields: fields{[]string{"not found"}, 404},
			args:   args{New(500, "internal error")},
			want:   false,
		},
		{
			name:   "non-Error type",
			fields: fields{[]string{"not found"}, 404},
			args:   args{errors.New("standard error")},
			want:   false,
		},
		{
			name:   "same error instance",
			fields: fields{[]string{"test1"}, 200},
			args:   args{New(200, "test")},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New(tt.fields.Code, tt.fields.Message...)
			if got := e.Is(tt.args.err); got != tt.want {
				t.Errorf("Error.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Wrap(t *testing.T) {
	originalErr := errors.New("original error")
	xecErr := New(500, "server error")

	wrappedErr := xecErr.Wrap(originalErr)

	if wrappedErr.Code != 500 {
		t.Errorf("Expected code 500, got %d", wrappedErr.Code)
	}

	if wrappedErr.Message != "[server error]" {
		t.Errorf("Expected message '[server error]', got %s", wrappedErr.Message)
	}

	if wrappedErr.cause != originalErr {
		t.Errorf("Expected cause to be original error")
	}
}

func TestError_ErrorWithCause(t *testing.T) {
	originalErr := errors.New("database connection failed")
	xecErr := New(500, "internal server error")
	wrappedErr := xecErr.Wrap(originalErr)

	errorMsg := wrappedErr.Error()
	expected := "[500] [internal server error] (caused by: database connection failed)"

	if errorMsg != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, errorMsg)
	}
}

func TestError_ErrorWithoutCause(t *testing.T) {
	xecErr := New(404, "not found")

	errorMsg := xecErr.Error()
	expected := "[404] [not found]"

	if errorMsg != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, errorMsg)
	}
}

func TestError_IsWithWrappedError(t *testing.T) {
	originalErr := errors.New("original error")
	xecErr := New(500, "server error")
	wrappedErr := xecErr.Wrap(originalErr)

	if !xecErr.Is(wrappedErr) {
		t.Errorf("Expected Is() to return true for errors with same code")
	}

	differentErr := New(404, "not found")
	if differentErr.Is(wrappedErr) {
		t.Errorf("Expected Is() to return false for errors with different code")
	}
}

func TestError_CauseChaining(t *testing.T) {
	rootCause := errors.New("root cause")
	firstWrap := New(500, "first error").Wrap(rootCause)
	secondWrap := New(500, "second error").Wrap(firstWrap)

	errorMsg := secondWrap.Error()

	if !strings.Contains(errorMsg, "root cause") {
		t.Errorf("Expected error message to contain root cause, got: %s", errorMsg)
	}

	if !strings.Contains(errorMsg, "first error") {
		t.Errorf("Expected error message to contain first error, got: %s", errorMsg)
	}
}

func TestError_ErrorWithMultiLevelCause(t *testing.T) {
	rootCause := errors.New("database connection timeout")
	firstLevel := New(500, "database error").Wrap(rootCause)
	secondLevel := New(503, "service unavailable").Wrap(firstLevel)
	thirdLevel := New(500, "internal server error").Wrap(secondLevel)

	errorMsg := thirdLevel.Error()
	expected := "[500] [internal server error] (caused by: [503] [service unavailable] (caused by: [500] [database error] (caused by: database connection timeout)))"

	if errorMsg != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, errorMsg)
	}
	t.Log(errorMsg)
	// Verify all levels are present in the output
	if !strings.Contains(errorMsg, "internal server error") {
		t.Errorf("Expected error message to contain 'internal server error', got: %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "service unavailable") {
		t.Errorf("Expected error message to contain 'service unavailable', got: %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "database error") {
		t.Errorf("Expected error message to contain 'database error', got: %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "database connection timeout") {
		t.Errorf("Expected error message to contain 'database connection timeout', got: %s", errorMsg)
	}
}
