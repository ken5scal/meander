package meander

var APIKey string

type Place struct {
	*googleGeometry `json:"geometry"` // Express Nested API data
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	Photos   []*googlePhoto `json:"photos"`
	Vicinity string `json:"vicinity"`
}

type googleResponse struct {
	Results [] *Place `json:"results"`
}

type googleGeometry struct {
	*googleLocation `json:"location"`
}

type googleLocation struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type googlePhoto struct {
	PhotoRef string `json:"photo_reference"`
	URL      string `json:"url"`
}

// Use this to return data from outside request
func (p *Place) Public() interface{} {
	return map[string]interface{}{
		"name": p.Name,
		"icon": p.Icon,
		"photos": p.Photos,
		"vicinity": p.Vicinity,
		"lat": p.Lat,
		"lng": p.Lng,
	}
}