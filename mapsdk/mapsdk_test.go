package mapsdk

import (
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		source Source
		ak     string
	}
	tests := []struct {
		name string
		args args
		want ISDK
	}{
		{
			name: "Baidu",
			args: args{
				source: SourceBaidu,
				ak:     "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sdk := New(tt.args.source, tt.args.ak)

			res, err := sdk.GetPlaceDetail("435d7aea036e54355abbbcc8")
			if err != nil {
				t.Errorf("GetPlaceDetail() error = %v", err)
				return
			}
			t.Logf("GetPlaceDetail() = %s", res.ToJSON())
		})
	}
}

func TestPoint_GetGeoHash(t *testing.T) {
	type fields struct {
		Type CoordinateType
		Lat  float64
		Lng  float64
	}
	type args struct {
		precision uint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			fields: fields{
				Lat: 39.9042,
				Lng: 116.4074,
			},
			args: args{
				precision: 5,
			},
			want: "wx4g0",
		},
		{
			fields: fields{
				Lat: 31.2304,
				Lng: 121.4737,
			},
			args: args{
				precision: 7,
			},
			want: "wtw3sjq",
		},
		{
			fields: fields{
				Lat: 37.7749,
				Lng: -122.4194,
			},
			args: args{
				precision: 4,
			},
			want: "9q8y",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Point{
				Lat: tt.fields.Lat,
				Lng: tt.fields.Lng,
			}
			if got := p.GetGeoHash(tt.args.precision); got != tt.want {
				t.Errorf("Point.GetGeoHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
