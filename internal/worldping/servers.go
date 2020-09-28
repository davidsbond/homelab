package worldping

import "fmt"

type (
	// The Server type describes a server that can be pinged at a certain
	// location.
	Server struct {
		Name string
		Code string
	}
)

// URL returns the pingable URL of the server.
func (s Server) URL() string {
	const serverFmt = "ftp.%s.debian.org"

	return fmt.Sprintf(serverFmt, s.Code)
}

// FTP servers that host the debian mirrors.
var servers = []Server{
	{
		Name: "Armenia",
		Code: "am",
	},
	{
		Name: "Australia",
		Code: "au",
	},
	{
		Name: "Austria",
		Code: "at",
	},
	{
		Name: "Belarus",
		Code: "by",
	},
	{
		Name: "Belgium",
		Code: "be",
	},
	{
		Name: "Brazil",
		Code: "br",
	},
	{
		Name: "Bulgaria",
		Code: "bg",
	},
	{
		Name: "Canada",
		Code: "ca",
	},
	{
		Name: "China",
		Code: "cn",
	},
	{
		Name: "Croatia",
		Code: "hr",
	},
	{
		Name: "Czech Republic",
		Code: "cz",
	},
	{
		Name: "Denmark",
		Code: "dk",
	},
	{
		Name: "El Salvador",
		Code: "sv",
	},
	{
		Name: "Estonia",
		Code: "ee",
	},
	{
		Name: "France",
		Code: "fr",
	},
	{
		Name: "Germany",
		Code: "de",
	},
	{
		Name: "Greece",
		Code: "hk",
	},
	{
		Name: "Hungary",
		Code: "hu",
	},
	{
		Name: "Italy",
		Code: "it",
	},
	{
		Name: "Japan",
		Code: "jp",
	},
	{
		Name: "Korea",
		Code: "kr",
	},
	{
		Name: "Lithuania",
		Code: "lt",
	},
	{
		Name: "Mexico",
		Code: "mx",
	},
	{
		Name: "Moldova",
		Code: "md",
	},
	{
		Name: "Netherlands",
		Code: "nl",
	},
	{
		Name: "New Caledonia",
		Code: "nc",
	},
	{
		Name: "New Zealand",
		Code: "nz",
	},
	{
		Name: "Norway",
		Code: "no",
	},
	{
		Name: "Poland",
		Code: "pl",
	},
	{
		Name: "Portugal",
		Code: "pt",
	},
	{
		Name: "Romania",
		Code: "ro",
	},
	{
		Name: "Russia",
		Code: "ru",
	},
	{
		Name: "Slovakia",
		Code: "sk",
	},
	{
		Name: "Slovenia",
		Code: "si",
	},
	{
		Name: "Spain",
		Code: "es",
	},
	{
		Name: "Sweden",
		Code: "fi",
	},
	{
		Name: "Switzerland",
		Code: "ch",
	},
	{
		Name: "Taiwan",
		Code: "tw",
	},
	{
		Name: "Turkey",
		Code: "tr",
	},
	{
		Name: "United Kingdom",
		Code: "uk",
	},
	{
		Name: "United States",
		Code: "us",
	},
}
