package mapsdk

import (
	"encoding/json"
	"strings"

	"github.com/monaco-io/lib/typing"
	"github.com/monaco-io/lib/typing/xec"
	"github.com/monaco-io/lib/typing/xjson"
	"github.com/monaco-io/lib/typing/xstr"
	"github.com/samber/lo"
)

type (
	BaiduPOI struct {
		Name       string              `json:"name"`
		Location   BaiduPOI_Location   `json:"location"`
		Address    string              `json:"address"`
		Province   string              `json:"province"`
		City       string              `json:"city"`
		Area       string              `json:"area"`
		Town       string              `json:"town"`
		TownCode   int                 `json:"town_code"`
		Adcode     int                 `json:"adcode"`
		Status     string              `json:"status"`
		Telephone  string              `json:"telephone"`
		StreetID   string              `json:"street_id"`
		Detail     int                 `json:"detail"`
		UID        string              `json:"uid"`
		DetailInfo BaiduPOI_Detail     `json:"detail_info,omitempty"`
		Children   []BaiduPOI_Children `json:"children,omitempty"`
	}
	BaiduPOI_Children struct {
		UID              string            `json:"uid"`
		ShowName         string            `json:"show_name"`
		Name             string            `json:"name"`
		ClassifiedPOITag string            `json:"classified_poi_tag"`
		Location         BaiduPOI_Location `json:"location"`
		Address          string            `json:"address"`
	}
	BaiduPOI_Location struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}
	BaiduPOI_Detail struct {
		ClassifiedPOITag string            `json:"classified_poi_tag"`
		NewAlias         string            `json:"new_alias"`
		Type             string            `json:"type"`
		DetailURL        string            `json:"detail_url"`
		ShopHours        string            `json:"shop_hours"`
		Price            string            `json:"price"`
		Label            string            `json:"label"`
		OverallRating    string            `json:"overall_rating"`
		ImageNum         string            `json:"image_num"`
		CommentNum       string            `json:"comment_num"`
		NaviLocation     BaiduPOI_Location `json:"navi_location"`
		Brand            string            `json:"brand"`
		IndoorFloor      string            `json:"indoor_floor"`
		Ranking          string            `json:"ranking"`
		ParentID         string            `json:"parent_id"`
		Photos           []string          `json:"photos"`
		BestTime         string            `json:"best_time"`
		SugTime          string            `json:"sug_time"`
		Description      string            `json:"description"`
	}
	BaiduResponse struct {
		Status  int             `json:"status"`
		Message string          `json:"message"`
		Result  json.RawMessage `json:"result"`
	}
)

func (item *BaiduPOI) GetTags() []string {
	var tags []string
	tags = append(tags, item.DetailInfo.Label,
		item.DetailInfo.Type, item.DetailInfo.Brand,
		item.DetailInfo.NewAlias,
		item.DetailInfo.Ranking,
	)
	tags = append(tags, strings.Split(item.DetailInfo.ClassifiedPOITag, ";")...)
	return lo.Filter(tags, func(item string, index int) bool {
		return item != ""
	})
}

type baiduSearchResponse3 struct {
	Status     int        `json:"status"`
	Message    string     `json:"message"`
	Total      int        `json:"total"`
	ResultType string     `json:"result_type"`
	QueryType  string     `json:"query_type"`
	Results    []BaiduPOI `json:"results"`
}

type baiduDetailResponse struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Result  []BaiduPOI `json:"result"`
}

func (v *baiduDetailResponse) ToJSON() json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}

func (v *baiduReverseGeocodingResponse3) ToJSON() json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}

func (v *baiduSearchResponse3) ToJSON() json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}

func (v *BaiduResponse) ToJSON() json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}

type (
	baiduReverseGeocodingResponse3 struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Result  struct {
			Location         BaiduPOI_Location      `json:"location"`
			FormattedAddress string                 `json:"formatted_address"`
			AddressComponent BAIDU_AddressComponent `json:"addressComponent"`
			POIs             []BaiduPOI             `json:"pois"`
		} `json:"result"`
	}
	BAIDU_AddressComponent struct {
		Country      string `json:"country"`
		Province     string `json:"province"`
		City         string `json:"city"`
		District     string `json:"district"`
		Town         string `json:"town"`
		Street       string `json:"street"`
		StreetNumber string `json:"street_number"`
	}
)

