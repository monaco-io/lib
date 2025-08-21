package mapsdk

import (
	"encoding/json"

	"github.com/monaco-io/lib/xec"
)

var (
	OK                      = xec.New(0o000, "ok")
	ErrorParam              = xec.New(-1000, "invalid parameter")
	ErrorServerInternal     = xec.New(-1001, "服务器内部错误")
	ErrorParameterInvalid   = xec.New(-1002, "请求参数非法")
	ErrorVerifyFailure      = xec.New(-1003, "权限校验失败")
	ErrorQuotaFailure       = xec.New(-1004, "配额校验失败")
	ErrorAKFailure          = xec.New(-1005, "ak不存在或者非法")
	ErrorParseProto         = xec.New(-1008, "数据解析失败")
	ErrorPermissionDenied   = xec.New(-1009, "高级权限校验失败")
	ErrorAKNotExist         = xec.New(-1101, "AK参数不存在")
	ErrorAPPNotExist        = xec.New(-1200, "APP不存在，AK有误请检查再重试")
	ErrorAPPDisabled        = xec.New(-1201, "APP被用户自己禁用，请在控制台解禁")
	ErrorAPPDeleted         = xec.New(-1202, "APP被管理员删除")
	ErrorAPPTypeWrong       = xec.New(-1203, "APP类型错误")
	ErrorAPPIPCheck         = xec.New(-1210, "APP IP校验失败")
	ErrorAPPSNCheck         = xec.New(-1211, "APP SN校验失败")
	ErrorAPPServiceDisabled = xec.New(-1240, "APP 服务被禁用")
	ErrorUserNotExist       = xec.New(-1250, "用户不存在")
	ErrorUserDeleted        = xec.New(-1251, "用户被自己删除")
	ErrorUserBanned         = xec.New(-1252, "用户被管理员删除")
	ErrorServiceNotExist    = xec.New(-1260, "服务不存在")
	ErrorServiceDisabled    = xec.New(-1261, "服务被禁用")
	ErrorQuotaExceeded      = xec.New(-1302, "天配额超限，限制访问")
	ErrorConcurrencyLimit   = xec.New(-1401, "当前并发量已经超过约定并发配额，限制访问")
)

type Source int

const (
	SourceBaidu Source = iota + 1
	SourceGaode
	SourceTencent
)

type ISDK interface {
	// SearchPlace(query string, region string) (*SearchResponse, error)
	GetPlaceDetail(id string) (*Response[*PlaceDetailData], error)
}

func New(source Source, ak string) ISDK {
	switch source {
	case SourceBaidu:
		return newBaidu(ak)
	case SourceGaode:
		// return newGaode(ak)
	case SourceTencent:
		// return newTencent(ak)
	}
	return nil
}

type Location struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	Province string   `json:"province"`
	City     string   `json:"city"`
	Area     string   `json:"area"`
	Tag      string   `json:"tag"`
	Points   []TPoint `json:"points"`
}

type Response[T any] struct {
	Source Source    `json:"source"`
	Status xec.Error `json:"status"`
	Data   T         `json:"data"`
}

func (r *Response[T]) IsOK() bool {
	return r.Status.Code == 0
}

func (r *Response[T]) ToJSON() string {
	data, _ := json.Marshal(r)
	return string(data)
}

type SearchPlaceData struct{}

type PlaceDetailData struct {
	Location *Location `json:"location"`
	Extra    string    `json:"extra"`
}
