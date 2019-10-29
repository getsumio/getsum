package validation

import (
	"errors"
	"os"

	. "github.com/getsumio/getsum/internal/config"
	"github.com/getsumio/getsum/internal/supplier"
)

var supportedAlgs string = "MD2,MD4,MD5,GOST,SHA1,SHA224,SHA256,SHA384,SHA512,RMD160,SHA3-224,SHA3-256,SHA3-384,SHA3-512,SHA512-224,SHA512-256,BLAKE2s256,BLAKE2b256,BLAKE2b384,BLAKE2b512,SHAKE128,SHAKE256,SM3"

func ValidateConfig(config *Config, onPremise bool) error {
	if *config.File == "" {
		return errors.New("No file path/url provided, example usage: getsum /tmp/file ")
	}
	if *config.RemoteOnly {
		if *config.LocalOnly {
			return errors.New("You can not set -localOnly and -remoteOnly at the same time")
		}
		if len(config.Servers.Servers) < 1 {
			return errors.New("No server recognized, create a ~/.getsum/servers.yml file or use -serverconfig parameter, for content of yaml file check documentation")
		}
	}
	if *config.All && (*config.RemoteOnly || len(config.Servers.Servers) > 1) && !*config.LocalOnly {
		return errors.New("On remote servers you can only run single algorithm set -localOnly or example usage: getsum -a MD5 /tmp/file")
	}
	if len(config.Algorithm) != 1 && (*config.RemoteOnly || len(config.Servers.Servers) > 1) && !*config.LocalOnly {
		return errors.New("On remote servers you can only run single algorithm set -localOnly or example usage: getsum -a MD5 /tmp/file")
	}
	if *config.Timeout < 1 {
		return errors.New("Invalid timeout value")
	}
	if *config.Supplier != "os" && *config.Supplier != "openssl" && *config.Supplier != "go" {
		return errors.New("Unrecognized library selection [" + *config.Supplier + "], supported: os , go , openssl ")
	}
	for _, a := range config.Algorithm {
		if supplier.ValueOf(&a) == 127 {
			return errors.New("Unrecognized algorithm, [" + a + "] supported types: " + supportedAlgs)
		}
	}
	if *config.Dir != "" {
		fi, err := os.Stat(*config.Dir)
		if os.IsNotExist(err) {
			return errors.New("Given -dir parameter " + *config.Dir + "doesnt exist")
		}
		if !fi.IsDir() {
			return errors.New("Given -dir parameter is not a directory!")
		}
	}
	if *config.TLSKey != "" && *config.TLSCert == "" {
		return errors.New("You specified -tlskey but not -tlscert, both parameter required for https/tls mode")
	}
	if *config.TLSKey == "" && *config.TLSCert != "" {
		return errors.New("You specified -tlscert but not -tlskey, both parameter required for https/tls mode")
	}

	return nil
}
