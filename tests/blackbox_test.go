package file_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

const (
	geturl              = "https://download.microsoft.com/download/8/b/4/8b4addd8-e957-4dea-bdb8-c4e00af5b94b/NDP1.1sp1-KB867460-X86.exe"
	help                = "-h"
	fileName            = "NDP1.1sp1-KB867460-X86.exe"
	all                 = "-all "
	multipleAlgs        = "-a MD5,SHA1,SHA512 "
	invalidServerConfig = "-sc /tmp/nothingFile "
	invalidAlgo         = "-a TestSum "
	invalidFile         = "/tmp/nothingFile"
	remoteOnly          = "-remoteOnly "
	localOnly           = "-localOnly "
	timeout             = "-t 1 -a RMD160 "
	invalidLib          = "-lib openGo "
	validLibGo          = "-a MD5 -lib Go "
	validLibOpenssl     = "-a MD5 -lib openssl "
	validLibOs          = "-a MD5 -lib os "
	invalidDir          = "-dir /tmp/nothingFile "
	validDir            = "-dir /tmp "
	keep                = "-keep "
	serve               = "-s " + validDir

	MD4        = "bb137fd4893ab9d85906257ede37dfaf"
	MD5        = "22e38a8a7d90c088064a0bbc882a69e5"
	SHA1       = "74a5b25d65a70b8ecd6a9c301a0aea10d8483a23"
	SHA224     = "18507f80722780ca477d7f10528ae28dd176f8d36cbce05a50cc7be0"
	SHA256     = "2c0a35409ff0873cfa28b70b8224e9aca2362241c1f0ed6f622fef8d4722fd9a"
	SHA384     = "c2372c71f93b5dc2a1c21c804bc74e27d82bfa45ee50fbc9037e713c156f1c591ffbe5e87f94022157906098916403b4"
	SHA512     = "bbe643f447f49636732b12d23a052d02681ad41f6920dc1038b073fa600f7589b378ed8e7de97e811543d93ae89ce52871a85ee58aa3b6aeaddc01bc1617ad85"
	RMD160     = "5bcec7ca2ff3b4b13db72cafb14c311a8fd281dd"
	SHA3224    = "b46b0b864c18029d4b90c4b16be7d9b96f6691d384c6aafb90fec059"
	SHA3256    = "5ea8f6ad9b197018aef9218cc6c372d7bff4e90eec8335cd884eeb50458e482f"
	SHA3384    = "3fee7161da83c2933d3fa6af675d74307bfee8c866a03a6e28775c441376db514359543e482d7fcb4b94ef7faf15dea7"
	SHA3512    = "28f136f06940f1325e4a081652159d9030ab384ec8e3201ff694b987643614f24b9625ad7212189ae706738be34d434629f54bb0a0303f549f7fa7b27160d409"
	SHA512224  = "63b2ffb0c5f1cd68abafba23997482b2087d486dcf60bec6fef7446d"
	SHA512256  = "7b44095feff471dee9366a2153dfe2654d70754c21b7e5204ed950cdf4a3f15a"
	BLAKE2s256 = "699c510d881bf3015dc027afbc32c8ae74b431342520a708fe9bf1760f4b26ef"
	BLAKE2b256 = "9d4d2f9ec65027b070163b3832070effa257f3f57d81ca226774625a90e27e0f"
	BLAKE2b384 = "afcfab976281f0b0a3c30213e9de1715274d811849dcaf155909aebef9d4facf7748b4141cb93133fa8d09587d81c7a8"
	BLAKE2b512 = "ec2de7c2c0c51dc1e016e1a4ef64b23830b088ca9aaf582a5531a6595a5339b61232168780b7e5d40e5bf92281a7272680650f518ae1ba1f33e7b7f8f3077219"
	SM3        = "05c4c5ac91721edec0bc332bcb3a5b5973291f433c8e0f69d23b5d5e354c37dd"
	SHAKE256   = "cfe847bc16528928ca7ab7bb0fb0ccd5443839335a8d5c5613ca777b4ab3859e"
	SHAKE128   = "9aac296a041c51b517223be2c31283d6"
)

