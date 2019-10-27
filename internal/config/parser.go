package config

import (
	"flag"
	"fmt"
	"strings"
)

var supportedAlgs string = "MD2,MD4,MD5,GOST,SHA1,SHA224,SHA256,SHA384,SHA512,RMD160,SHA3-224,SHA3-256,SHA3-384,SHA3-512,SHA512-224,SHA512-256,BLAKE2s256,BLAKE2b256,BLAKE2b384,BLAKE2b512,SHAKE128,SHAKE256,SM3"

func ParseConfig() *Config {
	c := new(Config)
	var algo *string
	c.LocalOnly = flag.Bool("localOnly", false, "Only calculate checksum locally \nif remote servers present in config app will ignore those servers")
	c.LocalOnly = flag.Bool("l", false, "Only calculate checksum locally \nif remote servers present in config app will ignore those servers")
	algo = flag.String("algorithm", "SHA512", fmt.Sprintf("Checksum algorithm, you can choose multiple by using MD5,SHA512... \nsupported algos: %s", supportedAlgs))
	algo = flag.String("a", "SHA512", fmt.Sprintf("Checksum algorithm, you can choose multiple by using MD5,SHA512... \nsupported algos: %s", supportedAlgs))
	c.LogLevel = flag.String("logLevel", "WARNING", "log level, supported: {TRACE,DEBUG,INFO,WARNING,ERROR}")
	c.Proxy = flag.String("proxy", "", "Proxy address to reach file or servers")
	c.Proxy = flag.String("p", "", "Proxy address to reach file or servers")
	c.RemoteOnly = flag.Bool("remoteOnly", false, "Only calculate checksum remotely, by default calculation will be done locally and remotely as well")
	c.RemoteOnly = flag.Bool("r", false, "Only calculate checksum remotely, by default calculation will be done locally and remotely as well")
	c.Validate = flag.Bool("validate", true, "Cross validates each calculated checksums, if any of them not matches doesnt download file with error output")
	c.Validate = flag.Bool("v", true, "Cross validates each calculated checksums, if any of checksum not matches doesnt download file with error output")
	c.Download = flag.Bool("download", true, "If all checksums are matches download file into current path, set false if you just want to retrieve file, if local allowed file still be downloaded")
	c.Download = flag.Bool("d", true, "If all checksums are matches download file into current path, set false if you just want to retrieve file, if local allowed file still be downloaded")
	c.Timeout = flag.Int("timeout", 60, "Timeout in secounds for each running calculation")
	c.Timeout = flag.Int("t", 60, "Timeout in secounds for each running calculation")
	c.All = flag.Bool("all", false, "Run all algorithms (MD5,SHA1 , SHA256 ...) for each running client")
	c.Key = flag.String("key", "", "Key for blake2 hashing")
	c.Key = flag.String("k", "", "Key for blake2 hashing")
	c.Supplier = flag.String("supplier", "go", "Algorithm supplier default is [GO] that core golang libraries used, if you want to use unix, win, mac default apps set to [OS], for openssl set [openssl] cloud providers support may vary")
	c.Supplier = flag.String("s", "go", "Algorithm supplier default is [GO] that core golang libraries used, if you want to use unix, win, mac default apps set to [OS], for openssl set [openssl] cloud providers support may vary")
	var empty string = ""
	c.File = &empty
	c.Cheksum = &empty

	upper := strings.ToUpper(*algo)
	flag.Parse()
	lower := strings.ToLower(*c.Supplier)
	c.Supplier = &lower
	c.Algorithm = strings.Split(upper, ",")
	args := flag.Args()
	if args != nil {
		if len(args) > 0 && args[0] != "" {
			c.File = &args[0]
		}
		if len(args) > 1 && args[1] != "" {
			c.Cheksum = &args[1]
		}

	}

	return c
}
