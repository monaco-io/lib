package mapsdk

// nolint
type tencent struct {
	Host   string
	APIKey string
}

// nolint
func newTencent(apiKey string) *tencent {
	return &tencent{
		Host:   "https://api.map.tencent.com",
		APIKey: apiKey,
	}
}
