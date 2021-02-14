package ftp

type (
	// The Option type is a function used to configure the connection to the FTP server.
	Option func(c *config)
)

// WithCredentials will cause the client to authenticate on open using the provided username and password
// combination.
func WithCredentials(username, password string) Option {
	return func(c *config) {
		c.username = username
		c.password = password
	}
}
