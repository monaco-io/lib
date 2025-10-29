package mapsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/monaco-io/lib/typing"
	"github.com/monaco-io/lib/xhttp"
)

type baidu struct {
	host string
	ak   string
}

func newBaidu(ak string) *baidu {
	if ak == "" {
		panic("Baidu AK is required")
	}
	return &baidu{
		host: "https://api.map.baidu.com",
		ak:   ak,
	}
}

type NativeDoResponse struct {
	URI  string
	Body json.RawMessage
}

func (b *baidu) NativeDo(uri string, params ...typing.KV[string, string]) (*NativeDoResponse, error) {
	values := url.Values{}
	for _, opt := range params {
		k, v := opt.Get()
		values.Set(k, v)
	}
	values.Set("ak", b.ak)
	values.Set("output", "json")
	response, err := xhttp.Do(context.Background(), b.host+uri,
		xhttp.URLRawQuery(values))
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}
	return &NativeDoResponse{
		URI:  response.URL.String(),
		Body: response.Body,
	}, nil
}

// https://lbs.baidu.com/faq/api?title=webapi/guide/webservice-placeapiV3/interfaceDocumentV3
func (b *baidu) SearchRegion(params SearchRegionParams, opts ...typing.KV[string, string]) (*Response[SearchPlaceData], error) {
	base := []typing.KV[string, string]{
		typing.NewKV("query", params.Keyword),
	}
	if params.Radius != 0 {
		base = append(base, typing.NewKV("radius", fmt.Sprintf("%d", params.Radius)))
		if params.Lat != 0 && params.Lng != 0 {
			base = append(base, typing.NewKV("location", params.GetPointString()))
		}
		opts = append(base, opts...)
		body, err := b.NativeDo("/place/v3/around", opts...)
		if err != nil {
			return nil, err
		}
		return unmarshal[*baiduSearchResponse3](body)
	}
	if params.Region != "" {
		base = append(base, typing.NewKV("region", params.Region))
		if params.Lat != 0 && params.Lng != 0 {
			base = append(base, typing.NewKV("center", params.GetPointString()))
		}
		opts = append(base, opts...)
		body, err := b.NativeDo("/place/v3/region", opts...)
		if err != nil {
			return nil, err
		}
		return unmarshal[*baiduSearchResponse3](body)
	} else { // 无区域时，使用附近检索
		if params.Lat != 0 && params.Lng != 0 {
			base = append(base, typing.NewKV("location", params.GetPointString()))
		}
		//query=银行&location=39.915,116.404&radius=2000&output=json&ak=您的密钥
		opts = append(base, opts...)
		body, err := b.NativeDo("/place/v2/search", opts...)
		if err != nil {
			return nil, err
		}
		return unmarshal[*baiduSearchResponse3](body)
	}
}

func (b *baidu) GetPlaceDetail(params GetPlaceDetailParams, opts ...typing.KV[string, string]) (*Response[[]PlaceDetailData], error) {
	opts = append(opts, typing.NewKV("uids", strings.Join(params.IDs, ",")))
	body, err := b.NativeDo("/place/v2/detail", opts...)
	if err != nil {
		return nil, err
	}
	return unmarshal[*baiduDetailResponse](body)
}

// 控制返回附近POI类型
// 以下内容需要 extensions_poi=1时才生效；
// 可以选择poi类型召回不同类型的poi，例如poi_types=酒店，如想召回多个POI类型数据，可以‘|’分割
// 例如poi_types=酒店|房地产 不添加该参数则默认召回全部POI分类数据。
// poi分类 https://lbsyun.baidu.com/index.php?title=open/poitags
func (b *baidu) GetReverseGeocoding(params GetReverseGeocodingParams, opts ...typing.KV[string, string]) (*Response[ReverseGeocodingData], error) {
	base := []typing.KV[string, string]{
		typing.NewKV("location", params.GetPointString()),
	}
	opts = append(base, opts...)
	body, err := b.NativeDo("/reverse_geocoding/v3", opts...)
	if err != nil {
		return nil, err
	}
	return unmarshal[*baiduReverseGeocodingResponse3](body)
}

func (b *baidu) GetTransitRoute(params GetTransitRouteParams, opts ...typing.KV[string, string]) (*Response[TransitRouteData], error) {
	base := []typing.KV[string, string]{
		typing.NewKV("origin", params.From.GetPointString()),
		typing.NewKV("destination", params.To.GetPointString()),
		typing.NewKV("steps_info", "1"),
	}
	opts = append(
		base, opts...)
	body, err := b.NativeDo("/direction/v2/transit", opts...)
	if err != nil {
		return nil, err
	}
	return unmarshal[*baiduTransitRouteDataResponse](body)
}

type IndustryClassification map[string][]string

func (p IndustryClassification) GetAllTypes1() []string {
	var types []string
	for k := range p {
		types = append(types, k)
	}
	return types
}

func (p IndustryClassification) GetAllTypes2() []string {
	var types []string
	for _, v := range p {
		types = append(types, v...)
	}
	return types
}

