package mapsdk

import (
	"math"
)

// 坐标系统常量
const (
	X_PI = 3.14159265358979324 * 3000.0 / 180.0 // π * 3000 / 180
	PI   = 3.1415926535897932384626             // π
	A    = 6378245.0                            // 长半轴
	EE   = 0.00669342162296594323               // 偏心率平方
)

// CoordinateType 坐标系统类型
type CoordinateType int

const (
	WGS84 CoordinateType = iota + 1 // GPS原始坐标
	GCJ02                           // 国测局加密坐标（火星坐标）
	BD09                            // 百度加密坐标
)

func (v CoordinateType) IsOK() bool {
	return v == WGS84 || v == GCJ02 || v == BD09
}

// CoordinateConverter 坐标转换器
type CoordinateConverter struct{}

// NewCoordinateConverter 创建坐标转换器实例
func NewCoordinateConverter() *CoordinateConverter {
	return &CoordinateConverter{}
}

// Transform 坐标转换主函数
func (c *CoordinateConverter) Transform(point Point, from, to CoordinateType) Point {
	if from == to {
		return Point{Lat: point.Lat, Lng: point.Lng}
	}

	switch from {
	case WGS84:
		switch to {
		case GCJ02:
			return c.WGS84ToGCJ02(point)
		case BD09:
			gcj02 := c.WGS84ToGCJ02(point)
			return c.GCJ02ToBD09(gcj02)
		}
	case GCJ02:
		switch to {
		case WGS84:
			return c.GCJ02ToWGS84(point)
		case BD09:
			return c.GCJ02ToBD09(point)
		}
	case BD09:
		switch to {
		case GCJ02:
			return c.BD09ToGCJ02(point)
		case WGS84:
			gcj02 := c.BD09ToGCJ02(point)
			return c.GCJ02ToWGS84(gcj02)
		}
	}

	return point
}

// WGS84ToGCJ02 WGS84坐标转GCJ02坐标
func (c *CoordinateConverter) WGS84ToGCJ02(point Point) Point {
	if c.outOfChina(point.Lat, point.Lng) {
		return Point{Lat: point.Lat, Lng: point.Lng}
	}

	dlat := c.transformLat(point.Lng-105.0, point.Lat-35.0)
	dlng := c.transformLng(point.Lng-105.0, point.Lat-35.0)

	radlat := point.Lat / 180.0 * PI
	magic := math.Sin(radlat)
	magic = 1 - EE*magic*magic
	sqrtmagic := math.Sqrt(magic)
	dlat = (dlat * 180.0) / ((A * (1 - EE)) / (magic * sqrtmagic) * PI)
	dlng = (dlng * 180.0) / (A / sqrtmagic * math.Cos(radlat) * PI)

	return Point{
		Lat: point.Lat + dlat,
		Lng: point.Lng + dlng,
	}
}

// GCJ02ToWGS84 GCJ02坐标转WGS84坐标
func (c *CoordinateConverter) GCJ02ToWGS84(point Point) Point {
	if c.outOfChina(point.Lat, point.Lng) {
		return Point{Lat: point.Lat, Lng: point.Lng}
	}

	dlat := c.transformLat(point.Lng-105.0, point.Lat-35.0)
	dlng := c.transformLng(point.Lng-105.0, point.Lat-35.0)

	radlat := point.Lat / 180.0 * PI
	magic := math.Sin(radlat)
	magic = 1 - EE*magic*magic
	sqrtmagic := math.Sqrt(magic)
	dlat = (dlat * 180.0) / ((A * (1 - EE)) / (magic * sqrtmagic) * PI)
	dlng = (dlng * 180.0) / (A / sqrtmagic * math.Cos(radlat) * PI)

	return Point{
		Lat: point.Lat - dlat,
		Lng: point.Lng - dlng,
	}
}

// GCJ02ToBD09 GCJ02坐标转BD09坐标
func (c *CoordinateConverter) GCJ02ToBD09(point Point) Point {
	z := math.Sqrt(point.Lng*point.Lng+point.Lat*point.Lat) + 0.00002*math.Sin(point.Lat*X_PI)
	theta := math.Atan2(point.Lat, point.Lng) + 0.000003*math.Cos(point.Lng*X_PI)

	return Point{
		Lng: z*math.Cos(theta) + 0.0065,
		Lat: z*math.Sin(theta) + 0.006,
	}
}