func (d *BaiduResponse) ResponseDTO(uri string) *Response[*BaiduResponse] {
	return &Response[*BaiduResponse]{
		Source:   Baidu,
		Status:   xec.New(d.Status, d.Message),
		MetaURI:  uri,
		MetaData: d.ToJSON(),
		Data:     d,
	}
}

func (d *baiduSearchResponse3) ResponseDTO(uri string) *Response[SearchPlaceData] {
	data := SearchPlaceData{
		Locations: make([]Location, 0, len(d.Results)),
	}
	for _, item := range d.Results {
		data.Locations = append(data.Locations, Location{
			ID:           item.UID,
			Name:         item.Name,
			Address:      item.Address,
			Country:      "",
			Province:     item.Province,
			City:         item.City,
			Area:         item.Area,
			Street:       item.StreetID,
			Town:         item.Town,
			StreetNumber: "",
			Telephone:    item.Telephone,
			Tags:         item.GetTags(),
			Point:        Point{Lat: item.Location.Lat, Lng: item.Location.Lng},
			Detail: LocationDetail{
				Classification: strings.Split(item.DetailInfo.ClassifiedPOITag, ";"),
				Type:           item.DetailInfo.Type,
				ShopHours:      item.DetailInfo.ShopHours,
				Price:          item.DetailInfo.Price,
				Label:          strings.Split(item.DetailInfo.Label, ";"),
				Photos:         item.DetailInfo.Photos,
			},
			Extra: xjson.MarshalStringX(typing.Map[string, any]{"status": item.Status}),
		})
	}
	return &Response[SearchPlaceData]{
		Source:   Baidu,
		Status:   xec.New(d.Status, d.Message),
		MetaURI:  uri,
		MetaData: d.ToJSON(),
		Data:     data,
	}
}

func (d *baiduDetailResponse) ResponseDTO(uri string) *Response[[]PlaceDetailData] {
	var data []PlaceDetailData
	for _, item := range d.Result {
		data = append(data, PlaceDetailData{
			Location: Location{
				ID:       item.UID,
				Name:     item.Name,
				Address:  item.Address,
				Province: item.Province,
				City:     item.City,
				Area:     item.Area,
				Tags:     item.GetTags(),
				Point:    Point{Lat: item.Location.Lat, Lng: item.Location.Lng},
			},
		})
	}
	return &Response[[]PlaceDetailData]{
		Source:   Baidu,
		Status:   xec.New(d.Status, d.Message),
		MetaData: d.ToJSON(),
		MetaURI:  uri,
		Data:     data,
	}
}

func (d *baiduReverseGeocodingResponse3) ResponseDTO(uri string) *Response[ReverseGeocodingData] {
	var child []Location
	ac := d.Result.AddressComponent
	for _, item := range d.Result.POIs {
		child = append(child, Location{
			ID:           item.UID,
			Name:         item.Name,
			Address:      xstr.DefaultIfBlank(item.Address, d.Result.FormattedAddress),
			Country:      "",
			Province:     xstr.DefaultIfBlank(item.Province, ac.Province),
			City:         xstr.DefaultIfBlank(item.City, ac.City),
			Area:         xstr.DefaultIfBlank(item.Area, ac.District),
			Street:       "",
			Town:         item.Town,
			StreetNumber: "",
			Tags:         []string{},
			Point:        Point{Lat: item.Location.Lat, Lng: item.Location.Lng},
		})
	}
	return &Response[ReverseGeocodingData]{
		Source:   Baidu,
		Status:   xec.New(d.Status, d.Message),
		MetaData: d.ToJSON(),
		MetaURI:  uri,
		Data: ReverseGeocodingData{
			Location: Location{
				ID: "",
				// Name:     d.Result.AddressComponent.Name,
				Address:      d.Result.FormattedAddress,
				Country:      ac.Country,
				Province:     ac.Province,
				City:         ac.City,
				Area:         ac.District,
				Town:         ac.Town,
				Street:       ac.Street,
				StreetNumber: ac.StreetNumber,
				Tags:         []string{},
				Point:        Point{Lat: d.Result.Location.Lat, Lng: d.Result.Location.Lng},
			},
			Child: child,
			Extra: "",
		},
	}
}
