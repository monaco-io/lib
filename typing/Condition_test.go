package typing_test

import (
	"testing"

	"github.com/monaco-io/lib/typing"
)

func TestIsEmpty(t *testing.T) {
	var empty string
	var emptyPtr *string
	noEmpty := "1"
	var noEmptyPtr = &noEmpty
	// nolint:staticcheck
	tests := []struct {
		// Named input parameters for target function.
		name string
		v    any
		want bool
	}{
		{"case-1", "", true},
		{"case-2", nil, true},
		{"case-3", empty, true},
		{"case-4", emptyPtr, true},
		{"case-5", &empty, true},
		{"case-6", &emptyPtr, true},
		{"case-7", noEmptyPtr, false},
		{"case-8", &noEmptyPtr, false},
		{"case-9", []string{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := typing.IsEmpty(tt.v)
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("IsEmpty('%v') = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestIsEmptyExtended(t *testing.T) {
	// 测试基本类型
	t.Run("basic types", func(t *testing.T) {
		if !typing.IsEmpty(0) {
			t.Error("Expected 0 to be empty")
		}
		if !typing.IsEmpty("") {
			t.Error("Expected empty string to be empty")
		}
		if !typing.IsEmpty(false) {
			t.Error("Expected false to be empty")
		}
		if typing.IsEmpty(1) {
			t.Error("Expected 1 not to be empty")
		}
		if typing.IsEmpty("hello") {
			t.Error("Expected 'hello' not to be empty")
		}
		if typing.IsEmpty(true) {
			t.Error("Expected true not to be empty")
		}
		if !typing.IsEmpty(0.0) {
			t.Error("Expected 0.0 to be empty")
		}
		if typing.IsEmpty(1.5) {
			t.Error("Expected 1.5 not to be empty")
		}
		if !typing.IsEmpty(complex(0, 0)) {
			t.Error("Expected 0+0i to be empty")
		}
		if typing.IsEmpty(complex(1, 2)) {
			t.Error("Expected 1+2i not to be empty")
		}
	})

	// 测试多级指针
	t.Run("multi-level pointers", func(t *testing.T) {
		var nilPtr *int
		if !typing.IsEmpty(nilPtr) {
			t.Error("Expected nil pointer to be empty")
		}

		val := 42
		ptr := &val
		if typing.IsEmpty(ptr) {
			t.Error("Expected non-nil pointer to not be empty")
		}

		zeroVal := 0
		zeroPtr := &zeroVal
		if !typing.IsEmpty(zeroPtr) {
			t.Error("Expected pointer to zero value to be empty")
		}

		// 双重指针
		ptrPtr := &ptr
		if typing.IsEmpty(ptrPtr) {
			t.Error("Expected pointer to pointer (non-zero) to not be empty")
		}

		nilPtrPtr := &nilPtr
		if !typing.IsEmpty(nilPtrPtr) {
			t.Error("Expected pointer to nil pointer to be empty")
		}
	})

	// 测试切片和数组
	t.Run("slices and arrays", func(t *testing.T) {
		var nilSlice []int
		if !typing.IsEmpty(nilSlice) {
			t.Error("Expected nil slice to be empty")
		}

		emptySlice := make([]int, 0)
		if !typing.IsEmpty(emptySlice) {
			t.Error("Expected empty slice to be empty")
		}

		nonEmptySlice := []int{1}
		if typing.IsEmpty(nonEmptySlice) {
			t.Error("Expected non-empty slice to not be empty")
		}

		// 测试数组
		var emptyArray [0]int
		if !typing.IsEmpty(emptyArray) {
			t.Error("Expected empty array to be empty")
		}

		nonEmptyArray := [1]int{1}
		if typing.IsEmpty(nonEmptyArray) {
			t.Error("Expected non-empty array to not be empty")
		}
	})

	// 测试映射
	t.Run("maps", func(t *testing.T) {
		var nilMap map[string]int
		if !typing.IsEmpty(nilMap) {
			t.Error("Expected nil map to be empty")
		}

		emptyMap := make(map[string]int)
		if !typing.IsEmpty(emptyMap) {
			t.Error("Expected empty map to be empty")
		}

		nonEmptyMap := map[string]int{"key": 1}
		if typing.IsEmpty(nonEmptyMap) {
			t.Error("Expected non-empty map to not be empty")
		}
	})

	// 测试结构体
	t.Run("structs", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		var zeroPerson Person
		if !typing.IsEmpty(zeroPerson) {
			t.Error("Expected zero value struct to be empty")
		}

		nonZeroPerson := Person{Name: "Alice", Age: 30}
		if typing.IsEmpty(nonZeroPerson) {
			t.Error("Expected non-zero struct to not be empty")
		}

		partialPerson := Person{Name: "Bob"}
		if typing.IsEmpty(partialPerson) {
			t.Error("Expected partially filled struct to not be empty")
		}
	})

	// 测试通道
	t.Run("channels", func(t *testing.T) {
		var nilChan chan int
		if !typing.IsEmpty(nilChan) {
			t.Error("Expected nil channel to be empty")
		}

		ch := make(chan int)
		// 非 nil 的通道不应该被认为是空的，即使没有数据
		if typing.IsEmpty(ch) {
			t.Error("Expected non-nil channel to not be empty")
		}
		close(ch)

		// 测试带缓冲的通道
		bufferedCh := make(chan int, 1)
		// 非 nil 的通道不应该被认为是空的，即使缓冲区是空的
		if typing.IsEmpty(bufferedCh) {
			t.Error("Expected non-nil buffered channel to not be empty")
		}
		close(bufferedCh)
	}) // 测试函数
	t.Run("functions", func(t *testing.T) {
		var nilFunc func()
		if !typing.IsEmpty(nilFunc) {
			t.Error("Expected nil function to be empty")
		}

		nonNilFunc := func() {}
		if typing.IsEmpty(nonNilFunc) {
			t.Error("Expected non-nil function to not be empty")
		}
	})

	// 测试接口
	t.Run("interfaces", func(t *testing.T) {
		var nilInterface interface{}
		if !typing.IsEmpty(nilInterface) {
			t.Error("Expected nil interface to be empty")
		}

		var nonNilInterface interface{} = "hello"
		if typing.IsEmpty(nonNilInterface) {
			t.Error("Expected non-nil interface to not be empty")
		}

		var zeroValueInterface interface{} = 0
		if !typing.IsEmpty(zeroValueInterface) {
			t.Error("Expected interface with zero value to be empty")
		}
	})
}

func TestDefaultIfEmpty(t *testing.T) {
	// 测试基本类型
	t.Run("basic types", func(t *testing.T) {
		// 字符串
		if got := typing.DefaultIfEmpty("", "default"); got != "default" {
			t.Errorf("DefaultIfEmpty('', 'default') = %v, want 'default'", got)
		}
		if got := typing.DefaultIfEmpty("hello", "default"); got != "hello" {
			t.Errorf("DefaultIfEmpty('hello', 'default') = %v, want 'hello'", got)
		}

		// 整数
		if got := typing.DefaultIfEmpty(0, 42); got != 42 {
			t.Errorf("DefaultIfEmpty(0, 42) = %v, want 42", got)
		}
		if got := typing.DefaultIfEmpty(10, 42); got != 10 {
			t.Errorf("DefaultIfEmpty(10, 42) = %v, want 10", got)
		}

		// 布尔值
		if got := typing.DefaultIfEmpty(false, true); got != true {
			t.Errorf("DefaultIfEmpty(false, true) = %v, want true", got)
		}
		if got := typing.DefaultIfEmpty(true, false); got != true {
			t.Errorf("DefaultIfEmpty(true, false) = %v, want true", got)
		}

		// 浮点数
		if got := typing.DefaultIfEmpty(0.0, 3.14); got != 3.14 {
			t.Errorf("DefaultIfEmpty(0.0, 3.14) = %v, want 3.14", got)
		}
		if got := typing.DefaultIfEmpty(2.5, 3.14); got != 2.5 {
			t.Errorf("DefaultIfEmpty(2.5, 3.14) = %v, want 2.5", got)
		}
	})

	// 测试切片
	t.Run("slices", func(t *testing.T) {
		defaultSlice := []int{1, 2, 3}

		// nil 切片
		var nilSlice []int
		if got := typing.DefaultIfEmpty(nilSlice, defaultSlice); len(got) != 3 {
			t.Errorf("DefaultIfEmpty(nil slice, default) should return default slice")
		}

		// 空切片
		emptySlice := make([]int, 0)
		if got := typing.DefaultIfEmpty(emptySlice, defaultSlice); len(got) != 3 {
			t.Errorf("DefaultIfEmpty(empty slice, default) should return default slice")
		}

		// 非空切片
		nonEmptySlice := []int{4, 5}
		if got := typing.DefaultIfEmpty(nonEmptySlice, defaultSlice); len(got) != 2 {
			t.Errorf("DefaultIfEmpty(non-empty slice, default) should return original slice")
		}
	})

	// 测试映射
	t.Run("maps", func(t *testing.T) {
		defaultMap := map[string]int{"key": 1}

		// nil 映射
		var nilMap map[string]int
		if got := typing.DefaultIfEmpty(nilMap, defaultMap); len(got) != 1 {
			t.Errorf("DefaultIfEmpty(nil map, default) should return default map")
		}

		// 空映射
		emptyMap := make(map[string]int)
		if got := typing.DefaultIfEmpty(emptyMap, defaultMap); len(got) != 1 {
			t.Errorf("DefaultIfEmpty(empty map, default) should return default map")
		}

		// 非空映射
		nonEmptyMap := map[string]int{"other": 2}
		if got := typing.DefaultIfEmpty(nonEmptyMap, defaultMap); len(got) != 1 || got["other"] != 2 {
			t.Errorf("DefaultIfEmpty(non-empty map, default) should return original map")
		}
	})

	// 测试结构体
	t.Run("structs", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		defaultPerson := Person{Name: "Default", Age: 0}

		// 零值结构体
		var zeroPerson Person
		if got := typing.DefaultIfEmpty(zeroPerson, defaultPerson); got.Name != "Default" {
			t.Errorf("DefaultIfEmpty(zero struct, default) should return default struct")
		}

		// 非零结构体
		nonZeroPerson := Person{Name: "Alice", Age: 30}
		if got := typing.DefaultIfEmpty(nonZeroPerson, defaultPerson); got.Name != "Alice" {
			t.Errorf("DefaultIfEmpty(non-zero struct, default) should return original struct")
		}
	})

	// 测试指针
	t.Run("pointers", func(t *testing.T) {
		defaultVal := 42
		defaultPtr := &defaultVal

		// nil 指针
		var nilPtr *int
		if got := typing.DefaultIfEmpty(nilPtr, defaultPtr); got != defaultPtr {
			t.Errorf("DefaultIfEmpty(nil pointer, default) should return default pointer")
		}

		// 非 nil 指针
		val := 10
		ptr := &val
		if got := typing.DefaultIfEmpty(ptr, defaultPtr); got != ptr {
			t.Errorf("DefaultIfEmpty(non-nil pointer, default) should return original pointer")
		}

		// 指向零值的指针
		zeroVal := 0
		zeroPtr := &zeroVal
		if got := typing.DefaultIfEmpty(zeroPtr, defaultPtr); got != defaultPtr {
			t.Errorf("DefaultIfEmpty(pointer to zero, default) should return default pointer")
		}
	})
}

func TestDefaultIfEmptyPtr(t *testing.T) {
	// 测试基本类型
	t.Run("basic types", func(t *testing.T) {
		// 字符串指针
		defaultStr := "default"

		// nil 指针，有默认值
		var nilStrPtr *string
		if got := typing.DefaultIfEmptyPtr(nilStrPtr, defaultStr); got != defaultStr {
			t.Errorf("DefaultIfEmptyPtr(nil, 'default') = %v, want 'default'", got)
		}

		// nil 指针，无默认值
		if got := typing.DefaultIfEmptyPtr(nilStrPtr); got != "" {
			t.Errorf("DefaultIfEmptyPtr(nil) = %v, want zero value", got)
		}

		// 指向空值的指针
		emptyStr := ""
		emptyStrPtr := &emptyStr
		if got := typing.DefaultIfEmptyPtr(emptyStrPtr, defaultStr); got != defaultStr {
			t.Errorf("DefaultIfEmptyPtr(empty string ptr, 'default') = %v, want 'default'", got)
		}

		// 指向非空值的指针
		nonEmptyStr := "hello"
		nonEmptyStrPtr := &nonEmptyStr
		if got := typing.DefaultIfEmptyPtr(nonEmptyStrPtr, defaultStr); got != "hello" {
			t.Errorf("DefaultIfEmptyPtr('hello' ptr, 'default') = %v, want 'hello'", got)
		}
	})

	// 测试整数指针
	t.Run("integer pointers", func(t *testing.T) {
		defaultInt := 42

		// nil 指针
		var nilIntPtr *int
		if got := typing.DefaultIfEmptyPtr(nilIntPtr, defaultInt); got != defaultInt {
			t.Errorf("DefaultIfEmptyPtr(nil int ptr, 42) = %v, want 42", got)
		}

		// 指向零值的指针
		zeroInt := 0
		zeroIntPtr := &zeroInt
		if got := typing.DefaultIfEmptyPtr(zeroIntPtr, defaultInt); got != defaultInt {
			t.Errorf("DefaultIfEmptyPtr(zero int ptr, 42) = %v, want 42", got)
		}

		// 指向非零值的指针
		nonZeroInt := 10
		nonZeroIntPtr := &nonZeroInt
		if got := typing.DefaultIfEmptyPtr(nonZeroIntPtr, defaultInt); got != 10 {
			t.Errorf("DefaultIfEmptyPtr(10 ptr, 42) = %v, want 10", got)
		}
	})

	// 测试切片指针
	t.Run("slice pointers", func(t *testing.T) {
		defaultSlice := []int{1, 2, 3}

		// nil 指针
		var nilSlicePtr *[]int
		if got := typing.DefaultIfEmptyPtr(nilSlicePtr, defaultSlice); len(got) != 3 {
			t.Errorf("DefaultIfEmptyPtr(nil slice ptr) should return default slice")
		}

		// 指向 nil 切片的指针
		var nilSlice []int
		nilSlicePtr2 := &nilSlice
		if got := typing.DefaultIfEmptyPtr(nilSlicePtr2, defaultSlice); len(got) != 3 {
			t.Errorf("DefaultIfEmptyPtr(nil slice ptr) should return default slice")
		}

		// 指向空切片的指针
		emptySlice := make([]int, 0)
		emptySlicePtr := &emptySlice
		if got := typing.DefaultIfEmptyPtr(emptySlicePtr, defaultSlice); len(got) != 3 {
			t.Errorf("DefaultIfEmptyPtr(empty slice ptr) should return default slice")
		}

		// 指向非空切片的指针
		nonEmptySlice := []int{4, 5}
		nonEmptySlicePtr := &nonEmptySlice
		if got := typing.DefaultIfEmptyPtr(nonEmptySlicePtr, defaultSlice); len(got) != 2 {
			t.Errorf("DefaultIfEmptyPtr(non-empty slice ptr) should return original slice")
		}
	})

	// 测试结构体指针
	t.Run("struct pointers", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		defaultPerson := Person{Name: "Default", Age: 0}

		// nil 指针
		var nilPersonPtr *Person
		if got := typing.DefaultIfEmptyPtr(nilPersonPtr, defaultPerson); got.Name != "Default" {
			t.Errorf("DefaultIfEmptyPtr(nil person ptr) should return default person")
		}

		// 指向零值结构体的指针
		zeroPerson := Person{}
		zeroPersonPtr := &zeroPerson
		if got := typing.DefaultIfEmptyPtr(zeroPersonPtr, defaultPerson); got.Name != "Default" {
			t.Errorf("DefaultIfEmptyPtr(zero person ptr) should return default person")
		}

		// 指向非零结构体的指针
		nonZeroPerson := Person{Name: "Alice", Age: 30}
		nonZeroPersonPtr := &nonZeroPerson
		if got := typing.DefaultIfEmptyPtr(nonZeroPersonPtr, defaultPerson); got.Name != "Alice" {
			t.Errorf("DefaultIfEmptyPtr(non-zero person ptr) should return original person")
		}
	})

	// 测试多个默认值参数
	t.Run("multiple default values", func(t *testing.T) {
		var nilStrPtr *string

		// 传递多个默认值，应该使用第一个
		if got := typing.DefaultIfEmptyPtr(nilStrPtr, "first", "second", "third"); got != "first" {
			t.Errorf("DefaultIfEmptyPtr should use first default value, got %v", got)
		}
	})

	// 测试无默认值参数
	t.Run("no default values", func(t *testing.T) {
		var nilStrPtr *string

		// 不传递默认值，应该返回零值
		if got := typing.DefaultIfEmptyPtr(nilStrPtr); got != "" {
			t.Errorf("DefaultIfEmptyPtr with no defaults should return zero value, got %v", got)
		}

		var nilIntPtr *int
		if got := typing.DefaultIfEmptyPtr(nilIntPtr); got != 0 {
			t.Errorf("DefaultIfEmptyPtr with no defaults should return zero value, got %v", got)
		}
	})

	// 测试嵌套指针
	t.Run("nested pointers", func(t *testing.T) {
		defaultStr := "default"

		// 指向 nil 指针的指针
		var nilStrPtr *string
		nilStrPtrPtr := &nilStrPtr
		if got := typing.DefaultIfEmptyPtr(nilStrPtrPtr, &defaultStr); got == nil || *got != "default" {
			t.Errorf("DefaultIfEmptyPtr(nested nil ptr) should return default pointer")
		}

		// 指向指向空值指针的指针
		emptyStr := ""
		emptyStrPtr := &emptyStr
		emptyStrPtrPtr := &emptyStrPtr
		if got := typing.DefaultIfEmptyPtr(emptyStrPtrPtr, &defaultStr); got == nil || *got != "default" {
			t.Errorf("DefaultIfEmptyPtr(nested empty ptr) should return default pointer")
		}
	})

	// 测试边界情况
	t.Run("edge cases", func(t *testing.T) {
		// 测试不同类型的零值
		var nilBoolPtr *bool
		if got := typing.DefaultIfEmptyPtr(nilBoolPtr, true); got != true {
			t.Errorf("DefaultIfEmptyPtr(nil bool ptr, true) = %v, want true", got)
		}

		falseBool := false
		falseBoolPtr := &falseBool
		if got := typing.DefaultIfEmptyPtr(falseBoolPtr, true); got != true {
			t.Errorf("DefaultIfEmptyPtr(false ptr, true) = %v, want true", got)
		}

		// 测试浮点数
		var nilFloatPtr *float64
		if got := typing.DefaultIfEmptyPtr(nilFloatPtr, 3.14); got != 3.14 {
			t.Errorf("DefaultIfEmptyPtr(nil float ptr, 3.14) = %v, want 3.14", got)
		}

		zeroFloat := 0.0
		zeroFloatPtr := &zeroFloat
		if got := typing.DefaultIfEmptyPtr(zeroFloatPtr, 3.14); got != 3.14 {
			t.Errorf("DefaultIfEmptyPtr(0.0 ptr, 3.14) = %v, want 3.14", got)
		}
	})
}

// TestCondition 测试三元运算符函数
func TestCondition(t *testing.T) {
	// 测试基本用法
	t.Run("basic usage", func(t *testing.T) {
		// 条件为真
		if got := typing.Condition(true, "yes", "no"); got != "yes" {
			t.Errorf("Condition(true, 'yes', 'no') = %v, want 'yes'", got)
		}

		// 条件为假
		if got := typing.Condition(false, "yes", "no"); got != "no" {
			t.Errorf("Condition(false, 'yes', 'no') = %v, want 'no'", got)
		}
	})

	// 测试不同类型
	t.Run("different types", func(t *testing.T) {
		// 整数
		if got := typing.Condition(5 > 3, 10, 20); got != 10 {
			t.Errorf("Condition(5 > 3, 10, 20) = %v, want 10", got)
		}

		// 字符串
		if got := typing.Condition("a" == "b", "equal", "not equal"); got != "not equal" {
			t.Errorf("Condition('a' == 'b', 'equal', 'not equal') = %v, want 'not equal'", got)
		}

		// 切片
		slice1 := []int{1, 2, 3}
		slice2 := []int{4, 5, 6}
		if got := typing.Condition(len(slice1) > 2, slice1, slice2); len(got) != 3 || got[0] != 1 {
			t.Errorf("Condition with slices failed")
		}
	})

	// 测试复杂表达式
	t.Run("complex expressions", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		person1 := Person{Name: "Alice", Age: 30}
		person2 := Person{Name: "Bob", Age: 25}

		if got := typing.Condition(person1.Age > person2.Age, person1, person2); got.Name != "Alice" {
			t.Errorf("Condition with structs failed")
		}
	})
}

