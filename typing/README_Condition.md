# Condition 函数使用说明

这个包提供了三个实用的条件处理函数：`Condition`、`DefaultIfEmpty` 和 `DefaultIfEmptyPtr`。

## 函数概览

### 1. `Condition[T any](cond bool, v1, v2 T) T`

三元运算符的 Go 实现，根据条件返回两个值中的一个。

```go
// 基本用法
result := typing.Condition(age >= 18, "adult", "minor")

// 不同类型
maxValue := typing.Condition(a > b, a, b)
message := typing.Condition(success, "操作成功", "操作失败")
```

### 2. `DefaultIfEmpty[T any](d T, defaultValue T) T`

如果值为空，则返回默认值，否则返回原值。

```go
// 字符串
name := typing.DefaultIfEmpty("", "Anonymous")  // 返回 "Anonymous"
name := typing.DefaultIfEmpty("Alice", "Anonymous")  // 返回 "Alice"

// 数值
port := typing.DefaultIfEmpty(0, 8080)  // 返回 8080
port := typing.DefaultIfEmpty(3000, 8080)  // 返回 3000

// 切片
data := typing.DefaultIfEmpty([]int{}, []int{1,2,3})  // 返回 [1,2,3]
data := typing.DefaultIfEmpty([]int{4,5}, []int{1,2,3})  // 返回 [4,5]

// 结构体
type Config struct {
    Host string
    Port int
}
config := typing.DefaultIfEmpty(Config{}, Config{Host: "localhost", Port: 8080})
```

### 3. `DefaultIfEmptyPtr[T any](d *T, defaultValue ...T) T`

处理指针类型，如果指针为 nil 或指向空值，则返回默认值。

```go
// 基本用法
var nilPtr *string
result := typing.DefaultIfEmptyPtr(nilPtr, "default")  // 返回 "default"

// 指向空值的指针
empty := ""
emptyPtr := &empty
result := typing.DefaultIfEmptyPtr(emptyPtr, "default")  // 返回 "default"

// 指向非空值的指针
value := "hello"
valuePtr := &value
result := typing.DefaultIfEmptyPtr(valuePtr, "default")  // 返回 "hello"

// 无默认值（返回零值）
result := typing.DefaultIfEmptyPtr(nilPtr)  // 返回 ""

// 多个默认值（使用第一个）
result := typing.DefaultIfEmptyPtr(nilPtr, "first", "second")  // 返回 "first"
```

## 实际使用场景

### 配置管理

```go
type ServerConfig struct {
    Host    string
    Port    int
    Debug   bool
    Timeout int
}

func NewServerConfig(userConfig *ServerConfig) ServerConfig {
    defaultConfig := ServerConfig{
        Host:    "localhost",
        Port:    8080,
        Debug:   false,
        Timeout: 30,
    }

    return typing.DefaultIfEmptyPtr(userConfig, defaultConfig)
}
```

### API 响应处理

```go
func GetUserMessage(success bool, data interface{}, err error) string {
    return typing.Condition(success,
        "数据获取成功",
        typing.DefaultIfEmpty(err.Error(), "未知错误"))
}
```

### 可选参数处理

```go
func ProcessData(data string, options *map[string]interface{}) map[string]interface{} {
    defaultOptions := map[string]interface{}{
        "timeout": 30,
        "retry":   3,
        "verbose": false,
    }

    return typing.DefaultIfEmptyPtr(options, defaultOptions)
}
```

### 字符串处理

```go
func FormatName(firstName, lastName string) string {
    fullName := strings.TrimSpace(firstName + " " + lastName)
    return typing.DefaultIfEmpty(fullName, "Unknown User")
}
```

## 空值定义

这些函数使用 `IsEmpty` 来判断值是否为空，支持以下空值定义：

- **基本类型**: 零值（0, "", false 等）
- **指针/接口**: nil 或指向空值
- **切片/映射**: nil 或长度为 0
- **通道**: nil
- **数组**: 长度为 0
- **函数**: nil
- **结构体**: 所有字段都为零值

## 性能特性

- 对于常见的基本类型（string、int、bool 等），使用了类型断言优化，避免反射开销
- 对于复杂类型，使用反射进行通用处理
- 基准测试显示基本类型处理性能在 2-5 ns/op 范围内
- 所有函数都是零内存分配的（0 B/op）
