package weather

type (
	// GetCurrentWeatherResponse represents the response body when requesting the current weather
	// from the weather API.
	GetCurrentWeatherResponse struct {
		Location Location `json:"location"`
		Current  Weather  `json:"current"`
	}

	// The Location type represents a unique location.
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	}

	// The Weather type contains fields describing the weather.
	Weather struct {
		LastUpdatedEpoch int       `json:"last_updated_epoch"`
		LastUpdated      string    `json:"last_updated"`
		TempC            float64   `json:"temp_c"`
		TempF            float64   `json:"temp_f"`
		IsDay            float64   `json:"is_day"`
		Condition        Condition `json:"condition"`
		WindMph          float64   `json:"wind_mph"`
		WindKph          float64   `json:"wind_kph"`
		WindDegree       float64   `json:"wind_degree"`
		WindDir          string    `json:"wind_dir"`
		PressureMb       float64   `json:"pressure_mb"`
		PressureIn       float64   `json:"pressure_in"`
		PrecipMm         float64   `json:"precip_mm"`
		PrecipIn         float64   `json:"precip_in"`
		Humidity         float64   `json:"humidity"`
		Cloud            float64   `json:"cloud"`
		FeelslikeC       float64   `json:"feelslike_c"`
		FeelslikeF       float64   `json:"feelslike_f"`
		VisKm            float64   `json:"vis_km"`
		VisMiles         float64   `json:"vis_miles"`
		Uv               float64   `json:"uv"`
		GustMph          float64   `json:"gust_mph"`
		GustKph          float64   `json:"gust_kph"`
	}

	// The Condition type represents a summary of the current weather condition.
	Condition struct {
		Text string `json:"text"`
		Icon string `json:"icon"`
		Code int    `json:"code"`
	}
)
