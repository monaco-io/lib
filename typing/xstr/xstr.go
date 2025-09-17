package xstr

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/exp/constraints"
)

const (
	EMPTY = ""
	SPACE = " "

	TAB             = "\t"
	NEWLINE         = "\n"
	CARRIAGE_RETURN = "\r"

	COMMA      = ","
	PERIOD     = "."
	SEMICOLON  = ";"
	HYPHEN     = "-"
	UNDERSCORE = "_"

	ARROW = "->"

	X_REQUEST_ID = "x-request-id"
)

func UUID() string {
	return uuid.NewString()
}

func UUIDX() string {
	return strings.ReplaceAll(UUID(), HYPHEN, EMPTY)
}

// IsEmpty 检查字符串是否为空
func IsEmpty(s string) bool {
	return s == EMPTY
}

// IsNotEmpty 检查字符串是否不为空
func IsNotEmpty(s string) bool {
	return s != EMPTY
}

// IsBlank 检查字符串是否为空白（空或只包含空白字符）
func IsBlank(s string) bool {
	return strings.TrimSpace(s) == EMPTY
}

// IsNotBlank 检查字符串是否不为空白
func IsNotBlank(s string) bool {
	return strings.TrimSpace(s) != EMPTY
}

// 如果字符串为空则返回默认值
func DefaultIfEmpty(s string, defaultValue ...string) string {
	if IsEmpty(s) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return EMPTY
	}
	return s
}

// 如果字符串为空则返回默认值
func DefaultIfEmptyPtr(s *string, defaultValue ...string) string {
	if s == nil || IsEmpty(*s) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return EMPTY
	}
	return *s
}

// DefaultIfBlank 如果字符串为空白则返回默认值
func DefaultIfBlank(s, defaultValue string) string {
	if IsBlank(s) {
		return defaultValue
	}
	return s
}

// Reverse 反转字符串
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Contains 检查字符串是否包含子字符串（忽略大小写）
func ContainsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// StartsWith 检查字符串是否以指定前缀开始（忽略大小写）
func StartsWithIgnoreCase(s, prefix string) bool {
	return strings.HasPrefix(strings.ToLower(s), strings.ToLower(prefix))
}

// EndsWith 检查字符串是否以指定后缀结束（忽略大小写）
func EndsWithIgnoreCase(s, suffix string) bool {
	return strings.HasSuffix(strings.ToLower(s), strings.ToLower(suffix))
}

// Truncate 截断字符串到指定长度
func Truncate(s string, length int) string {
	if length < 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= length {
		return s
	}
	return string(runes[:length])
}

// CollapseWhitespace 将连续的空白字符替换为单个空格
func CollapseWhitespace(s string) string {
	re := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(re.ReplaceAllString(s, " "))
}

// ToSnakeCase 转换为蛇形命名法
func ToSnakeCase(s string) string {
	re := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	snake := re.ReplaceAllString(s, `${1}_${2}`)
	return strings.ToLower(snake)
}

// ToCamelCase 转换为驼峰命名法
func ToCamelCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	if len(words) == 0 {
		return s
	}

	result := strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		if len(words[i]) > 0 {
			result += strings.ToUpper(string(words[i][0])) + strings.ToLower(words[i][1:])
		}
	}
	return result
}

// ToPascalCase 转换为帕斯卡命名法
func ToPascalCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	var result strings.Builder
	for _, word := range words {
		if len(word) > 0 {
			result.WriteString(strings.ToUpper(string(word[0])) + strings.ToLower(word[1:]))
		}
	}
	return result.String()
}

// ToKebabCase 转换为短横线命名法
func ToKebabCase(s string) string {
	re := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	kebab := re.ReplaceAllString(s, `${1}-${2}`)
	kebab = strings.ToLower(kebab)
	// 替换非字母数字字符为短横线
	re2 := regexp.MustCompile(`[^a-z0-9]+`)
	kebab = re2.ReplaceAllString(kebab, "-")
	// 移除开头和结尾的短横线
	return strings.Trim(kebab, "-")
}

// SwapCase 交换大小写
func SwapCase(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsUpper(r) {
			return unicode.ToLower(r)
		}
		if unicode.IsLower(r) {
			return unicode.ToUpper(r)
		}
		return r
	}, s)
}

// CountWords 统计单词数量
func CountWords(s string) int {
	fields := strings.Fields(strings.TrimSpace(s))
	return len(fields)
}

// CountLines 统计行数
func CountLines(s string) int {
	if s == "" {
		return 0
	}
	return strings.Count(s, "\n") + 1
}

