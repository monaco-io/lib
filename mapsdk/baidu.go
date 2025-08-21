package mapsdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/monaco-io/lib/mapsdk/coordinate"
	"github.com/monaco-io/lib/typing/xstring"
	"github.com/monaco-io/lib/xec"
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

type BaiduPOI struct {
	Name     string `json:"name"`
	Location struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"location"`
	Address  string `json:"address"`
	Province string `json:"province"`
	City     string `json:"city"`
	Area     string `json:"area"`
	StreetID string `json:"street_id"`
	Detail   int    `json:"detail"`
	UID      string `json:"uid"`
}

type BaiduSearchResponse struct {
	Status     int        `json:"status"`
	Message    string     `json:"message"`
	ResultType string     `json:"result_type"`
	Results    []BaiduPOI `json:"results"`
}

type baiduDetailResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Result  []struct {
		UID      string `json:"uid"`
		StreetID string `json:"street_id"`
		Name     string `json:"name"`
		Location struct {
			Lng float64 `json:"lng"`
			Lat float64 `json:"lat"`
		} `json:"location"`
		Address    string `json:"address"`
		Province   string `json:"province"`
		City       string `json:"city"`
		Area       string `json:"area"`
		DetailInfo struct {
			Tag          string `json:"tag"`
			NaviLocation struct {
				Lng float64 `json:"lng"`
				Lat float64 `json:"lat"`
			} `json:"navi_location"`
			ShopHours     string   `json:"shop_hours"`
			Alias         []string `json:"alias"`
			DetailURL     string   `json:"detail_url"`
			Type          string   `json:"type"`
			OverallRating string   `json:"overall_rating"`
			ImageNum      string   `json:"image_num"`
			CommentNum    string   `json:"comment_num"`
			ContentTag    string   `json:"content_tag"`
		} `json:"detail_info"`
		Detail int `json:"detail"`
	} `json:"result"`
}

func (v *baiduDetailResponse) ToJSON() json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}

func (v *baiduReverseGeocodingResponse) ToJSON() json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}

func (v *BaiduSearchResponse) ToJSON() json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}

type baiduReverseGeocodingResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Result  struct {
		Location struct {
			Lng float64 `json:"lng"`
			Lat float64 `json:"lat"`
		} `json:"location"`
		FormattedAddress string `json:"formatted_address"`
		AddressComponent struct {
			Country      string `json:"country"`
			Province     string `json:"province"`
			City         string `json:"city"`
			District     string `json:"district"`
			Town         string `json:"town"`
			Street       string `json:"street"`
			StreetNumber string `json:"street_number"`
		} `json:"addressComponent"`
		POIs []struct {
			Name     string `json:"name"`
			Location struct {
				Lng float64 `json:"lng"`
				Lat float64 `json:"lat"`
			} `json:"location"`
			Address  string `json:"address"`
			Province string `json:"province"`
			City     string `json:"city"`
			Area     string `json:"area"`
			StreetID string `json:"street_id"`
			Detail   int    `json:"detail"`
			UID      string `json:"uid"`
		} `json:"pois"`
	} `json:"result"`
}

func (d *BaiduSearchResponse) dto() *Response[SearchPlaceData] {
	data := SearchPlaceData{
		Locations: make([]Location, 0, len(d.Results)),
	}
	for _, item := range d.Results {
		wgs84 := TPoint{TYPE: coordinate.WGS84, Point: Point{Lat: item.Location.Lat, Lng: item.Location.Lng}}
		data.Locations = append(data.Locations, Location{
			ID:       item.UID,
			Name:     item.Name,
			Address:  item.Address,
			Province: item.Province,
			City:     item.City,
			Area:     item.Area,
			Tag:      []string{},
			Points: []TPoint{
				wgs84, wgs84.Transform(coordinate.BD09), wgs84.Transform(coordinate.GCJ02),
			},
		})
	}
	return &Response[SearchPlaceData]{
		Source: Baidu,
		Status: xec.New(d.Status, d.Message),
		Meta:   d.ToJSON(),
		Data:   data,
	}
}

func (d *baiduDetailResponse) dto() *Response[[]PlaceDetailData] {
	var data []PlaceDetailData
	for _, item := range d.Result {
		wgs84 := TPoint{TYPE: coordinate.WGS84, Point: Point{Lat: item.Location.Lat, Lng: item.Location.Lng}}
		data = append(data, PlaceDetailData{
			Location: Location{
				ID:       item.UID,
				Name:     item.Name,
				Address:  item.Address,
				Province: item.Province,
				City:     item.City,
				Area:     item.Area,
				Tag:      strings.Split(item.DetailInfo.Tag, ";"),
				Points: []TPoint{
					wgs84, wgs84.Transform(coordinate.BD09), wgs84.Transform(coordinate.GCJ02),
				},
			},
		})
	}
	return &Response[[]PlaceDetailData]{
		Source: Baidu,
		Status: xec.New(d.Status, d.Message),
		Meta:   d.ToJSON(),
		Data:   data,
	}
}

func (d *baiduReverseGeocodingResponse) dto() *Response[ReverseGeocodingData] {
	wgs84 := TPoint{TYPE: coordinate.WGS84, Point: Point{Lat: d.Result.Location.Lat, Lng: d.Result.Location.Lng}}
	var child []Location
	ac := d.Result.AddressComponent
	for _, item := range d.Result.POIs {
		child = append(child, Location{
			ID:           item.UID,
			Name:         item.Name,
			Address:      xstring.Pick(item.Address, d.Result.FormattedAddress),
			Country:      "",
			Province:     xstring.Pick(item.Province, ac.Province),
			City:         xstring.Pick(item.City, ac.City),
			Area:         xstring.Pick(item.Area, ac.District),
			Street:       "",
			Town:         "",
			StreetNumber: "",
			Tag:          []string{},
			Points:       []TPoint{{TYPE: coordinate.WGS84, Point: Point{Lat: item.Location.Lat, Lng: item.Location.Lng}}},
		})
	}
	return &Response[ReverseGeocodingData]{
		Source: Baidu,
		Status: xec.New(d.Status, d.Message),
		Meta:   d.ToJSON(),
		Data: ReverseGeocodingData{
			Location: Location{
				ID: "",
				// Name:     d.Result.AddressComponent.Name,
				Address:      d.Result.FormattedAddress,
				Country:      ac.Country,
				Province:     ac.Province,
				City:         ac.City,
				Area:         ac.District,
				Town:         ac.Town,
				Street:       ac.Street,
				StreetNumber: ac.StreetNumber,
				Tag:          []string{},
				Points:       []TPoint{wgs84, wgs84.Transform(coordinate.GCJ02), wgs84.Transform(coordinate.BD09)},
			},
			Child: child,
			Extra: "",
		},
	}
}

func (b *baidu) SearchPlace(query string, region string) (*Response[SearchPlaceData], error) {
	// 接口地址
	uri := "/place/v2/search"

	// 设置请求参数
	params := url.Values{
		"query":         []string{query},
		"region":        []string{region},
		"output":        []string{"json"},
		"ret_coordtype": []string{"wgs84ll"},
		"ak":            []string{b.ak},
	}

	// 发起请求
	request, err := url.Parse(b.host + uri + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("host error: %v", err)
	}

	resp, err := http.Get(request.String())
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response error: %v", err)
	}

	var response BaiduSearchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}

	return response.dto(), nil
}

func (b *baidu) GetPlaceDetail(id []string) (*Response[[]PlaceDetailData], error) {
	// 接口地址
	uri := "/place/v2/detail"

	// 设置请求参数
	params := url.Values{
		"uids":              []string{strings.Join(id, ",")},
		"output":            []string{"json"},
		"ret_coordtype":     []string{"wgs84ll"},
		"scope":             []string{"2"},
		"ak":                []string{b.ak},
		"extensions_adcode": []string{"true"},
	}

	// 发起请求
	request, err := url.Parse(b.host + uri + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("host error: %v", err)
	}

	resp, err := http.Get(request.String())
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response error: %v", err)
	}

	var response baiduDetailResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}

	return response.dto(), nil
}

// 控制返回附近POI类型
// 以下内容需要 extensions_poi=1时才生效；
// 可以选择poi类型召回不同类型的poi，例如poi_types=酒店，如想召回多个POI类型数据，可以‘|’分割
// 例如poi_types=酒店|房地产 不添加该参数则默认召回全部POI分类数据。
// poi分类 https://lbsyun.baidu.com/index.php?title=open/poitags
func (b *baidu) GetReverseGeocoding(poi TPoint, poiTypes []string) (*Response[ReverseGeocodingData], error) {
	// 服务地址
	// 接口地址
	uri := "/reverse_geocoding/v3"
	wgs84 := poi.Transform(coordinate.WGS84)
	// 设置请求参数
	params := url.Values{
		"ak":             []string{b.ak},
		"output":         []string{"json"},
		"coordtype":      []string{"wgs84ll"},
		"ret_coordtype":  []string{"wgs84ll"},
		"extensions_poi": []string{"1"},
		"location":       []string{fmt.Sprintf("%f,%f", wgs84.Lat, wgs84.Lng)},
	}
	if len(poiTypes) > 0 {
		params["poi_types"] = []string{strings.Join(poiTypes, "|")}
	}
	// 发起请求
	request, err := url.Parse(b.host + uri + "?" + params.Encode())
	if nil != err {
		return nil, fmt.Errorf("host error: %v", err)
	}

	resp, err := http.Get(request.String())
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response error: %v", err)
	}

	var response baiduReverseGeocodingResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}
	return response.dto(), nil
}
