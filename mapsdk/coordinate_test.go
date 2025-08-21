package mapsdk

import (
	"math"
	"testing"
)

func TestCoordinateConverter_WGS84ToGCJ02(t *testing.T) {
	converter := NewCoordinateConverter()

	// 测试北京天安门坐标 (WGS84: 116.391272, 39.906641)
	wgs84Point := Point{Lat: 39.906641, Lng: 116.391272}
	gcj02Point := converter.WGS84ToGCJ02(wgs84Point)

	// 验证转换后坐标发生了偏移
	if gcj02Point.Lat == wgs84Point.Lat && gcj02Point.Lng == wgs84Point.Lng {
		t.Error("Coordinate should be transformed")
	}

	// 验证转换结果在合理范围内
	if math.Abs(gcj02Point.Lat-wgs84Point.Lat) > 0.01 || math.Abs(gcj02Point.Lng-wgs84Point.Lng) > 0.01 {
		t.Logf("WGS84: %+v, GCJ02: %+v", wgs84Point, gcj02Point)
	}
}

func TestCoordinateConverter_GCJ02ToWGS84(t *testing.T) {
	converter := NewCoordinateConverter()

	// 先从WGS84转换到GCJ02，再转换回来
	originalPoint := Point{Lat: 39.906641, Lng: 116.391272}
	gcj02Point := converter.WGS84ToGCJ02(originalPoint)
	wgs84Point := converter.GCJ02ToWGS84(gcj02Point)

	// 验证往返转换的精度
	if math.Abs(wgs84Point.Lat-originalPoint.Lat) > 0.01 {
		t.Errorf("Round-trip conversion error too large: %f vs %f", wgs84Point.Lat, originalPoint.Lat)
	}
	if math.Abs(wgs84Point.Lng-originalPoint.Lng) > 0.01 {
		t.Errorf("Round-trip conversion error too large: %f vs %f", wgs84Point.Lng, originalPoint.Lng)
	}
}

func TestCoordinateConverter_GCJ02ToBD09(t *testing.T) {
	converter := NewCoordinateConverter()

	// 测试GCJ02转BD09
	gcj02Point := Point{Lat: 39.908042, Lng: 116.397769}
	bd09Point := converter.GCJ02ToBD09(gcj02Point)

	// 验证转换后坐标发生了偏移
	if bd09Point.Lat == gcj02Point.Lat && bd09Point.Lng == gcj02Point.Lng {
		t.Error("Coordinate should be transformed")
	}

	t.Logf("GCJ02: %+v, BD09: %+v", gcj02Point, bd09Point)
}

func TestCoordinateConverter_BD09ToGCJ02(t *testing.T) {
	converter := NewCoordinateConverter()

	// 先从GCJ02转换到BD09，再转换回来
	originalPoint := Point{Lat: 39.908042, Lng: 116.397769}
	bd09Point := converter.GCJ02ToBD09(originalPoint)
	gcj02Point := converter.BD09ToGCJ02(bd09Point)

	// 验证往返转换的精度
	if math.Abs(gcj02Point.Lat-originalPoint.Lat) > 0.001 {
		t.Errorf("Round-trip conversion error: %f vs %f", gcj02Point.Lat, originalPoint.Lat)
	}
	if math.Abs(gcj02Point.Lng-originalPoint.Lng) > 0.001 {
		t.Errorf("Round-trip conversion error: %f vs %f", gcj02Point.Lng, originalPoint.Lng)
	}
}

func TestCoordinateConverter_Transform(t *testing.T) {
	converter := NewCoordinateConverter()

	// 测试WGS84到BD09的直接转换
	wgs84Point := Point{Lat: 39.906641, Lng: 116.391272}
	bd09Point := converter.Transform(wgs84Point, WGS84, BD09)

	// 验证转换后坐标发生了偏移
	if bd09Point.Lat == wgs84Point.Lat && bd09Point.Lng == wgs84Point.Lng {
		t.Error("Coordinate should be transformed")
	}

	// 测试相同坐标系转换
	samePoint := converter.Transform(wgs84Point, WGS84, WGS84)
	if samePoint.Lat != wgs84Point.Lat || samePoint.Lng != wgs84Point.Lng {
		t.Errorf("Same coordinate system transform should return same point")
	}
}

