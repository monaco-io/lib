package cnarea

import (
	"testing"

	"github.com/monaco-io/lib/typing/xjson"
)

func Test_GetProvinceList(t *testing.T) {
	t.Log(len(GetProvinceList()), xjson.MarshalIndentStringX(GetProvinceList()))
}

func Test_GetCityList(t *testing.T) {
	citys := GetCityList()
	t.Log(len(citys), xjson.MarshalIndentStringX(citys))
}

func Test_GetDistrictList(t *testing.T) {
	t.Log(len(GetDistrictList()), xjson.MarshalIndentStringX(GetDistrictList()))
}

func Test_GetAreaByCode(t *testing.T) {
	t.Log(xjson.MarshalIndentStringX(GetAreaByCodeX("110101"))) // 东城区
}
