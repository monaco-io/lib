package xtime

import (
	"reflect"
	"testing"
	"time"
)

func TestDateCN(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test1", args{time.Date(2022, 1, 1, 0, 0, 0, 0, Shanghai)}, "2022-01-01"},
		{"test2", args{time.Date(2022, 1, 1, 12, 0, 0, 0, Shanghai)}, "2022-01-01"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DateCN(tt.args.t); got != tt.want {
				t.Errorf("DateCN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateTimeCN(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test1", args{time.Date(2022, 1, 1, 12, 30, 45, 0, Shanghai)}, "2022-01-01 12:30:45"},
		{"test1", args{time.Date(2022, 1, 1, 12, 30, 45, 0, time.UTC)}, "2022-01-01 20:30:45"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DateTimeCN(tt.args.t); got != tt.want {
				t.Errorf("DateTimeCN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseLocal(t *testing.T) {
	type args struct {
		t      string
		format []string
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			"test1",
			args{"2022-01-01 12:30:45", []string{time.RFC3339, time.RFC1123}},
			time.Date(2022, 1, 1, 12, 30, 45, 0, time.Local),
		},
		{
			"test2_RFC3339",
			args{t: "2023-05-15T14:30:45Z"},
			time.Date(2023, 5, 15, 14, 30, 45, 0, time.UTC).In(time.Local),
		},
		{
			"test3_RFC822",
			args{t: "15 May 23 14:30 UTC", format: []string{time.RFC822}},
			time.Date(2023, 5, 15, 14, 30, 0, 0, time.UTC).In(time.Local),
		},
		{
			"test4_Kitchen",
			args{t: "2:30PM", format: []string{time.Kitchen}},
			time.Date(0, 1, 1, 14, 30, 0, 0, time.Local),
		},
		{
			"test5_Stamp",
			args{t: "May 15 14:30:45"},
			time.Date(0, 5, 15, 14, 30, 45, 0, time.Local),
		},
		{
			"test6_custom_format",
			args{"2023/12/25 09:15:30", []string{"2006/01/02 15:04:05"}},
			time.Date(2023, 12, 25, 9, 15, 30, 0, time.Local),
		},
		{
			"test7_date_only",
			args{"2023-07-04", []string{"2006-01-02"}},
			time.Date(2023, 7, 4, 0, 0, 0, 0, time.Local),
		},
		{
			"test8_multiple_formats",
			args{"Mon, 02 Jan 2006 15:04:05 MST", []string{"2006-01-02", time.RFC1123}},
			time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("MST", -7*3600)).In(time.Local),
		},
		{
			"test9_unix_timestamp_string",
			args{"1640995200", []string{"unix"}},
			time.Unix(1640995200, 0).In(time.Local),
		},
		{
			"test10_iso8601",
			args{"2023-08-20T16:45:30+08:00", []string{time.RFC3339}},
			time.Date(2023, 8, 20, 16, 45, 30, 0, time.FixedZone("", 8*3600)).In(time.Local),
		},
		{
			"test11_leap_year",
			args{t: "2024-02-29 23:59:59"},
			time.Date(2024, 2, 29, 23, 59, 59, 0, Shanghai),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseLocal(tt.args.t, tt.args.format...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseLocal() = %v, want %v", got, tt.want)
			}
		})
	}
}
