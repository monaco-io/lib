package mapsdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/monaco-io/lib/typing/xstring"
)

type baidu struct {
	host string
	ak   string
}

func newBaidu(ak string) *baidu {
	if ak == "" {
		panic("Baidu AK is required")
	}
	return &baidu{
		host: "https://api.map.baidu.com",
		ak:   ak,
	}
}

func (b *baidu) NativeDo(uri string, params ...KV) (json.RawMessage, error) {
	values := url.Values{}
	for _, opt := range params {
		k, v := opt()
		values.Set(k, v)
	}
	values.Set("ak", b.ak)
	values.Set("output", "json")
	url, err := url.Parse(b.host + uri + "?" + values.Encode())
	if err != nil {
		return nil, fmt.Errorf("host error: %v", err)
	}
	resp, err := http.Get(url.String())
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error: %v", err)
	}
	return body, nil
}

// https://lbs.baidu.com/faq/api?title=webapi/guide/webservice-placeapiV3/interfaceDocumentV3
func (b *baidu) SearchRegion(params SearchRegionParams, opts ...KV) (*Response[SearchPlaceData], error) {
	opts = append(opts, NewKV("query", params.Keyword))
	opts = append(opts, NewKV("region", params.Region))
	if params.Lat != 0 && params.Lng != 0 {
		opts = append(opts, NewKV("center", fmt.Sprintf("%f,%f", params.Lat, params.Lng)))
	}

	body, err := b.NativeDo("/place/v3/region", opts...)
	if err != nil {
		return nil, err
	}
	return unmarshal[*baiduSearchResponse3](body)
}

func (b *baidu) GetPlaceDetail(params GetPlaceDetailParams, opts ...KV) (*Response[[]PlaceDetailData], error) {
	opts = append(opts, NewKV("uids", strings.Join(params.IDs, ",")))
	body, err := b.NativeDo("/place/v2/detail", opts...)
	if err != nil {
		return nil, err
	}
	return unmarshal[*baiduDetailResponse](body)
}

// 控制返回附近POI类型
// 以下内容需要 extensions_poi=1时才生效；
// 可以选择poi类型召回不同类型的poi，例如poi_types=酒店，如想召回多个POI类型数据，可以‘|’分割
// 例如poi_types=酒店|房地产 不添加该参数则默认召回全部POI分类数据。
// poi分类 https://lbsyun.baidu.com/index.php?title=open/poitags
func (b *baidu) GetReverseGeocoding(params GetReverseGeocodingParams, opts ...KV) (*Response[ReverseGeocodingData], error) {
	opts = append(opts, NewKV("location", fmt.Sprintf("%f,%f", params.Lat, params.Lng)))
	if len(params.PoiTypes) > 0 {
		opts = append(opts, NewKV("poi_types", strings.Join(params.PoiTypes, "|")))
	}
	opts = append(opts, NewKV("radius", xstring.Pick("1000", fmt.Sprintf("%d", params.Radius)))) // 确保开启poi扩展
	body, err := b.NativeDo("/reverse_geocoding/v3", opts...)
	if err != nil {
		return nil, err
	}
	return unmarshal[*baiduReverseGeocodingResponse3](body)
}
