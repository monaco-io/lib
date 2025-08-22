package mapsdk

import (
	"encoding/json"
	"fmt"

	"github.com/monaco-io/lib/xec"
	"github.com/monaco-io/lib/xjson"
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

type Source string

const (
	Baidu   Source = "baidu"
	Gaode   Source = "gaode"
	Tencent Source = "tencent"
)

type ISDK interface {
	// 地点检索
	// baidu 地点检索3.0 https://lbs.baidu.com/faq/api?title=webapi/guide/webservice-placeapiV3/interfaceDocumentV3
	SearchRegion(params SearchRegionParams, opts ...KV) (*Response[SearchPlaceData], error)

	// 地点详情检索
	// baidu https://lbs.baidu.com/faq/api?title=webapi/guide/webservice-placeapi/detail
	GetPlaceDetail(params GetPlaceDetailParams, opts ...KV) (*Response[[]PlaceDetailData], error)

	// 全球逆地理编码
	// baidu https://lbs.baidu.com/faq/api?title=webapi/guide/webservice-geocoding-abroad-base
	GetReverseGeocoding(params GetReverseGeocodingParams, opts ...KV) (*Response[ReverseGeocodingData], error)

	NativeDo(uri string, kv ...KV) (json.RawMessage, error)
}

func New(source Source, ak string) ISDK {
	switch source {
	case Baidu:
		return newBaidu(ak)
	case Gaode:
		// return newGaode(ak)
	case Tencent:
		// return newTencent(ak)
	}
	return nil
}

type KV func() (k, v string)

func NewKV(k, v string) func() (string, string) {
	return func() (string, string) {
		return k, v
	}
}

type IResponseDTO[T any] interface {
	ResponseDTO() *Response[T]
}
type Location struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`          // 百度大厦
	Address      string         `json:"address"`       // 北京市海淀区上地十街10号
	Country      string         `json:"country"`       // 中国
	Province     string         `json:"province"`      // 北京市
	City         string         `json:"city"`          // 北京市
	Area         string         `json:"area"`          // 海淀区
	Street       string         `json:"street"`        // 上地十街
	Town         string         `json:"town"`          // 上地街道
	StreetNumber string         `json:"street_number"` // 10号
	Telephone    string         `json:"telephone"`     // 电话(021)38751245
	Tags         []string       `json:"tags"`          // 房地产,写字楼
	Point        Point          `json:"points"`
	Detail       LocationDetail `json:"detail"`
}

type LocationDetail struct {
	Classification []string `json:"classification"` // 类别
	Type           string   `json:"type"`
	ShopHours      string   `json:"shop_hours"` // 营业时间
	Price          string   `json:"price"`
	Label          []string `json:"label"`  // 标签
	Photos         []string `json:"photos"` // poi的图片
}

type Response[T any] struct {
	Source `json:"source"`
	Status xec.Error       `json:"status"`
	Data   T               `json:"data"`
	Meta   json.RawMessage `json:"meta"`
}

func (r *Response[T]) IsOK() bool {
	return r.Status.Code == 0
}

func (r *Response[T]) ToJSON() string {
	data, _ := json.Marshal(r)
	return string(data)
}

type SearchPlaceData struct {
	Locations []Location `json:"locations"`
	Extra     string     `json:"extra"`
}

type PlaceDetailData struct {
	Location Location `json:"location"`
	Extra    string   `json:"extra"`
}

type ReverseGeocodingData struct {
	Location Location   `json:"location"`
	Child    []Location `json:"child"`
	Extra    string     `json:"extra"`
}

type (
	GetPlaceDetailParams struct {
		IDs []string
	}
	GetReverseGeocodingParams struct {
		Point             // 圆形区域检索中心点，必填 格式：lat,lng
		Radius   int      // 圆形区域检索半径，单位米，默认1000
		PoiTypes []string // 控制返回附近POI类型
	}
	SearchRegionParams struct {
		Keyword string // 检索关键字，必填
		Region  string // 检索区域，必填
		Point          // 圆形区域检索中心点，必填 格式：lat,lng
	}
)

func unmarshal[IDTO IResponseDTO[T], T any](body json.RawMessage) (*Response[T], error) {
	target, err := xjson.UnmarshalT[IDTO](body)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}
	return target.ResponseDTO(), nil
}