func TestGet(t *testing.T) {

	execCommand(localOnly+geturl, fileName, true, t, SHA512)
}

func TestLibOpenSSL(t *testing.T) {

	commandStr := localOnly + validLibOpenssl + geturl
	execCommand(commandStr, fileName, true, t, MD5)

}

func TestLibOS(t *testing.T) {

	commandStr := localOnly + validLibOs + geturl
	execCommand(commandStr, fileName, true, t, MD5)

}

func TestLibGO(t *testing.T) {

	commandStr := localOnly + validLibGo + geturl
	execCommand(commandStr, fileName, true, t, MD5)

}

func TestAllAlgosGoLib(t *testing.T) {

	commandStr := all + localOnly + geturl
	execCommand(commandStr, fileName, true, t, MD4, MD5, SHA1, SHA224, SHA384, SHA256, SHA512, SHA3224, SHA3384, SHA3256, SHA3512, SHA512224, SHA512256, RMD160, BLAKE2s256, BLAKE2b256, BLAKE2b384, BLAKE2b512)

}
func TestAllAlgosOpenSSLLib(t *testing.T) {

	commandStr := all + localOnly + validLibOpenssl + geturl
	execCommand(commandStr, fileName, true, t, MD4, MD5, SHA1, SHA224, SHA384, SHA256, SHA512, SHA3224, SHA3384, SHA3256, SHA3512, SHA512224, SHA512256, RMD160, BLAKE2s256, BLAKE2b512, SHAKE128, SHAKE256, SM3)

}
func TestAllAlgosOSLib(t *testing.T) {

	commandStr := all + localOnly + validLibOs + geturl
	execCommand(commandStr, fileName, true, t, MD5, SHA1, SHA224, SHA384, SHA256, SHA512)

}
func TestMultipleAlgos(t *testing.T) {

	commandStr := multipleAlgs + localOnly + geturl
	execCommand(commandStr, fileName, true, t, MD5, SHA1, SHA512)

}
func TestValidDir(t *testing.T) {

	commandStr := localOnly + validDir + geturl
	execCommand(commandStr, "/tmp/"+fileName, true, t, SHA512)

}

func TestInValidDir(t *testing.T) {

	commandStr := localOnly + invalidDir + geturl
	execForError(commandStr, fileName, false, t, "Given -dir parameter")

}

func TestInValidAlgo(t *testing.T) {

	commandStr := localOnly + invalidAlgo + geturl
	execForError(commandStr, fileName, false, t, "Unrecognized algorithm")

}

func TestInValidLib(t *testing.T) {

	commandStr := localOnly + invalidLib + geturl
	execForError(commandStr, fileName, false, t, "Unrecognized library selection")
}

func TestInValidFile(t *testing.T) {

	commandStr := localOnly + invalidFile
	execForError(commandStr, fileName, false, t, "Given url")

}
func TestValidationFail(t *testing.T) {

	commandStr := localOnly + geturl + " " + MD5
	execForError(commandStr, fileName, false, t, "MISMATCH")

}

func TestValidation(t *testing.T) {

	commandStr := localOnly + geturl + " " + SHA512
	execCommand(commandStr, fileName, true, t, SHA512)

}

func TestKeep(t *testing.T) {
	commandStr := localOnly + keep + geturl + " " + MD5
	execForError(commandStr, fileName, true, t, "MISMATCH")

}

func TestRemoteOnly(t *testing.T) {
	commandStr := serve
	cmd := getCommand(commandStr)
	err := cmd.Start()
	defer killServer(cmd, t)
	if err != nil {
		t.Errorf("Can not start server instance! %s", err.Error())
	}
	commandStr = "-remoteOnly -a MD5 -sc servers.yml " + geturl
	execCommand(commandStr, fileName, false, t, MD5, "server1")

}

