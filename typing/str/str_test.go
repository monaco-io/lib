package str

import (
	"reflect"
	"strings"
	"testing"
)

func TestPick(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want string
	}{
		{"first not empty", "hello", "world", "hello"},
		{"first empty", "", "world", "world"},
		{"both empty", "", "", ""},
		{"both not empty", "hello", "world", "hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Pick(tt.a, tt.b); got != tt.want {
				t.Errorf("Pick() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUUID(t *testing.T) {
	got := UUID()
	if len(got) != 36 {
		t.Errorf("UUID() length = %v, want 36", len(got))
	}
	if strings.Count(got, "-") != 4 {
		t.Errorf("UUID() hyphens count = %v, want 4", strings.Count(got, "-"))
	}
}

func TestUUIDX(t *testing.T) {
	got := UUIDX()
	if len(got) != 32 {
		t.Errorf("UUIDX() length = %v, want 32", len(got))
	}
	if strings.Contains(got, "-") {
		t.Errorf("UUIDX() should not contain hyphens")
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"empty string", "", true},
		{"non-empty string", "hello", false},
		{"space", " ", false},
		{"tab", "\t", false},
		{"newline", "\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmpty(tt.s); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNotEmpty(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"empty string", "", false},
		{"non-empty string", "hello", true},
		{"space", " ", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotEmpty(tt.s); got != tt.want {
				t.Errorf("IsNotEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsBlank(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"empty string", "", true},
		{"space", " ", true},
		{"tab", "\t", true},
		{"newline", "\n", true},
		{"multiple spaces", "   ", true},
		{"mixed whitespace", " \t\n ", true},
		{"non-blank string", "hello", false},
		{"string with spaces", " hello ", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsBlank(tt.s); got != tt.want {
				t.Errorf("IsBlank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNotBlank(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"empty string", "", false},
		{"space", " ", false},
		{"non-blank string", "hello", true},
		{"string with spaces", " hello ", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotBlank(tt.s); got != tt.want {
				t.Errorf("IsNotBlank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultIfEmpty(t *testing.T) {
	tests := []struct {
		name         string
		s            string
		defaultValue string
		want         string
	}{
		{"empty returns default", "", "default", "default"},
		{"non-empty returns original", "hello", "default", "hello"},
		{"space returns space", " ", "default", " "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultIfEmpty(tt.s, tt.defaultValue); got != tt.want {
				t.Errorf("DefaultIfEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultIfBlank(t *testing.T) {
	tests := []struct {
		name         string
		s            string
		defaultValue string
		want         string
	}{
		{"empty returns default", "", "default", "default"},
		{"blank returns default", " ", "default", "default"},
		{"tab returns default", "\t", "default", "default"},
		{"non-blank returns original", "hello", "default", "hello"},
		{"string with spaces returns original", " hello ", "default", " hello "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultIfBlank(tt.s, tt.defaultValue); got != tt.want {
				t.Errorf("DefaultIfBlank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"empty string", "", ""},
		{"single char", "a", "a"},
		{"simple string", "hello", "olleh"},
		{"unicode", "你好", "好你"},
		{"mixed", "a你b好c", "c好b你a"},
		{"palindrome", "aba", "aba"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reverse(tt.s); got != tt.want {
				t.Errorf("Reverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsIgnoreCase(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		substr string
		want   bool
	}{
		{"exact match", "hello", "hello", true},
		{"case insensitive", "Hello", "hello", true},
		{"uppercase", "HELLO", "hello", true},
		{"mixed case", "HeLLo", "ell", true},
		{"not found", "hello", "world", false},
		{"empty substr", "hello", "", true},
		{"empty string", "", "hello", false},
		{"both empty", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsIgnoreCase(tt.s, tt.substr); got != tt.want {
				t.Errorf("ContainsIgnoreCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStartsWithIgnoreCase(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		prefix string
		want   bool
	}{
		{"exact match", "hello", "hello", true},
		{"case insensitive", "Hello", "hello", true},
		{"prefix match", "Hello World", "hello", true},
		{"uppercase", "HELLO", "hello", true},
		{"not match", "hello", "world", false},
		{"empty prefix", "hello", "", true},
		{"empty string", "", "hello", false},
		{"both empty", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StartsWithIgnoreCase(tt.s, tt.prefix); got != tt.want {
				t.Errorf("StartsWithIgnoreCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndsWithIgnoreCase(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		suffix string
		want   bool
	}{
		{"exact match", "hello", "hello", true},
		{"case insensitive", "Hello", "hello", true},
		{"suffix match", "Hello World", "world", true},
		{"uppercase", "HELLO", "hello", true},
		{"not match", "hello", "world", false},
		{"empty suffix", "hello", "", true},
		{"empty string", "", "hello", false},
		{"both empty", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EndsWithIgnoreCase(tt.s, tt.suffix); got != tt.want {
				t.Errorf("EndsWithIgnoreCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		length int
		want   string
	}{
		{"truncate needed", "hello world", 5, "hello"},
		{"no truncate needed", "hello", 10, "hello"},
		{"exact length", "hello", 5, "hello"},
		{"zero length", "hello", 0, ""},
		{"negative length", "hello", -1, ""},
		{"empty string", "", 5, ""},
		{"unicode", "你好世界", 2, "你好"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Truncate(tt.s, tt.length); got != tt.want {
				t.Errorf("Truncate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollapseWhitespace(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"multiple spaces", "hello    world", "hello world"},
		{"tabs and spaces", "hello\t  \tworld", "hello world"},
		{"newlines", "hello\n\nworld", "hello world"},
		{"mixed whitespace", " \t hello \n\n world \t ", "hello world"},
		{"single space", "hello world", "hello world"},
		{"no whitespace", "helloworld", "helloworld"},
		{"only whitespace", "   \t\n  ", ""},
		{"empty string", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CollapseWhitespace(tt.s); got != tt.want {
				t.Errorf("CollapseWhitespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"camelCase", "helloWorld", "hello_world"},
		{"PascalCase", "HelloWorld", "hello_world"},
		{"multiple words", "HelloWorldFoo", "hello_world_foo"},
		{"single word", "hello", "hello"},
		{"uppercase", "HELLO", "hello"},
		{"with numbers", "hello2World", "hello2_world"},
		{"already snake", "hello_world", "hello_world"},
		{"empty string", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSnakeCase(tt.s); got != tt.want {
				t.Errorf("ToSnakeCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"snake_case", "hello_world", "helloWorld"},
		{"kebab-case", "hello-world", "helloWorld"},
		{"space separated", "hello world", "helloWorld"},
		{"multiple separators", "hello_world-foo bar", "helloWorldFooBar"},
		{"single word", "hello", "hello"},
		{"PascalCase", "HelloWorld", "helloworld"},
		{"mixed", "hello_World-FOO", "helloWorldFoo"},
		{"empty string", "", ""},
		{"numbers", "hello2_world3", "hello2World3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToCamelCase(tt.s); got != tt.want {
				t.Errorf("ToCamelCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"snake_case", "hello_world", "HelloWorld"},
		{"kebab-case", "hello-world", "HelloWorld"},
		{"space separated", "hello world", "HelloWorld"},
		{"camelCase", "helloWorld", "Helloworld"},
		{"single word", "hello", "Hello"},
		{"multiple separators", "hello_world-foo bar", "HelloWorldFooBar"},
		{"empty string", "", ""},
		{"numbers", "hello2_world3", "Hello2World3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToPascalCase(tt.s); got != tt.want {
				t.Errorf("ToPascalCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToKebabCase(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"camelCase", "helloWorld", "hello-world"},
		{"PascalCase", "HelloWorld", "hello-world"},
		{"snake_case", "hello_world", "hello-world"},
		{"space separated", "hello world", "hello-world"},
		{"multiple separators", "hello_world FOO-bar", "hello-world-foo-bar"},
		{"single word", "hello", "hello"},
		{"empty string", "", ""},
		{"numbers", "hello2World3", "hello2-world3"},
		{"special chars", "hello@world#foo", "hello-world-foo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToKebabCase(tt.s); got != tt.want {
				t.Errorf("ToKebabCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSwapCase(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"mixed case", "Hello World", "hELLO wORLD"},
		{"all lowercase", "hello", "HELLO"},
		{"all uppercase", "HELLO", "hello"},
		{"with numbers", "Hello123", "hELLO123"},
		{"with symbols", "Hello, World!", "hELLO, wORLD!"},
		{"empty string", "", ""},
		{"unicode", "你好Hello", "你好hELLO"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SwapCase(tt.s); got != tt.want {
				t.Errorf("SwapCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountWords(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{"single word", "hello", 1},
		{"two words", "hello world", 2},
		{"multiple spaces", "hello    world", 2},
		{"with tabs", "hello\tworld", 2},
		{"with newlines", "hello\nworld", 2},
		{"mixed whitespace", " hello \t world \n ", 2},
		{"empty string", "", 0},
		{"only spaces", "   ", 0},
		{"three words", "hello beautiful world", 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CountWords(tt.s); got != tt.want {
				t.Errorf("CountWords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountLines(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{"single line", "hello", 1},
		{"two lines", "hello\nworld", 2},
		{"three lines", "hello\nworld\nfoo", 3},
		{"empty string", "", 0},
		{"only newline", "\n", 2},
		{"trailing newline", "hello\n", 2},
		{"multiple newlines", "hello\n\nworld", 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CountLines(tt.s); got != tt.want {
				t.Errorf("CountLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepeat(t *testing.T) {
	tests := []struct {
		name  string
		s     string
		count int
		want  string
	}{
		{"repeat hello", "hello", 3, "hellohellohello"},
		{"repeat single char", "a", 5, "aaaaa"},
		{"repeat zero", "hello", 0, ""},
		{"repeat negative", "hello", -1, ""},
		{"repeat empty", "", 3, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Repeat(tt.s, tt.count); got != tt.want {
				t.Errorf("Repeat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndexOfIgnoreCase(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		substr string
		want   int
	}{
		{"exact match", "hello", "hello", 0},
		{"case insensitive", "Hello", "hello", 0},
		{"found in middle", "Hello World", "world", 6},
		{"uppercase", "HELLO WORLD", "world", 6},
		{"not found", "hello", "world", -1},
		{"empty substr", "hello", "", 0},
		{"empty string", "", "hello", -1},
		{"both empty", "", "", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IndexOfIgnoreCase(tt.s, tt.substr); got != tt.want {
				t.Errorf("IndexOfIgnoreCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLastIndexOfIgnoreCase(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		substr string
		want   int
	}{
		{"exact match", "hello", "hello", 0},
		{"case insensitive", "Hello", "hello", 0},
		{"multiple occurrences", "Hello World Hello", "hello", 12},
		{"uppercase", "HELLO WORLD HELLO", "hello", 12},
		{"not found", "hello", "world", -1},
		{"empty substr", "hello", "", 5},
		{"empty string", "", "hello", -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LastIndexOfIgnoreCase(tt.s, tt.substr); got != tt.want {
				t.Errorf("LastIndexOfIgnoreCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplitAndTrim(t *testing.T) {
	tests := []struct {
		name string
		s    string
		sep  string
		want []string
	}{
		{"simple split", "hello,world", ",", []string{"hello", "world"}},
		{"with spaces", " hello , world ", ",", []string{"hello", "world"}},
		{"empty parts", "hello,,world", ",", []string{"hello", "world"}},
		{"only spaces", " , , ", ",", []string{}},
		{"no separator", "hello", ",", []string{"hello"}},
		{"empty string", "", ",", []string{}},
		{"trailing separator", "hello,world,", ",", []string{"hello", "world"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitAndTrim(tt.s, tt.sep); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitAndTrim() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInteger(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    int
		wantErr bool
	}{
		{"valid integer", "123", 123, false},
		{"negative integer", "-123", -123, false},
		{"with spaces", " 123 ", 123, false},
		{"zero", "0", 0, false},
		{"invalid", "abc", 0, true},
		{"empty", "", 0, true},
		{"float", "123.45", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToInteger[int](tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToInteger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToInteger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToFloat(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    float64
		wantErr bool
	}{
		{"valid float", "123.45", 123.45, false},
		{"integer", "123", 123.0, false},
		{"negative", "-123.45", -123.45, false},
		{"with spaces", " 123.45 ", 123.45, false},
		{"zero", "0", 0.0, false},
		{"invalid", "abc", 0.0, true},
		{"empty", "", 0.0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToFloat[float64](tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToBool(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    bool
		wantErr bool
	}{
		{"true", "true", true, false},
		{"false", "false", false, false},
		{"True", "True", true, false},
		{"FALSE", "FALSE", false, false},
		{"1", "1", true, false},
		{"0", "0", false, false},
		{"with spaces", " true ", true, false},
		{"invalid", "yes", false, true},
		{"empty", "", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToBool(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromInteger(t *testing.T) {
	tests := []struct {
		name string
		i    int
		want string
	}{
		{"positive", 123, "123"},
		{"negative", -123, "-123"},
		{"zero", 0, "0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromInteger(tt.i); got != tt.want {
				t.Errorf("FromInteger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromFloat(t *testing.T) {
	tests := []struct {
		name      string
		f         float64
		precision int
		want      string
	}{
		{"simple float", 123.45, 2, "123.45"},
		{"no decimal", 123.0, 2, "123.00"},
		{"high precision", 123.456789, 4, "123.4568"},
		{"zero precision", 123.456, 0, "123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromFloat(tt.f, tt.precision); got != tt.want {
				t.Errorf("FromFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromBool(t *testing.T) {
	tests := []struct {
		name string
		b    bool
		want string
	}{
		{"true", true, "true"},
		{"false", false, "false"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromBool(tt.b); got != tt.want {
				t.Errorf("FromBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMD5(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"hello", "hello", "5d41402abc4b2a76b9719d911017c592"},
		{"empty", "", "d41d8cd98f00b204e9800998ecf8427e"},
		{"unicode", "你好", "7eca689f0d3389d9dea66ae112e5cfd7"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MD5(tt.s); got != tt.want {
				t.Errorf("MD5() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSHA256(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"hello", "hello", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
		{"empty", "", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SHA256(tt.s); got != tt.want {
				t.Errorf("SHA256() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMask(t *testing.T) {
	tests := []struct {
		name  string
		s     string
		start int
		end   int
		mask  string
		want  string
	}{
		{"basic mask", "1234567890", 3, 3, "*", "123****890"},
		{"custom mask", "1234567890", 2, 2, "x", "12xxxxxx90"},
		{"empty mask", "1234567890", 3, 3, "", "123****890"},
		{"no mask needed", "123", 1, 1, "*", "1*3"},
		{"negative values", "1234567890", -1, -1, "*", "**********"},
		{"start+end >= length", "123", 2, 2, "*", "123"},
		{"unicode", "你好世界", 1, 1, "*", "你**界"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Mask(tt.s, tt.start, tt.end, tt.mask); got != tt.want {
				t.Errorf("Mask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandom(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		charset string
	}{
		{"default charset", 10, ""},
		{"custom charset", 5, "abc"},
		{"zero length", 0, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Random(tt.length, tt.charset)
			if len(got) != tt.length {
				t.Errorf("Random() length = %v, want %v", len(got), tt.length)
			}
		})
	}
}

func TestSimilarity(t *testing.T) {
	tests := []struct {
		name string
		s1   string
		s2   string
		want float64
	}{
		{"identical", "hello", "hello", 1.0},
		{"empty strings", "", "", 1.0},
		{"one empty", "hello", "", 0.0},
		{"completely different", "abc", "xyz", 0.0},
		{"similar", "hello", "hallo", 0.8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Similarity(tt.s1, tt.s2)
			if got != tt.want {
				t.Errorf("Similarity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		name   string
		format string
		args   []interface{}
		want   string
	}{
		{"string format", "Hello %s", []interface{}{"world"}, "Hello world"},
		{"integer format", "Number: %d", []interface{}{42}, "Number: 42"},
		{"multiple args", "%s: %d", []interface{}{"Count", 42}, "Count: 42"},
		{"no args", "Hello", []interface{}{}, "Hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Format(tt.format, tt.args...); got != tt.want {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Benchmarks
func BenchmarkPick(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Pick("hello", "world")
	}
}

func BenchmarkUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UUID()
	}
}

func BenchmarkUUIDX(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UUIDX()
	}
}

func BenchmarkToSnakeCase(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ToSnakeCase("HelloWorldFoo")
	}
}

func BenchmarkToCamelCase(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ToCamelCase("hello_world_foo")
	}
}

func BenchmarkMD5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MD5("hello world")
	}
}

func BenchmarkSHA256(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SHA256("hello world")
	}
}

func BenchmarkSimilarity(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Similarity("hello world", "hello world!")
	}
}
