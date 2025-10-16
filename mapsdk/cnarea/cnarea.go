package cnarea

import (
	"fmt"
	"regexp"
	"strings"

	_ "embed"

	"github.com/monaco-io/lib/typing/xstr"
	"github.com/samber/lo"
)

//go:embed sql/2023_insert.sql
var areasql string

var areas struct {
	list          []Area
	levelGroup    map[uint][]Area
	areaCodeGroup map[string]Area
}

type Area struct {
	Level      uint    `json:"level" db:"level"`
	ParentCode string  `json:"parent_code" db:"parent_code"`
	AreaCode   string  `json:"area_code" db:"area_code"`
	ZipCode    string  `json:"zip_code" db:"zip_code"`
	CityCode   string  `json:"city_code" db:"city_code"`
	Name       string  `json:"name" db:"name"`
	ShortName  string  `json:"short_name" db:"short_name"`
	MergerName string  `json:"merger_name" db:"merger_name"`
	Pinyin     string  `json:"pinyin" db:"pinyin"`
	Lng        float64 `json:"lng" db:"lng"`
	Lat        float64 `json:"lat" db:"lat"`
}

func load() []Area {
	line := strings.Split(strings.Trim(areasql, "\n"), "\n")
	// `level`, `parent_code`, `area_code`, `zip_code`, `city_code`, `name`, `short_name`, `merger_name`, `pinyin`, `lng`, `lat`
	areas := make([]Area, 0, len(line))

	// Regular expression to match VALUES clause
	re := regexp.MustCompile(`VALUES \s*\(([^)]+)\);`)

	for _, v := range line {

		matches := re.FindStringSubmatch(v)
		if len(matches) < 2 {
			panic(fmt.Errorf("failed to match VALUES in line: %s", v))
		}

		values := strings.Split(matches[1], ", ")
		if len(values) != 11 {
			panic(fmt.Errorf("unexpected number of values: %d in line: %s", len(values), v))
		}
		areas = append(areas, Area{
			Level:      xstr.ToIntegerX[uint](values[0]),
			ParentCode: strings.Trim(values[1], "'"),
			AreaCode:   strings.Trim(values[2], "'"),
			ZipCode:    strings.Trim(values[3], "'"),
			CityCode:   strings.Trim(values[4], "'"),
			Name:       strings.Trim(values[5], "'"),
			ShortName:  strings.Trim(values[6], "'"),
			MergerName: strings.Trim(values[7], "'"),
			Pinyin:     strings.Trim(values[8], "'"),
			Lng:        xstr.ToFloatX[float64](values[9]),
			Lat:        xstr.ToFloatX[float64](values[10]),
		})
	}
	return areas
}

func init() {
	areas.list = load()
	areas.levelGroup = lo.GroupBy(areas.list, func(item Area) uint {
		return item.Level
	})
	areas.areaCodeGroup = lo.SliceToMap(areas.list, func(item Area) (string, Area) {
		return item.AreaCode, item
	})
}

// 省
func GetProvinceList() []Area {
	return areas.levelGroup[1]
}

// 市
func GetCityList() []Area {
	return areas.levelGroup[2]
}

// 区
func GetDistrictList() []Area {
	return areas.levelGroup[3]
}

func GetAreaByCode(areaCode string) (Area, bool) {
	a, ok := areas.areaCodeGroup[areaCode]
	return a, ok
}

func GetAreaByCodeX(areaCode string) Area {
	return areas.areaCodeGroup[areaCode]
}
