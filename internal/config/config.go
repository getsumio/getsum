package config

type Config struct {
	File         *string
	Remote       *bool
	LocalOnly    *bool
	Proxy        *string
	Algorithm    *string
	Cheksum      *string
	RemoteOnly   *bool
	OnlyChecksum *bool
	LogLevel     *string
	Validate     *bool
	Download     *bool
}