func TestCoordinateConverter_OutOfChina(t *testing.T) {
	converter := NewCoordinateConverter()

	// 测试中国境内坐标
	chinaPoint := Point{Lat: 39.906641, Lng: 116.391272}
	gcj02Point := converter.WGS84ToGCJ02(chinaPoint)

	// 应该有偏移
	if gcj02Point.Lat == chinaPoint.Lat && gcj02Point.Lng == chinaPoint.Lng {
		t.Error("China coordinate should be transformed")
	}

	// 测试中国境外坐标（纽约）
	outsidePoint := Point{Lat: 40.7128, Lng: -74.0060}
	gcj02Outside := converter.WGS84ToGCJ02(outsidePoint)

	// 应该没有偏移
	if gcj02Outside.Lat != outsidePoint.Lat || gcj02Outside.Lng != outsidePoint.Lng {
		t.Error("Outside China coordinate should not be transformed")
	}
}

func TestCoordinateConverter_Distance(t *testing.T) {
	converter := NewCoordinateConverter()

	// 北京天安门和故宫的距离
	tiananmen := Point{Lat: 39.906641, Lng: 116.391272}
	gugong := Point{Lat: 39.918030, Lng: 116.396969}

	distance := converter.Distance(tiananmen, gugong)

	// 天安门到故宫大约1.3公里
	if distance < 1000 || distance > 1600 {
		t.Errorf("Expected distance around 1300m, got %f", distance)
	}
}

func TestCoordinateConverter_Bearing(t *testing.T) {
	converter := NewCoordinateConverter()

	// 测试方位角
	point1 := Point{Lat: 39.906641, Lng: 116.391272}
	point2 := Point{Lat: 39.918030, Lng: 116.396969}

	bearing := converter.Bearing(point1, point2)

	// 验证方位角在0-360度之间
	if bearing < 0 || bearing >= 360 {
		t.Errorf("Bearing should be between 0 and 360, got %f", bearing)
	}
}

func TestCoordinateConverter_IsValidCoordinate(t *testing.T) {
	converter := NewCoordinateConverter()

	// 测试有效坐标
	validPoint := Point{Lat: 39.906641, Lng: 116.391272}
	if !converter.IsValidCoordinate(validPoint) {
		t.Error("Valid coordinate should return true")
	}

	// 测试无效坐标 - 纬度超出范围
	invalidLat := Point{Lat: 91.0, Lng: 116.391272}
	if converter.IsValidCoordinate(invalidLat) {
		t.Error("Invalid latitude should return false")
	}

	// 测试无效坐标 - 经度超出范围
	invalidLng := Point{Lat: 39.906641, Lng: 181.0}
	if converter.IsValidCoordinate(invalidLng) {
		t.Error("Invalid longitude should return false")
	}

	// 测试nil坐标
	if converter.IsValidCoordinate(Point{}) {
		t.Error("Nil coordinate should return false")
	}
}

func BenchmarkWGS84ToGCJ02(b *testing.B) {
	converter := NewCoordinateConverter()
	point := Point{Lat: 39.906641, Lng: 116.391272}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		converter.WGS84ToGCJ02(point)
	}
}

func BenchmarkGCJ02ToBD09(b *testing.B) {
	converter := NewCoordinateConverter()
	point := Point{Lat: 39.908042, Lng: 116.397769}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		converter.GCJ02ToBD09(point)
	}
}

func BenchmarkTransform(b *testing.B) {
	converter := NewCoordinateConverter()
	point := Point{Lat: 39.906641, Lng: 116.391272}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		converter.Transform(point, WGS84, BD09)
	}
}
