package config

import (
	"flag"
)

func ParseConfig() *Config {
	c := new(Config)
	c.LocalOnly = flag.Bool("localOnly", false, "Only calculate checksum locally if remote servers present in config it will use those servers also local resources as well")
	c.LocalOnly = flag.Bool("l", false, "Only calculate checksum locally, if remote servers present in config it will use those servers also local resources as well")
	c.Algorithm = flag.String("algo", "SHA512", "Checksum algorithm, supported: {MD5,SHA-0,SHA-1,SHA256,SHA384,SHA512,SHA-3}")
	c.Algorithm = flag.String("a", "SHA512", "Checksum algorithm, supported: {MD5,SHA-0,SHA-1,SHA256,SHA384,SHA512,SHA-3}")
	c.LogLevel = flag.String("logLevel", "WARNING", "log level, supported: {TRACE,DEBUG,INFO,WARNING,ERROR}")
	c.Proxy = flag.String("proxy", "", "Proxy address to reach file or servers")
	c.Proxy = flag.String("p", "", "Proxy address to reach file or servers")
	c.RemoteOnly = flag.Bool("remoteOnly", false, "Only calculate checksum remotely, by default calculation will be done locally and remotely as well")
	c.RemoteOnly = flag.Bool("r", false, "Only calculate checksum remotely, by default calculation will be done locally and remotely as well")
	c.Validate = flag.Bool("validate", true, "Cross validates each calculated checksums, if any of them not matches doesnt download file with error output")
	c.Validate = flag.Bool("v", true, "Cross validates each calculated checksums, if any of checksum not matches doesnt download file with error output")
	c.Download = flag.Bool("download", true, "If all checksums are matches download file into current path, set false if you just want to retrieve file, if local allowed file still be downloaded")
	c.Download = flag.Bool("d", true, "If all checksums are matches download file into current path, set false if you just want to retrieve file, if local allowed file still be downloaded")

	flag.Parse()
	return c
}
