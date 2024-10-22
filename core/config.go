package core

type Config struct {
	CIDRv4  string `toml:"CIDRv4"`
	TunName string `toml:"tunName"`

	ServiceAddr string `toml:"serviceAddr"`
	ServerMode  bool   `toml:"serverMode"`

	HttpService HttpService `toml:"httpService"`
	Extra       Extra       `toml:"extra"`
}

type Extra struct {
	ObfName bool `toml:"obfName"`
}

type HttpService struct {
	Enable     bool   `toml:"enable"`
	Port       int    `toml:"port"`
	ServiceKey string `toml:"key"`
}
