package speedtest

type (
	// GetClientsResponse is the response DTO used when getting the speed test client info.
	GetClientsResponse struct {
		Clients []Client `xml:"client"`
	}

	// The Client type contains information about the client performing a speed test.
	Client struct {
		IP        string `xml:"ip,attr"`
		Latitude  string `xml:"lat,attr"`
		Longitude string `xml:"lon,attr"`
		ISP       string `xml:"isp,attr"`
	}

	// GetServersResponse is the response DTO used when getting test server info.
	GetServersResponse struct {
		Servers []Server `xml:"servers>server"`
	}

	// The Server type contains information about a server that can be used to perform
	// a speed test.
	Server struct {
		URL       string `xml:"url,attr"`
		Latitude  string `xml:"lat,attr"`
		Longitude string `xml:"lon,attr"`
		Name      string `xml:"name,attr"`
		Country   string `xml:"country,attr"`
		Sponsor   string `xml:"sponsor,attr"`
		ID        string `xml:"id,attr"`
		URL2      string `xml:"url2,attr"`
		Host      string `xml:"host,attr"`
	}
)