// BD09ToGCJ02 BD09坐标转GCJ02坐标
func (c *CoordinateConverter) BD09ToGCJ02(point Point) Point {
	x := point.Lng - 0.0065
	y := point.Lat - 0.006
	z := math.Sqrt(x*x+y*y) - 0.00002*math.Sin(y*X_PI)
	theta := math.Atan2(y, x) - 0.000003*math.Cos(x*X_PI)

	return Point{
		Lng: z * math.Cos(theta),
		Lat: z * math.Sin(theta),
	}
}

// WGS84ToBD09 WGS84坐标直接转BD09坐标
func (c *CoordinateConverter) WGS84ToBD09(point Point) Point {
	gcj02 := c.WGS84ToGCJ02(point)
	return c.GCJ02ToBD09(gcj02)
}

// BD09ToWGS84 BD09坐标直接转WGS84坐标
func (c *CoordinateConverter) BD09ToWGS84(point Point) Point {
	gcj02 := c.BD09ToGCJ02(point)
	return c.GCJ02ToWGS84(gcj02)
}

// transformLat 纬度转换
func (c *CoordinateConverter) transformLat(lng, lat float64) float64 {
	ret := -100.0 + 2.0*lng + 3.0*lat + 0.2*lat*lat + 0.1*lng*lat + 0.2*math.Sqrt(math.Abs(lng))
	ret += (20.0*math.Sin(6.0*lng*PI) + 20.0*math.Sin(2.0*lng*PI)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lat*PI) + 40.0*math.Sin(lat/3.0*PI)) * 2.0 / 3.0
	ret += (160.0*math.Sin(lat/12.0*PI) + 320*math.Sin(lat*PI/30.0)) * 2.0 / 3.0
	return ret
}

// transformLng 经度转换
func (c *CoordinateConverter) transformLng(lng, lat float64) float64 {
	ret := 300.0 + lng + 2.0*lat + 0.1*lng*lng + 0.1*lng*lat + 0.1*math.Sqrt(math.Abs(lng))
	ret += (20.0*math.Sin(6.0*lng*PI) + 20.0*math.Sin(2.0*lng*PI)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lng*PI) + 40.0*math.Sin(lng/3.0*PI)) * 2.0 / 3.0
	ret += (150.0*math.Sin(lng/12.0*PI) + 300.0*math.Sin(lng/30.0*PI)) * 2.0 / 3.0
	return ret
}

// outOfChina 判断是否在中国境外
func (c *CoordinateConverter) outOfChina(lat, lng float64) bool {
	return lng < 72.004 || lng > 137.8347 || lat < 0.8293 || lat > 55.8271
}

// Distance 计算两点间距离（米）
func (c *CoordinateConverter) Distance(point1, point2 Point) float64 {
	radLat1 := point1.Lat * PI / 180.0
	radLat2 := point2.Lat * PI / 180.0
	deltaLat := radLat1 - radLat2
	deltaLng := (point1.Lng - point2.Lng) * PI / 180.0

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(radLat1)*math.Cos(radLat2)*math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	distance := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return 6378137.0 * distance // 地球半径（米）
}

// Bearing 计算方位角（度）
func (c *CoordinateConverter) Bearing(point1, point2 Point) float64 {
	radLat1 := point1.Lat * PI / 180.0
	radLat2 := point2.Lat * PI / 180.0
	deltaLng := (point2.Lng - point1.Lng) * PI / 180.0

	y := math.Sin(deltaLng) * math.Cos(radLat2)
	x := math.Cos(radLat1)*math.Sin(radLat2) - math.Sin(radLat1)*math.Cos(radLat2)*math.Cos(deltaLng)

	bearing := math.Atan2(y, x) * 180.0 / PI
	return math.Mod(bearing+360.0, 360.0)
}

// IsValidCoordinate 验证坐标是否有效
func (c *CoordinateConverter) IsValidCoordinate(point Point) bool {
	return point.Lat >= -90.0 && point.Lat <= 90.0 &&
		point.Lng >= -180.0 && point.Lng <= 180.0
}
