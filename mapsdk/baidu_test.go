package mapsdk

import (
	"os"
	"testing"

	. "github.com/monaco-io/lib/typing"
	"github.com/monaco-io/lib/typing/xjson"
)

var TEST_BAIDU_SDK = newBaidu(os.Getenv("BAIDU_MAP_AK"))

func Test_baidu_SearchRegion(t *testing.T) {
	type fields struct{}
	type args struct {
		params SearchRegionParams
		opts   []KV[string, string]
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Response[SearchPlaceData]
		wantErr bool
	}{
		{
			name:   "",
			fields: fields{},
			args: args{
				params: SearchRegionParams{
					Keyword: "洗浴",
					// Region:  "上海",
					Point: Point{Lat: 31.2304, Lng: 121.4737}, // 上海市中心点
				},
				opts: []KV[string, string]{NewKV("scope", "2")},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TEST_BAIDU_SDK.SearchRegion(tt.args.params, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("baidu.SearchRegion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("baidu.SearchRegion() = %v", got.ToJSON())
		})
	}
}

func Test_baidu_GetPlaceDetail(t *testing.T) {
	type fields struct{}
	type args struct {
		params GetPlaceDetailParams
		opts   []KV[string, string]
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Response[SearchPlaceData]
		wantErr bool
	}{
		{
			name:   "",
			fields: fields{},
			args: args{
				params: GetPlaceDetailParams{
					IDs: []string{"8ee4560cf91d160e6cc02cd7", "435d7aea036e54355abbbcc8"},
				},
				opts: []KV[string, string]{NewKV("scope", "2")},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TEST_BAIDU_SDK.GetPlaceDetail(tt.args.params, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("baidu.GetPlaceDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("baidu.GetPlaceDetail() = %v", got.ToJSON())
		})
	}
}

func Test_baidu_GetReverseGeocoding(t *testing.T) {
	type fields struct{}
	type args struct {
		params GetReverseGeocodingParams
		opts   []KV[string, string]
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Response[SearchPlaceData]
		wantErr bool
	}{
		{
			name:   "",
			fields: fields{},
			args: args{
				params: GetReverseGeocodingParams{
					Point: Point{Lat: 31.2304, Lng: 121.4737}, // 上海市中心点
				},
				opts: []KV[string, string]{NewKV("scope", "2"), NewKV("extensions_poi", "1"), NewKV("entire_poi", "1")},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TEST_BAIDU_SDK.GetReverseGeocoding(tt.args.params, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("baidu.GetReverseGeocoding() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("baidu.GetReverseGeocoding() = %v", got.ToJSON())
		})
	}
}

func Test_baidu_GetTransitRoute(t *testing.T) {
	type fields struct{}
	type args struct {
		params GetTransitRouteParams
		opts   []KV[string, string]
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Response[SearchPlaceData]
		wantErr bool
	}{
		{
			name:   "",
			fields: fields{},
			args: args{
				params: GetTransitRouteParams{
					From: Point{Lat: 40.056878, Lng: 116.30815},
					To:   Point{Lat: 39.909263, Lng: 116.39269},
				},
				// origin=40.056878,116.30815&destination=39.909263,116.39269
				opts: []KV[string, string]{NewKV("scope", "2"), NewKV("extensions_poi", "1"), NewKV("entire_poi", "1")},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TEST_BAIDU_SDK.GetTransitRoute(tt.args.params, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("baidu.GetTransitRoute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("baidu.GetTransitRoute() = %v", got.ToJSON())
			t.Logf("baidu.GetTransitRoute() = %v", xjson.MarshalIndentStringX(got.MetaData))
		})
	}
}
