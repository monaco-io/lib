package mapsdk

import (
	"os"
	"testing"

	"github.com/monaco-io/lib/mapsdk/coordinate"
)

var ak = os.Getenv("BAIDU_MAP_AK")

func Test_baidu_GetReverseGeocoding(t *testing.T) {
	type fields struct{}
	type args struct {
		poi      TPoint
		poiTypes []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Response[ReverseGeocodingData]
		wantErr bool
	}{
		{
			name:   "",
			fields: fields{},
			args: args{
				poi:      TPoint{coordinate.BD09, Point{40.049612216322984, 116.29535438352558}},
				poiTypes: []string{},
			},
			want:    &Response[ReverseGeocodingData]{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := newBaidu(ak)
			got, err := b.GetReverseGeocoding(tt.args.poi, tt.args.poiTypes)
			if (err != nil) != tt.wantErr {
				t.Errorf("baidu.GetReverseGeocoding() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("baidu.GetReverseGeocoding() = %s", got.ToJSON())
		})
	}
}

func Test_baidu_GetPlaceDetail(t *testing.T) {
	type fields struct{}
	type args struct {
		id []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Response[ReverseGeocodingData]
		wantErr bool
	}{
		{
			name:   "",
			fields: fields{},
			args: args{
				id: []string{"435d7aea036e54355abbbcc8"},
			},
			want:    &Response[ReverseGeocodingData]{},
			wantErr: false,
		},
		{
			name:   "",
			fields: fields{},
			args: args{
				id: []string{"435d7aea036e54355abbbcc8", "ec27de7f4b178af1a2d23075"},
			},
			want:    &Response[ReverseGeocodingData]{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := newBaidu(ak)
			got, err := b.GetPlaceDetail(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("baidu.GetPlaceDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("baidu.GetPlaceDetail() = %s", got.ToJSON())
		})
	}
}

func Test_baidu_SearchPlace(t *testing.T) {
	type fields struct{}
	type args struct {
		query  string
		region string
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
				query:  "洗浴",
				region: "上海",
			},
			want:    &Response[SearchPlaceData]{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := newBaidu(ak)
			got, err := b.SearchPlace(tt.args.query, tt.args.region)
			if (err != nil) != tt.wantErr {
				t.Errorf("baidu.SearchPlace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("baidu.SearchPlace() = %v", got.ToJSON())
		})
	}
}
