package mapsdk

import (
	"math"
	"testing"

	"github.com/monaco-io/lib/mapsdk/coordinate"
)

func TestPointDistanceTo(t *testing.T) {
	// 测试北京天安门广场到上海外滩的距离
	beijing := &Point{
		Lat: 39.9042,  // 北京天安门广场纬度
		Lng: 116.4074, // 北京天安门广场经度
	}

	shanghai := &Point{
		Lat: 31.2304,  // 上海外滩纬度
		Lng: 121.4737, // 上海外滩经度
	}

	distance := beijing.DistanceTo(shanghai)

	// 北京到上海的实际距离大约是1067公里
	expectedDistance := 1067000.0 // 1067公里转换为米
	tolerance := 50000.0          // 允许50公里的误差

	if math.Abs(distance-expectedDistance) > tolerance {
		t.Errorf("Expected distance around %.0f meters, got %.0f meters", expectedDistance, distance)
	}

	t.Logf("Distance from Beijing to Shanghai: %.2f km", distance/1000)
}

func TestPointDistanceToSamePoint(t *testing.T) {
	point := &Point{
		Lat: 39.9042,
		Lng: 116.4074,
	}

	distance := point.DistanceTo(point)
	if distance != 0 {
		t.Errorf("Expected distance to same point to be 0, got %.2f", distance)
	}
}

func TestPointDistanceToNil(t *testing.T) {
	point := &Point{
		Lat: 39.9042,
		Lng: 116.4074,
	}

	distance := point.DistanceTo(nil)
	if distance != 0 {
		t.Errorf("Expected distance to nil point to be 0, got %.2f", distance)
	}

	var nilPoint *Point
	distance = nilPoint.DistanceTo(point)
	if distance != 0 {
		t.Errorf("Expected distance from nil point to be 0, got %.2f", distance)
	}
}

func TestTPointDistanceTo(t *testing.T) {
	// 测试不同坐标系统下的距离计算
	// 北京天安门广场 (WGS84)
	beijing := &TPoint{
		TYPE: coordinate.WGS84,
		Point: Point{
			Lat: 39.9042,
			Lng: 116.4074,
		},
	}

	// 上海外滩 (GCJ02)
	shanghai := &TPoint{
		TYPE: coordinate.GCJ02,
		Point: Point{
			Lat: 31.2304,
			Lng: 121.4737,
		},
	}

	distance := beijing.DistanceTo(shanghai)

	// 由于坐标系统的转换，距离应该与WGS84下的距离相近
	expectedDistance := 1067000.0 // 1067公里转换为米
	tolerance := 50000.0          // 允许50公里的误差

	if math.Abs(distance-expectedDistance) > tolerance {
		t.Errorf("Expected distance around %.0f meters, got %.0f meters", expectedDistance, distance)
	}

	t.Logf("Distance from Beijing (WGS84) to Shanghai (GCJ02): %.2f km", distance/1000)
}

func TestTPointDistanceToNil(t *testing.T) {
	point := &TPoint{
		TYPE: coordinate.WGS84,
		Point: Point{
			Lat: 39.9042,
			Lng: 116.4074,
		},
	}

	distance := point.DistanceTo(nil)
	if distance != 0 {
		t.Errorf("Expected distance to nil TPoint to be 0, got %.2f", distance)
	}

	var nilPoint *TPoint
	distance = nilPoint.DistanceTo(point)
	if distance != 0 {
		t.Errorf("Expected distance from nil TPoint to be 0, got %.2f", distance)
	}
}

func TestShortDistance(t *testing.T) {
	// 测试短距离计算的精度
	point1 := &Point{
		Lat: 39.9042,
		Lng: 116.4074,
	}

	// 相距约100米的点
	point2 := &Point{
		Lat: 39.9051, // 大约向北移动100米
		Lng: 116.4074,
	}

	distance := point1.DistanceTo(point2)

	// 预期距离约100米
	expectedDistance := 100.0
	tolerance := 10.0 // 允许10米的误差

	if math.Abs(distance-expectedDistance) > tolerance {
		t.Errorf("Expected distance around %.0f meters, got %.0f meters", expectedDistance, distance)
	}

	t.Logf("Short distance: %.2f meters", distance)
}
