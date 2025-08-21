package mapsdk

import "github.com/monaco-io/lib/mapsdk/geohash"

type TPoint struct {
	Point          `json:"point"`
	CoordinateType `json:"type"`
}

type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func (p *Point) GetGeoHash(precision uint) string {
	if p != nil {
		return geohash.EncodeWithPrecision(p.Lat, p.Lng, precision)
	}
	return ""
}

func (p *TPoint) GetGeoHash(precision uint) string {
	if p == nil {
		return ""
	}
	return p.Point.GetGeoHash(precision)
}

func (p TPoint) Transform(t CoordinateType) TPoint {
	return TPoint{
		CoordinateType: t,
		Point:          coordinateConverter.Transform(p.Point, p.CoordinateType, t),
	}
}
