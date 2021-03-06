# getsum : Tool for validating and calculating checksums

***getsum*** calculates and validates checksum of files remotely or locally. According to user choice, if their checksum mismatch local downloads can be prevented . You can also run application in listen mode, so you can remotely deploy on your server or cloud provider. Then you can use another getsum as client on host pc (Please see also below docker image section) or browser addon. I get the idea from https://blog.linuxmint.com/?p=2994 so I thought it would be great fit for people who host binaries as well as users to validate their checksum. In validation mode if remote servers are present, then application first calculates checksum on remote servers and if there is a match it will download the file and run another calculation locally.

 [![Watch the full record](docs/main.gif)](https://asciinema.org/a/ovpGNqNS56qlrKevUllOks1qT)
 
**Installation**

 Current binaries are stored on [release page](https://github.com/getsumio/getsum/releases). Please consider application only tested on Fedora 30. 
 
 ```
 cd /location/to/store
 tar xzvf getsum-linux-amd64-XXXX.tar.gz
 cd builds/linux/amd64/
 
 ./getsum -h
 ```
 

 add binary location to /etc/profile or ~/.bashrc or if you have alternatives installed:
 ```
 alternatives --install /usr/bin/getsum getsum /location/to/store/getsum 0
 ```

**How to run**

Run 'getsum -h' for all parameters
```
getsum https://some.server.address/binary
getsum /tmp/path/to/file
getsum -a MD5,SHA512 https://some.server.address/binary cf1a31c3acf3a1c3f2a13cfa13
getsum -remoteOnly https://some.server.address/binary cf1a31c3acf3a1c3f2a13cfa13
``` 
**Features**

* Run in server mode to listen requests from another getsum client
* Multiple libraries/application support for calculating checksums
* Calculate checksums by using multiple/all supported algorithms at once
* Download file from internet and calculate checksum
* Prevent download if checksum doesnt match
* Local/Remote only calculations
* Proxy support
* TLS support on listen mode or reaching files behind untrusted certificate
* Browser addons

**Selecting library/applications**

For checksum calculations core Golang libraries will be used as default. If you have openssl installed set *-lib openssl* or if you want to use applications from operating system then set *-lib os*.

```
getsum -a MD5 -lib openssl https://some.server.address/binary
getsum -a MD5 -lib go https://some.server.address/binary
getsum -a MD5 -lib os https://some.server.address/binary
``` 

[![Watch the full record](docs/libs.gif)](https://asciinema.org/a/sy0OSLL8IWUOED2DFk1yFLiOB)


 
**Running Multiple Algorithms** 

Default algorithm is ***SHA512***. Use *-a* parameter to specify different algorithms. Algorithms are comma separated. *-all* runs all algorithms at once (if selected library doesnt support some of them, only supported ones will run)

```
getsum -a MD5,SHA512,SHA1 https://some.server.address/binary
getsum -all /tmp/path/to/file
``` 

[![Watch the full record](docs/multiple.gif)](https://asciinema.org/a/nejfc4N0vLJhkxpqikEfHIBCe)
[![Watch the full record](docs/all.gif)](https://asciinema.org/a/KA4sT6xTNN9iTzKHJhdgnybrB)

**Validation**
 If another checksum provided application will compare generated one with the given one. If there is mismatch file will be removed from host pc. Use ***-keep*** parameter if you want to keep file even there is mismatch. 
 ```
getsum -a MD5 -keep https://some.server.address/binary cf1a31c3acf3a1c3f2a13cfa13
getsum -remoteOnly https://some.server.address/binary cf1a31c3acf3a1c3f2a13cfa13
getsum -localOnly https://some.server.address/binary cf1a31c3acf3a1c3f2a13cfa13
``` 

**Running in serve mode**

Running in serve mode param is *-s* default listen address is *127.0.0.1* and port is *8088*. In serve mode files are not stored that they are removed after calculation. Set *-dir* param to change save folder, default is current location. There is no authentication method provided by this application, you need to handle it if you are planning to run servers in public.
```
getsum -s 
getsum -s -l 0.0.0.0 -p 9099
getsum -s -l 0.0.0.0 -p 9099 -tk /tmp/tlskeyfile -tc /tmp/tlscertfile
``` 
**Updating client to run with remote servers**

Create a config file at **$HOME/.getsum/servers.yml** with addresses of your servers, i.e.:
```
servers:
  - name: server1
    address: http://127.0.0.1:8088
  - name: server2
    address: http://127.0.0.1:8089
  - name: server3
    address: http://127.0.0.1:8090
```
"**servers,name and address" field names should be same** you just need to update values. 
Also use *-serverconfig* parameter for custom config location:
```
getsum -serverconfig /tmp/servers.yml /path/to/file
getsum -sc /tmp/servers.yml /path/to/file
``` 
A quick 3 server 1 client example:
 ```
 cd /location/to/store
 tar xzvf getsum-linux-amd64-XXXX.tar.gz
 cd builds/linux/amd64/
 mkdir {tmp{1,2,3},~/.getsum}
 
 cat > ~/.getsum/servers.yml << EOF
 servers:
   - name: server1
     address: http://127.0.0.1:8088
   - name: server2
     address: http://127.0.0.1:8089
   - name: server3
     address: http://127.0.0.1:8090
 EOF

 
./getsum -s -dir ./tmp1 > ./tmp1/server1.log 2>&1 &
./getsum -s -p 8089 -dir ./tmp2 > ./tmp2/server2.log 2>&1 &
./getsum -s -p 8090 -dir ./tmp3 > ./tmp3/server3.log 2>&1 &

./getsum -a MD5  -lib openssl https://download.microsoft.com/download/8/b/4/8b4addd8-e957-4dea-bdb8-c4e00af5b94b/NDP1.1sp1-KB867460-X86.exe 22e38a8a7d90c088064a0bbc882a69e5asd

./getsum -a MD5  -lib openssl https://download.microsoft.com/download/8/b/4/8b4addd8-e957-4dea-bdb8-c4e00af5b94b/NDP1.1sp1-KB867460-X86.exe 22e38a8a7d90c088064a0bbc882a69e5

 
 killall getsum
 rm -Rf ./tmp{1,2,3}
 
 ``` 

[![Watch the full record](docs/server.gif)](https://asciinema.org/a/KA4sT6xTNN9iTzKHJhdgnybrB)


**Docker image**
```
docker pull getsum/getsum
docker run -p127.0.0.1:8088:8088 getsum/getsum
```

Running docker in tls mode (assuming your key is server.key and cert is server.crt):
```
docker pull getsum/getsum
docker create -p127.0.0.1:8088:8088 -e tlskey='server.key' -e tlscert='server.crt' --name getsum getsum/getsum
docker cp /path/to/server.key getsum:/app/
docker cp /path/to/server.crt getsum:/app/
docker start -i getsum

//then on client machine

 cat > ~/.getsum/servers.yml << EOF
 servers:
   - name: server1
     address: https://127.0.0.1:8088
 EOF
//on client use -skipVerify in case of self signed cert
getsum -skipVerify -a MD5  -lib openssl https://download.microsoft.com/download/8/b/4/8b4addd8-e957-4dea-bdb8-c4e00af5b94b/NDP1.1sp1-KB867460-X86.exe 22e38a8a7d90c088064a0bbc882a69e5


```

**Browser Addons**

Chrome: https://chrome.google.com/webstore/detail/getsum/mbkjjcfdhbhpjmhpkkligkaifjmkakge

Firefox: https://addons.mozilla.org/en-US/firefox/addon/getsum/?src=search

Then start your server, extensions also allows you to use 127.0.0.1 / localhost ports. i.e. 

```
docker pull getsum/getsum
docker run -p127.0.0.1:8088:8088 getsum/getsum
```
 
on Firefox: about:addons -> then click extenstions -> select getsum and click preferences tab

on Chrome: On right top click Getsum icon -> select options

i.e.:

* hostname: http://127.0.0.1:8088
* library: openssl
* keep proxy empty (proxy for telling server to use it for reaching file)
* set timeout (server will use timeout to download and calculate)
* and save and restart browser

Then on any page  right click on a download link and select GetSum then select desired algorithm.

If you want to validate checksum of download file then first select text on the page and right click on download link and select Getsum->algorithm so extension will check if any selected text exist and if it is valid it will use it for comparation.

Here is a quick video:

https://www.youtube.com/watch?v=7f2hMyI38Lo .

 Addons just validates checksum so they are informative, even it is valid still you need to download file. If you use self signed/untrusted certificate on server still you need to convince browser to allow this certificate

  
  
**Serverless support**
 I really wanted to add native lambda, cloud functions support for different providers but each provider has their own limits i.e. 200mb storage space or 2GB memory, so its currently postponed.
 
 **Issues**
 Application tested only on linux. If you had issues please raise here. 
 
 **How to support**
  Code review, pull requests, raise issues, promote :) 

**In case of 'os' selected**:
below commands will be called:
* For ***Linux/Mac*** :  *md5sum,sha1sum,sha224sum,sha256sum,sha384sum,sha512sum*
* For ***Windows*** : *certUtil* will be called 

**Supported Algorithms**:
* ***Windows***: MD2, MD4, MD5, SHA1, SHA256, SHA384, SHA512
* ***Linux/MAC***: MD5, SHA1, SHA224, SHA256, SHA384, SHA512
* ***GO***: MD4,MD5,SHA1,SHA224,SHA256,SHA384,SHA512,RMD160,SHA3-224,SHA3-256,SHA3-384,SHA3-512,SHA512-224,SHA512-256,BLAKE2s256,BLAKE2b256,BLAKE2b384,BLAKE2b512
* ***OPENSSL***: MD4,MD5,SHA1,SHA224,SHA256,SHA384,SHA512,RMD160,SHA3-224,SHA3-256,SHA3-384,SHA3-512,SHA512-224,SHA512-256,BLAKE2s256,BLAKE2b512,SHAKE128,SHAKE256,SM3

icon credit: https://www.flaticon.com/authors/freepik
