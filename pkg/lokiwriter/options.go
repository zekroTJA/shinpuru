package lokiwriter

// Options holds connection and authentication
// information for the loki connection.
type Options struct {
	Address  string            `json:"address"`
	Username string            `json:"username"`
	Password string            `json:"password"`
	Labels   map[string]string `json:"labels"`
}
