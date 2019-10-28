package config

import (
	"flag"
	"fmt"
	"strings"
)

var supportedAlgs string = "MD2,MD4,MD5,GOST,SHA1,SHA224,SHA256,SHA384,SHA512,RMD160,SHA3-224,SHA3-256,SHA3-384,SHA3-512,SHA512-224,SHA512-256,BLAKE2s256,BLAKE2b256,BLAKE2b384,BLAKE2b512,SHAKE128,SHAKE256,SM3"

const (
	defaultServe      = false
	defaultListen     = "http://127.0.0.1"
	defaultPort       = 8088
	defaultLocalOnly  = false
	defaultAlgo       = "SHA512"
	defaultRemoteOnly = false
	defaultValidate   = true
	defaultDownload   = true
	defaultTimeout    = 60
	defaultKey        = ""
	defaultSupplier   = "go"
)

func ParseConfig() *Config {
	c := new(Config)
	var algo *string
	c.Serve = flag.Bool("serve", defaultServe, "Run in server mode default address http://127.0.0.1:8088 otherwise set -listen and -port params")
	flag.BoolVar(c.Serve, "s", defaultServe, "shorthand of -serve")
	c.Listen = flag.String("listen", defaultListen, "listen address only setted if -serve is true")
	flag.StringVar(c.Listen, "l", defaultListen, "shorthand of -listen")
	c.Port = flag.Int("port", defaultPort, "Listen port, only enabled if -serve is true")
	flag.IntVar(c.Port, "p", defaultPort, "shorthand of -port")
	c.LocalOnly = flag.Bool("localOnly", defaultLocalOnly, "Only calculate checksum locally \nif remote servers present in config app will ignore those servers")
	flag.BoolVar(c.LocalOnly, "local", defaultLocalOnly, "shorthand of -localOnly")
	algo = flag.String("algorithm", defaultAlgo, fmt.Sprintf("Checksum algorithm, you can choose multiple by using MD5,SHA512... \nsupported algos: %s", supportedAlgs))
	flag.StringVar(algo, "a", defaultAlgo, "shorthand of -algorithm")
	c.LogLevel = flag.String("logLevel", "WARNING", "log level, supported: {TRACE,DEBUG,INFO,WARNING,ERROR}")
	c.Proxy = flag.String("proxy", "", "Proxy address to reach file or servers")
	c.RemoteOnly = flag.Bool("remoteOnly", defaultRemoteOnly, "Only calculate checksum remotely, by default calculation will be done locally and remotely as well")
	flag.BoolVar(c.RemoteOnly, "r", defaultRemoteOnly, "shorthand of -remoteOnly")
	c.Validate = flag.Bool("validate", defaultValidate, "Cross validates each calculated checksums, if any of them not matches doesnt download file with error output")
	flag.BoolVar(c.Validate, "v", defaultValidate, "shorthand of -validate")
	c.Download = flag.Bool("download", defaultDownload, "If all checksums are matches download file into current path, set false if you just want to retrieve file, if local allowed file still be downloaded")
	flag.BoolVar(c.Download, "d", defaultDownload, "shorthand of -download")
	c.Timeout = flag.Int("timeout", defaultTimeout, "Timeout in secounds for each running calculation")
	flag.IntVar(c.Timeout, "t", defaultTimeout, "shorthand of -timeout")
	c.All = flag.Bool("all", false, "Run all algorithms (MD5,SHA1 , SHA256 ...) for each running client")
	c.Key = flag.String("key", defaultKey, "Key for blake2 hashing")
	flag.StringVar(c.Key, "k", defaultKey, "shorthand of -key")
	c.Dir = flag.String("dir", ".", "Default folder to save files, default is current folder")
	c.Supplier = flag.String("supplier", defaultSupplier, "Algorithm supplier default is [GO] that core golang libraries used, if you want to use unix, win, mac default apps set to [OS], for openssl set [openssl] cloud providers support may vary")
	flag.StringVar(c.Supplier, "su", defaultSupplier, "shorthand of -supplier")
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
