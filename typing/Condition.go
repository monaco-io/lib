package typing

import (
	"reflect"
)

func Condition[T any](cond bool, v1, v2 T) T {
	if cond {
		return v1
	}
	return v2
}

func DefaultIfEmpty[T any](d T, defaultValue T) T {
	if IsEmpty(d) {
		return defaultValue
	}
	return d
}

func DefaultIfEmptyPtr[T any](d *T, defaultValue ...T) (zero T) {
	if IsEmpty(d) {
		return FirstNonEmpty(defaultValue...)
	}
	return *d
}

func FirstNonEmpty[T any](values ...T) (zero T) {
	for _, v := range values {
		if !IsEmpty(v) {
			return v
		}
	}
	return zero
}

// IsEmpty 检查给定值是否为空（零值）
// 支持所有 Go 类型，包括基本类型、复合类型、指针、接口等
//
// 对于不同类型的空值定义：
//   - 基本类型：零值（0, "", false 等）
//   - 指针/接口：nil 或指向空值
//   - 切片/映射：nil 或长度为 0
//   - 通道：nil
//   - 数组：长度为 0
//   - 函数：nil
//   - 结构体：所有字段都为零值
func IsEmpty[T any](d T) bool {
	// 使用类型断言优化常见类型的性能
	switch v := any(d).(type) {
	case string:
		return v == ""
	case int:
		return v == 0
	case int8:
		return v == 0
	case int16:
		return v == 0
	case int32:
		return v == 0
	case int64:
		return v == 0
	case uint:
		return v == 0
	case uint8:
		return v == 0
	case uint16:
		return v == 0
	case uint32:
		return v == 0
	case uint64:
		return v == 0
	case uintptr:
		return v == 0
	case float32:
		return v == 0
	case float64:
		return v == 0
	case complex64:
		return v == 0
	case complex128:
		return v == 0
	case bool:
		return !v
	default:
		// 对于其他类型，使用反射
		return isEmptyReflect(reflect.ValueOf(d))
	}
}

// isEmptyReflect 使用反射检查值是否为空
func isEmptyReflect(rv reflect.Value) bool {
	if !rv.IsValid() {
		return true
	}

	switch rv.Kind() {
	case reflect.Pointer, reflect.Interface:
		if rv.IsNil() {
			return true
		}
		// 递归检查指针指向的值
		return isEmptyReflect(rv.Elem())
	case reflect.Slice, reflect.Map:
		return rv.IsNil() || rv.Len() == 0
	case reflect.Chan:
		return rv.IsNil()
	case reflect.Array:
		return rv.Len() == 0
	case reflect.String:
		return rv.Len() == 0
	case reflect.Bool:
		return !rv.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return rv.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() == 0
	case reflect.Complex64, reflect.Complex128:
		return rv.Complex() == 0
	case reflect.Func:
		return rv.IsNil()
	case reflect.Struct:
		return rv.IsZero()
	default:
		return rv.IsZero()
	}
}
