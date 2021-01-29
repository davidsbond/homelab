package minecraft

type (
	// ServerStatus represents the JSON response when querying the status of a Minecraft server.
	ServerStatus struct {
		IP       string  `json:"ip"`
		Port     int     `json:"port"`
		MOTD     MOTD    `json:"motd"`
		Players  Players `json:"players"`
		Version  string  `json:"version"`
		Online   bool    `json:"online"`
		Protocol int     `json:"protocol"`
	}

	// MOTD contains the "message of the day" data for a Minecraft server.
	MOTD struct {
		Raw   []string `json:"raw"`
		Clean []string `json:"clean"`
		HTML  []string `json:"html"`
	}

	// Players contains the current players on a Minecraft server.
	Players struct {
		Online int      `json:"online"`
		Max    int      `json:"max"`
		List   []string `json:"list"`
	}
)
