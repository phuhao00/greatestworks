package consul

import (
	"encoding/json"
	"errors"
	"net"
	"os/user"
	"strings"
)

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

func LoadJSONFromConsulKV(key string, cfg interface{}) {
	configKeyParameterValue := strings.SplitN(key, ":", 2)
	if len(configKeyParameterValue) < 2 {
	}
	configKey := configKeyParameterValue[1]

	kv := GetConsul().KV()
	kvPair, _, err := kv.Get(configKey, nil)
	if err != nil {
	}
	if kvPair == nil {
	}
	if len(kvPair.Value) == 0 {
	}
	if err = json.Unmarshal(kvPair.Value, &cfg); err != nil {
	}
}

func isPrivateIPv4(ip net.IP) bool {
	return ip != nil &&
		(ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168)
}

func GetPrivateIPv4() (string, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, a := range as {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}

		ip := ipnet.IP.To4()
		if isPrivateIPv4(ip) {
			return ip.String(), nil
		}
	}
	return "", errors.New("no private ip address")
}

func GetConsulConfigName() string {
	iPv4, err := GetPrivateIPv4()
	if err != nil {
		return ""
	}
	confName := "consul:" + iPv4 + "-" + GetUser() + "-" + "login.json"

	return confName
}

func GetUser() string {
	var userName string
	u, err := user.Current()
	if err != nil {
		userName = "unknow"
	} else {
		userName = u.Username
	}
	sl := strings.Split(userName, "\\")
	userName = sl[len(sl)-1]
	return userName
}
