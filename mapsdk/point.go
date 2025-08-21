package mapsdk

import (
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

func (p TPoint) Transform(t coordinate.TYPE) TPoint {
	lng, lat := coordinate.Transform(p.Point.Lng, p.Point.Lat, p.TYPE, t)
	return TPoint{
		TYPE:  t,
		Point: Point{Lng: lng, Lat: lat},
	}
}
