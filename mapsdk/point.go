package mapsdk

import (
	"fmt"
	"math"

	"github.com/monaco-io/lib/mapsdk/coordinate"
	"github.com/monaco-io/lib/mapsdk/geohash"
)

type TPoint struct {
	coordinate.TYPE `json:"type"`
	Point           `json:"point"`
}

type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func (p *Point) GetPointString() string {
	if p != nil {
		return fmt.Sprintf("%f,%f", p.Lat, p.Lng)
	}
	return ""
}

func (p *Point) GetGeoHash(precision uint) string {
	if p != nil {
		return geohash.EncodeWithPrecision(p.Lat, p.Lng, precision)
	}
	return ""
}

// DistanceTo 计算两个点之间的距离（单位：米）
// 使用 Haversine 公式计算地球表面两点间的大圆距离
func (p *Point) DistanceTo(other *Point) float64 {
	if p == nil || other == nil {
		return 0
	}

	const earthRadius = 6371000 // 地球半径，单位：米

	// 将度数转换为弧度
	lat1Rad := p.Lat * math.Pi / 180
	lng1Rad := p.Lng * math.Pi / 180
	lat2Rad := other.Lat * math.Pi / 180
	lng2Rad := other.Lng * math.Pi / 180

	// 计算纬度和经度的差值
	deltaLat := lat2Rad - lat1Rad
	deltaLng := lng2Rad - lng1Rad

	// Haversine 公式
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// 计算距离
	distance := earthRadius * c
	return distance
}

func (p *TPoint) GetGeoHash(precision uint) string {
	if p == nil {
		return ""
	}
	return p.Point.GetGeoHash(precision)
}

// DistanceTo 计算两个 TPoint 之间的距离（单位：米）
// 会将坐标统一转换为 WGS84 后再计算距离
func (p *TPoint) DistanceTo(other *TPoint) float64 {
	if p == nil || other == nil {
		return 0
	}

	// 将两个点都转换为 WGS84 坐标系
	p1 := p.Transform(coordinate.WGS84)
	p2 := other.Transform(coordinate.WGS84)

	return p1.Point.DistanceTo(&p2.Point)
}

func (p TPoint) Transform(t coordinate.TYPE) TPoint {
	lng, lat := coordinate.Transform(p.Point.Lng, p.Point.Lat, p.TYPE, t)
	return TPoint{
		TYPE:  t,
		Point: Point{Lng: lng, Lat: lat},
	}
}
