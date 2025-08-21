package mapsdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/monaco-io/lib/xec"
)

var coordinateConverter = NewCoordinateConverter()

type baidu struct {
	host string
	ak   string
}

func newBaidu(ak string) *baidu {
	return &baidu{
		host: "https://api.map.baidu.com",
		ak:   ak,
	}
}

type BaiduSearchResponse struct {
	Status     int    `json:"status"`
	Message    string `json:"message"`
	ResultType string `json:"result_type"`
	Results    []struct {
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
	} `json:"results"`
}

type baiduDetailResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Result  struct {
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

func (d *BaiduSearchResponse) dto() *Response[*SearchPlaceData] {
	return &Response[*SearchPlaceData]{
		Source: SourceBaidu,
		Status: xec.New(d.Status, d.Message),
		Data:   &SearchPlaceData{},
	}
}

func (d *baiduDetailResponse) dto() *Response[*PlaceDetailData] {
	ext, _ := json.Marshal(d.Result.DetailInfo)
	bd09 := TPoint{CoordinateType: BD09, Point: Point{Lat: d.Result.Location.Lat, Lng: d.Result.Location.Lng}}
	return &Response[*PlaceDetailData]{
		Source: SourceBaidu,
		Status: xec.New(d.Status, d.Message),
		Data: &PlaceDetailData{
			Location: &Location{
				ID:       d.Result.UID,
				Name:     d.Result.Name,
				Address:  d.Result.Address,
				Province: d.Result.Province,
				City:     d.Result.City,
				Area:     d.Result.Area,
				Tag:      d.Result.DetailInfo.Tag,
				Points: []TPoint{
					bd09, bd09.Transform(GCJ02), bd09.Transform(WGS84),
				},
			},
			Extra: string(ext),
		},
	}
}

func (b *baidu) SearchPlace(query string, region string) (*Response[*SearchPlaceData], error) {
	// 接口地址
	uri := "/place/v2/search"

	// 设置请求参数
	params := url.Values{
		"query":  []string{query},
		"region": []string{region},
		"output": []string{"json"},
		"ak":     []string{b.ak},
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

func (b *baidu) GetPlaceDetail(id string) (*Response[*PlaceDetailData], error) {
	// 接口地址
	uri := "/place/v2/detail"

	// 设置请求参数
	params := url.Values{
		"uid":    []string{id},
		"output": []string{"json"},
		"scope":  []string{"2"},
		"ak":     []string{b.ak},
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
