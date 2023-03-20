package consul

type Nodes []string

type Config struct {
	Nodes        Nodes  `toml:"nodes"`
	Scheme       string `toml:"scheme"`
	ClientCert   string `toml:"client_cert"`
	ClientKey    string `toml:"client_key"`
	ClientCaKeys string `toml:"client_cakeys"`
	BasicAuth    bool   `toml:"basic_auth"`
	Username     string `toml:"username"`
	Password     string `toml:"password"`
	Token        string `toml:"token"`
}
