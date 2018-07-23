package model

type DirectionResponse struct {
	Status            string    `json:"status"`
	GeocodedWaypoints []Geocode `json:"geocoded_waypoints"`
	Routes            []Route   `json:routes"`
}

type Route struct {
	Bound Bounds `json:"bounds"`
	Legs  []Leg  `json:"legs"`
}

type Geocode struct {
	status  string   `json:"geocoder_status"`
	PlaceID string   `json:"place_id"`
	Types   []string `json:"types"`
}

type Bounds struct {
	NE Northeast `json:"northeast"`
	SW Southwest `json:"southwest"`
}

type Leg struct {
	Distance struct {
		Text  string  `json:"text"`
		Value float64 `json:"value"`
	}
}

type Northeast struct {
	LatLong
}

type Southwest struct {
	LatLong
}

type LatLong struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"lng"`
}
