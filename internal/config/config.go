package config

//config dto
type Config struct {
	File               *string `json:"file"`
	LocalOnly          *bool
	Proxy              *string  `json:"proxy"`
	Algorithm          []string `json:"algorithm"`
	Cheksum            *string  `json:"cheksum"`
	RemoteOnly         *bool
	LogLevel           *string
	Timeout            *int    `json:"timeout,string"`
	All                *bool   `json:"all"`
	Key                *string `json:"key"`
	Supplier           *string `json:"supplier"`
	Serve              *bool
	Listen             *string
	Port               *int
	Servers            ServerConfigs
	Dir                *string
	TLSKey             *string
	TLSCert            *string
	ServerConfig       *string
	Keep               *bool
	InsecureSkipVerify *bool
	Quite              *bool
}

//this is for collecting server info from yaml files
type ServerConfig struct {
	Name    string
	Address string
}

//wrapper for servers
type ServerConfigs struct {
	Servers []ServerConfig
}