// TestRealWorldUsage 测试实际使用场景
func TestRealWorldUsage(t *testing.T) {
	t.Run("config with defaults", func(t *testing.T) {
		type Config struct {
			Host  string
			Port  int
			Debug bool
		}

		// 模拟从配置文件读取的配置（可能有空值）
		var emptyConfig Config
		defaultConfig := Config{
			Host:  "localhost",
			Port:  8080,
			Debug: false,
		}

		// 使用 DefaultIfEmpty 设置默认值
		finalConfig := typing.DefaultIfEmpty(emptyConfig, defaultConfig)

		if finalConfig.Host != "localhost" || finalConfig.Port != 8080 {
			t.Errorf("Config defaults not applied correctly")
		}
	})

	t.Run("optional parameters", func(t *testing.T) {
		// 模拟可选参数的函数
		processData := func(data string, options *map[string]interface{}) map[string]interface{} {
			defaultOptions := map[string]interface{}{
				"timeout": 30,
				"retry":   3,
				"verbose": false,
			}

			finalOptions := typing.DefaultIfEmptyPtr(options, defaultOptions)
			return finalOptions
		}

		// 不传递选项
		result1 := processData("test", nil)
		if result1["timeout"] != 30 {
			t.Errorf("Default options not applied for nil pointer")
		}

		// 传递空选项
		emptyOptions := make(map[string]interface{})
		result2 := processData("test", &emptyOptions)
		if result2["timeout"] != 30 {
			t.Errorf("Default options not applied for empty map")
		}

		// 传递自定义选项
		customOptions := map[string]interface{}{
			"timeout": 60,
			"retry":   5,
		}
		result3 := processData("test", &customOptions)
		if result3["timeout"] != 60 {
			t.Errorf("Custom options not preserved")
		}
	})

	t.Run("conditional assignment", func(t *testing.T) {
		// 使用 Condition 进行条件赋值
		getStatusMessage := func(success bool, data interface{}) string {
			return typing.Condition(success,
				"Operation completed successfully",
				"Operation failed")
		}

		if got := getStatusMessage(true, nil); got != "Operation completed successfully" {
			t.Errorf("Condition function failed for success case")
		}

		if got := getStatusMessage(false, nil); got != "Operation failed" {
			t.Errorf("Condition function failed for failure case")
		}
	})
}