var BaiduIndustryClassification = IndustryClassification{
	"美食":    {"中餐厅", "外国餐厅", "小吃快餐店", "蛋糕甜品店", "咖啡厅", "茶座", "酒吧", "其他"},
	"酒店":    {"星级酒店", "快捷酒店", "公寓式酒店", "民宿", "其他"},
	"购物":    {"购物中心", "百货商场", "超市", "便利店", "家居建材", "家电数码", "商铺", "市场", "其他"},
	"生活服务":  {"通讯营业厅", "邮局", "物流公司", "售票处", "洗衣店", "图文快印店", "照相馆", "房产中介机构", "公用事业", "维修点", "家政服务", "殡葬服务", "彩票销售点", "宠物服务", "报刊亭", "公共厕所", "步骑行专用道驿站", "其他"},
	"丽人":    {"美容", "美发", "美甲", "美体", "其他"},
	"旅游景点":  {"公园", "动物园", "植物园", "游乐园", "博物馆", "水族馆", "海滨浴场", "文物古迹", "教堂", "风景区", "景点", "寺庙", "其他"},
	"休闲娱乐":  {"度假村", "农家院", "电影院", "ktv", "剧院", "歌舞厅", "网吧", "游戏场所", "洗浴按摩", "休闲广场", "其他"},
	"运动健身":  {"体育场馆", "极限运动场所", "健身中心", "其他"},
	"教育培训":  {"高等院校", "中学", "小学", "幼儿园", "成人教育", "亲子教育", "特殊教育学校", "留学中介机构", "科研机构", "培训机构", "图书馆", "科技馆", "其他"},
	"文化传媒":  {"新闻出版", "广播电视", "艺术团体", "美术馆", "展览馆", "文化宫", "其他"},
	"医疗":    {"综合医院", "专科医院", "诊所", "药店", "体检机构", "疗养院", "急救中心", "疾控中心", "医疗器械", "医疗保健", "核酸检测点", "新冠疫苗接种点", "风险点", "方舱医院", "发热门诊", "其他"},
	"汽车服务":  {"汽车销售", "汽车维修", "汽车美容", "汽车配件", "汽车租赁", "汽车检测场", "其他"},
	"交通设施":  {"飞机场", "火车站", "地铁站", "地铁线路", "长途汽车站", "公交车站", "港口", "停车场", "停车区", "停车位", "加油加气站", "服务区", "收费站", "桥", "充电站", "路侧停车位", "普通停车位", "接送点", "电动自行车充电站", "高速公路停车区", "其他"},
	"金融":    {"银行", "ATM", "信用社", "投资理财", "典当行", "其他"},
	"房地产":   {"写字楼", "住宅区", "宿舍", "内部楼栋", "其他"},
	"公司企业":  {"公司", "园区", "农林园艺", "厂矿", "其他"},
	"政府机构":  {"中央机构", "各级政府", "行政单位", "公检法机构", "涉外机构", "党派团体", "福利机构", "政治教育机构", "社会团体", "民主党派", "居民委员会", "其他"},
	"出入口":   {"高速公路出口", "高速公路入口", "机场出口", "机场入口", "车站出口", "车站入口", "门", "停车场出入口", "自行车高速出口", "自行车高速入口", "自行车高速出入口", "停车场出口", "停车场入口", "其他"},
	"自然地物":  {"岛屿", "山峰", "水系", "其他"},
	"行政地标":  {"省", "省级城市", "地级市", "区县", "商圈", "乡镇", "村庄", "其他"},
	"门址":    {"门址点", "其他"},
	"道路":    {"高速公路", "国道", "省道", "县道", "乡道", "城市快速路", "城市主干道", "城市次干道", "城市支路", "车渡线", "路口", "其他"},
	"铁路":    {"铁路", "地铁/轻轨", "磁悬浮列车", "有轨电车", "城际快轨", "其他"},
	"行政界线":  {"其他国家国界", "已定国界", "未定国界", "港澳界线", "南海范围线", "已定省界", "未定省界", "海岸线", "其他"},
	"其他线要素": {"桥梁", "隧道", "行政假想线", "水域假想线", "绿地假想线", "岛屿假想线", "疫情管控区", "其他"},
	"行政区划":  {"世界级", "国家级", "省级", "市级", "区县级", "热点区域", "建成区", "智能区域", "其他"},
	"水系":    {"双线河", "湖沼", "海洋", "其他"},
	"绿地":    {"绿地公园", "高尔夫球场", "岛", "绿化带", "机场", "机场道路", "其他"},
	"标注":    {"大洲标注", "大洋标注", "海域标注", "水系标注", "岛屿标注", "非水系标注", "其他"},
	"公交线路":  {"普通日行公交车", "地铁\\轻轨", "有轨电车", "机场巴士（前往机场）", "机场巴士（从机场返回）", "机场巴士（机场之间）", "旅游线路车", "夜班车", "轮渡", "快车", "慢车", "机场快轨（前往机场）", "机场快轨（从机场返回）", "机场轨道交通环路", "其他"},
	"电子眼":   {"限速电子眼", "应急车道电子眼", "公交车道电子眼", "外地车辆电子眼", "违章电子眼", "其他"},
}
