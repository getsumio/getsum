package config

type Config struct {
	File         *string        `json:"file"`
	Remote       *bool          `json:"remote"`
	LocalOnly    *bool          `json:"local_only"`
	Proxy        *string        `json:"proxy"`
	Algorithm    []string       `json:"algorithm"`
	Cheksum      *string        `json:"cheksum"`
	RemoteOnly   *bool          `json:"remote_only"`
	OnlyChecksum *bool          `json:"only_checksum"`
	LogLevel     *string        `json:"log_level"`
	Validate     *bool          `json:"validate"`
	Download     *bool          `json:"download"`
	Timeout      *int           `json:"timeout"`
	All          *bool          `json:"all"`
	Key          *string        `json:"key"`
	Supplier     *string        `json:"supplier"`
	OnFailure    *string        `json:"on_failure"`
	Serve        *bool          `json:"serve"`
	Listen       *string        `json:"listen"`
	Port         *int           `json:"port"`
	Servers      []ServerConfig `json:"servers"`
	Dir          *string        `json:"dir"`
	TLSKey       *string        `json:"tls_key"`
	TLSCert      *string        `json:"tls_cert"`
}

type ServerConfig struct {
	Name          string
	ListenAddress string
	Port          int
}