func BenchmarkDefaultIfEmpty(b *testing.B) {
	b.Run("string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			typing.DefaultIfEmpty("", "default")
		}
	})

	b.Run("int", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			typing.DefaultIfEmpty(0, 42)
		}
	})

	b.Run("slice", func(b *testing.B) {
		var nilSlice []int
		defaultSlice := []int{1, 2, 3}
		for i := 0; i < b.N; i++ {
			typing.DefaultIfEmpty(nilSlice, defaultSlice)
		}
	})
}

func BenchmarkDefaultIfEmptyPtr(b *testing.B) {
	b.Run("string", func(b *testing.B) {
		var nilPtr *string
		for i := 0; i < b.N; i++ {
			typing.DefaultIfEmptyPtr(nilPtr, "default")
		}
	})

	b.Run("int", func(b *testing.B) {
		var nilPtr *int
		for i := 0; i < b.N; i++ {
			typing.DefaultIfEmptyPtr(nilPtr, 42)
		}
	})
}

func BenchmarkIsEmpty(b *testing.B) {
	b.Run("string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			typing.IsEmpty("")
		}
	})

	b.Run("int", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			typing.IsEmpty(0)
		}
	})

	b.Run("pointer", func(b *testing.B) {
		var ptr *int
		for i := 0; i < b.N; i++ {
			typing.IsEmpty(ptr)
		}
	})

	b.Run("slice", func(b *testing.B) {
		var slice []int
		for i := 0; i < b.N; i++ {
			typing.IsEmpty(slice)
		}
	})

	b.Run("struct", func(b *testing.B) {
		type Person struct {
			Name string
			Age  int
		}
		var person Person
		for i := 0; i < b.N; i++ {
			typing.IsEmpty(person)
		}
	})
}
