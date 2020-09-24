package coronavirus

type (
	// GetDataResponse represents the response format from the government coronavirus API.
	GetDataResponse struct {
		Length       int `json:"length"`
		MaxPageLimit int `json:"maxPageLimit"`
		Data         []struct {
			Date      string `json:"date"`
			NewCases  int    `json:"newCases"`
			NewDeaths *int   `json:"newDeaths"` // This can be null, so we use a pointer.
		} `json:"data"`
		Pagination struct {
			Current  string      `json:"current"`
			Next     interface{} `json:"next"`
			Previous interface{} `json:"previous"`
			First    string      `json:"first"`
			Last     string      `json:"last"`
		} `json:"pagination"`
	}
)
