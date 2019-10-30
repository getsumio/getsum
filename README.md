# getsum : Tool for validating and calculating checksums

***getsum*** calculates and validates checksum of files remotely or locally. According to user choice local downloads can be prevented if checksum mismatch. You can also run ***getsum*** in listen mode so you can run remotely deploy on your server or cloud provider of your choice.

 [![Watch the full record](docs/main.gif)](https://asciinema.org/a/ovpGNqNS56qlrKevUllOks1qT)
 
**Installation**

**How to run**
```
getsum https://speed.hetzner.de/100MB.bin
getsum /tmp/path/to/file
getsum -a MD5,SHA512 https://speed.hetzner.de/100MB.bin cf1a31c3acf3a1c3f2a13cfa13
getsum -remoteOnly https://speed.hetzner.de/100MB.bin cf1a31c3acf3a1c3f2a13cfa13
``` 
**Features**
* Run in server mode to listen requests from another getsum client
* Multiple libraries/application support for calculating checksums
* Calculate checksums by using multiple/all supported algorithms at once
* Download file from internet and calculate checksum
* Prevent download if checksum doesnt match
* Local/Remote only calculations
* Proxy/TLS support on listen mode

**Selecting library/applications**

For checksum calculations core Golang libraries will be used as default. If you have installed openssl set '-lib openssl'. If you want to use applications from operating system set '-lib os'.

```
getsum -a -lib openssl https://speed.hetzner.de/100MB.bin
``` 

[![Watch the full record](docs/libs.gif)](https://asciinema.org/a/sy0OSLL8IWUOED2DFk1yFLiOB)


 
**Running Multiple Algorithms** 

Default algorithm is ***SHA512***. Use '-a' parameter to specify different algorithm. '-all' runs all algorithms at (if selected library doesnt support some of them only supported ones will run)

```
getsum -a MD5,SHA512,SHA1 https://speed.hetzner.de/100MB.bin
getsum -all /tmp/path/to/file
``` 

[![Watch the full record](docs/multiple.gif)](https://asciinema.org/a/nejfc4N0vLJhkxpqikEfHIBCe)
[![Watch the full record](docs/all.gif)](https://asciinema.org/a/KA4sT6xTNN9iTzKHJhdgnybrB)

**Running in listen mode**
Running in serve mode param is '-s' default listen address is 127.0.0.1 and port is 8088 
```
getsum -s 
getsum -s -l 0.0.0.0 -p 9099
getsum -s -l 0.0.0.0 -p 9099 -tlskey /tmp/tlskeyfile -tlscert /tmptlscertfile
``` 
[![Watch the full record](docs/server.gif)](https://asciinema.org/a/KA4sT6xTNN9iTzKHJhdgnybrB)

In case of 'os' selected:
* For Linux/Mac commands will be called: md5sum,sha1sum,sha224sum,sha256sum,sha384sum,sha512sum
* For Windows : certUtil will be called

Supported Algorithms:
* Windows: MD2, MD4, MD5, SHA1, SHA224, SHA256, SHA384, SHA512
* Linux/MAC: MD5, SHA1, SHA224, SHA256, SHA384, SHA512
* GO: MD4,MD5,SHA1,SHA224,SHA256,SHA384,SHA512,RMD160,SHA3_224,SHA3_256,SHA3_384,SHA3_512,SHA512_224,SHA512_256,BLAKE2s256,BLAKE2b256,BLAKE2b384,BLAKE2b512
* OPENSSL: MD4,MD5,SHA1,SHA224,SHA256,SHA384,SHA512,RMD160,SHA3_224,SHA3_256,SHA3_384,SHA3_512,SHA512_224,SHA512_256,BLAKE2s256,BLAKE2b512,SHAKE128,SHAKE256,SM3
**