func TestServeLocalRemote(t *testing.T) {
	commandStr := serve
	cmd := getCommand(commandStr)
	err := cmd.Start()
	defer killServer(cmd, t)
	if err != nil {
		t.Errorf("Can not start server instance! %s", err.Error())
	}
	commandStr = "-a MD5 -sc servers.yml " + geturl
	execCommand(commandStr, fileName, true, t, MD5, "server1", "local")

}

func TestServeRemoteValidation(t *testing.T) {
	commandStr := serve
	cmd := getCommand(commandStr)
	err := cmd.Start()
	defer killServer(cmd, t)
	if err != nil {
		t.Errorf("Can not start server instance! %s", err.Error())
	}
	commandStr = "-a MD5 -sc servers.yml " + geturl + " " + MD5
	execCommand(commandStr, fileName, true, t, MD5, "server1", "local")
}

func TestServeRemoteValidationFail(t *testing.T) {
	commandStr := serve
	cmd := getCommand(commandStr)
	err := cmd.Start()
	defer killServer(cmd, t)
	if err != nil {
		t.Errorf("Can not start server instance! %s", err.Error())
	}
	commandStr = "-a MD5 -sc servers.yml " + geturl + " " + MD4
	execForError(commandStr, fileName, false, t, "SUSPENDED", "server1", "local")

}

func TestServeAlgoFail(t *testing.T) {
	commandStr := serve
	cmd := getCommand(commandStr)
	err := cmd.Start()
	defer killServer(cmd, t)
	if err != nil {
		t.Errorf("Can not start server instance! %s", err.Error())
	}
	commandStr = "-a MD5,SHA512 -sc servers.yml " + geturl + " " + MD4
	execForError(commandStr, fileName, false, t, "you can only run single algorithm")
}

func killServer(cmd *exec.Cmd, t *testing.T) {
	if cmd != nil && cmd.Process != nil {
		err := cmd.Process.Kill()
		if err != nil {
			t.Errorf("Can not kill process! %s", err.Error())
		}
	}
}

func execCommand(commandStr string, filename string, keep bool, t *testing.T, contains ...string) {
	cmd := getCommand(commandStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("An error occured while calling %s error: %s , \noutput: %s", commandStr, err.Error(), string(out))
	} else {
		strOut := string(out)
		t.Logf("Calling getsum %s returned success", commandStr)
		t.Logf("\n%s", strOut)
		for _, contain := range contains {
			if !strings.Contains(strOut, contain) {
				t.Errorf("Output doesnt contain sum %s", contain)
			}

		}

		if keep && !fileExist(filename) {
			t.Errorf("Command successfull but file %s not present!", filename)
		} else if !keep && fileExist(filename) {
			t.Errorf("-keep is false but file is still present!")
		}
		defer deleteFile(filename)
	}
}
func execForError(commandStr string, filename string, keep bool, t *testing.T, contains ...string) {
	cmd := getCommand(commandStr)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Errorf("An error expected while calling %s , \noutput: %s", commandStr, string(out))
	} else {
		strOut := string(out)
		t.Logf("Calling getsum %s returned error %s", commandStr, err.Error())
		t.Logf("\n%s", strOut)
		for _, contain := range contains {
			if !strings.Contains(strOut, contain) {
				t.Errorf("Output doesnt contain error %s", contain)
			}

		}

		defer deleteFile(filename)
		if !keep && fileExist(filename) {
			t.Errorf("Command returned errpr but file %s is present!", filename)
		} else if keep {
			if !fileExist(filename) {
				t.Errorf("Command returned error and keep param used but file %s is not present!", filename)

			} else {
				t.Logf("-keep param success!")
			}

		}
	}
}

func deleteFile(path string) {
	os.Remove(path)
}

func fileExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func getCommand(command string) *exec.Cmd {
	fields := strings.Split(command, " ")
	return exec.Command("getsum", fields...)
}
