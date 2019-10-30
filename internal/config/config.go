package config

//config dto
type Config struct {
	File         *string       `json:"file"`
	LocalOnly    *bool         `json:"local_only"`
	Proxy        *string       `json:"proxy"`
	Algorithm    []string      `json:"algorithm"`
	Cheksum      *string       `json:"cheksum"`
	RemoteOnly   *bool         `json:"remote_only"`
	OnlyChecksum *bool         `json:"only_checksum"`
	LogLevel     *string       `json:"log_level"`
	Timeout      *int          `json:"timeout"`
	All          *bool         `json:"all"`
	Key          *string       `json:"key"`
	Supplier     *string       `json:"supplier"`
	OnFailure    *string       `json:"on_failure"`
	Serve        *bool         `json:"serve"`
	Listen       *string       `json:"listen"`
	Port         *int          `json:"port"`
	Servers      ServerConfigs `json:"servers"`
	Dir          *string       `json:"dir"`
	TLSKey       *string       `json:"tls_key"`
	TLSCert      *string       `json:"tls_cert"`
	ServerConfig *string       `json:"server_config"`
	Keep         *bool         `json:"keep"`
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