// Repeat 重复字符串
func Repeat(s string, count int) string {
	if count < 0 {
		return ""
	}
	return strings.Repeat(s, count)
}

// IndexOfIgnoreCase 查找子字符串位置（忽略大小写）
func IndexOfIgnoreCase(s, substr string) int {
	return strings.Index(strings.ToLower(s), strings.ToLower(substr))
}

// LastIndexOfIgnoreCase 查找子字符串最后位置（忽略大小写）
func LastIndexOfIgnoreCase(s, substr string) int {
	return strings.LastIndex(strings.ToLower(s), strings.ToLower(substr))
}

// SplitAndTrim 分割字符串并去除每部分的空白
func SplitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// ToInt 字符串转换为整数
func ToInteger[T constraints.Integer](s string) (T, error) {
	val, err := strconv.Atoi(strings.TrimSpace(s))
	return T(val), err
}

// ToInt 字符串转换为整数
func ToIntegerX[T constraints.Integer](s string) T {
	val, _ := strconv.Atoi(strings.TrimSpace(s))
	return T(val)
}

// ToInt64 字符串转换为64位整数
func ToFloatX[T constraints.Float](s string) T {
	val, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return T(val)
}

// ToInt64 字符串转换为64位整数
func ToFloat[T constraints.Float](s string) (T, error) {
	val, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return T(val), err
}

// ToBool 字符串转换为布尔值
func ToBool(s string) (bool, error) {
	return strconv.ParseBool(strings.TrimSpace(s))
}

// ToBoolX 字符串转换为布尔值
func ToBoolX(s string) bool {
	val, _ := strconv.ParseBool(strings.TrimSpace(s))
	return val
}

// FromInteger 整数转换为字符串
func FromInteger[T constraints.Integer](i T) string {
	return strconv.Itoa(int(i))
}

// FromFloat 32位浮点数转换为字符串
func FromFloat[T constraints.Float](f T, precision int) string {
	return strconv.FormatFloat(float64(f), 'f', precision, 64)
}

// FromBool 布尔值转换为字符串
func FromBool(b bool) string {
	return strconv.FormatBool(b)
}

// MD5 计算字符串的MD5哈希值
func MD5(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

// SHA256 计算字符串的SHA256哈希值
func SHA256(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}

// Mask 掩码字符串（保留前后指定字符数，中间用*替换）
func Mask(s string, start, end int, mask string) string {
	if mask == "" {
		mask = "*"
	}

	runes := []rune(s)
	length := len(runes)

	if start < 0 {
		start = 0
	}
	if end < 0 {
		end = 0
	}
	if start+end >= length {
		return s
	}

	maskLength := length - start - end
	if maskLength <= 0 {
		return s
	}

	maskString := strings.Repeat(mask, maskLength)

	return string(runes[:start]) + maskString + string(runes[length-end:])
}

// Random 生成指定长度的随机字符串
func Random(length int, charset string) string {
	if length <= 0 {
		return ""
	}
	if charset == "" {
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}

	result := make([]byte, length)
	for i := range result {
		// 简单的伪随机，基于位置计算索引
		idx := (i*7 + 13) % len(charset)
		result[i] = charset[idx]
	}
	return string(result)
}

// Similarity 计算两个字符串的相似度（简单版本）
func Similarity(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}
	if len(s1) == 0 || len(s2) == 0 {
		return 0.0
	}

	longer, shorter := s1, s2
	if len(s1) < len(s2) {
		longer, shorter = s2, s1
	}

	longerLength := len(longer)
	editDistance := levenshteinDistance(longer, shorter)

	return float64(longerLength-editDistance) / float64(longerLength)
}

// levenshteinDistance 计算编辑距离
func levenshteinDistance(s1, s2 string) int {
	r1, r2 := []rune(s1), []rune(s2)
	rows := len(r1) + 1
	cols := len(r2) + 1

	matrix := make([][]int, rows)
	for i := range matrix {
		matrix[i] = make([]int, cols)
	}

	for i := 0; i < rows; i++ {
		matrix[i][0] = i
	}
	for j := 0; j < cols; j++ {
		matrix[0][j] = j
	}

	for i := 1; i < rows; i++ {
		for j := 1; j < cols; j++ {
			cost := 0
			if r1[i-1] != r2[j-1] {
				cost = 1
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[rows-1][cols-1]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// Format 格式化字符串（类似fmt.Sprintf的简化版本）
func Format(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}
