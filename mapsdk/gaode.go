package mapsdk

// nolint
type gaode struct {
	Host   string
	APIKey string
}

// nolint
func newGaode(apiKey string) *gaode {
	return &gaode{
		Host:   "https://api.map.gaode.com",
		APIKey: apiKey,
	}
}
