package file

import (
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/getsumio/getsum/internal/status"
)

//file interface
type IFile interface {
	Path() string
	Data() ([]byte, error)
	Fetch(timeout int) error
	IsRemote() bool
	Delete()
	Reset()
}

//file struct
type File struct {
	path        string
	data        []byte
	Url         string
	Status      *status.Status
	Size        int64
	Proxy       string
	StoragePath string
	SkipVerify  bool
}

//file location on local host
//if file is remote file something like http/https
//path will be present after file is downloaded
//so you need to call Fetch method first
func (f *File) Path() string {
	return f.path
}

//simply removes file if its already fetched
func (f *File) Delete() {
	if f.path != "" {
		os.Remove(f.path)
	}
}

//checks file is not on local
//TODO add file:/// support
func (f *File) IsRemote() bool {
	return strings.HasPrefix(f.Url, "http")
}

//a bit of ugly but resets global variable
//check variable comment for details
func (f *File) Reset() {
	f.Url = ""
	f.Size = 0
	f.path = ""
	f.data = nil
	f.Proxy = ""
	f.StoragePath = ""
	f.Status = nil
	fetchedSize = -1
}

//read the file data in bytes
func (f *File) Data() ([]byte, error) {
	if f.path == "" {
		return nil, errors.New("File not fetched yet you need to first call Fetch()")
	}
	bytes, err := ioutil.ReadFile(f.path)
	if err != nil {
		return nil, err
	}
	f.Status.Type = status.ALLOCATED
	return bytes, nil

}

//validates file
//if remote fetches file
//sets path location forthe stored file
func (f *File) Fetch(timeout int) error {
	err := validateUrl(f)
	if err != nil {
		return err
	}
	isRemote := f.IsRemote()
	if !isRemote {
		return fetchLocal(f)
	} else {
		return fetchRemote(f, timeout)
	}

	return nil
}

//TODO: sync.ONCE maybe better?
var mux sync.Mutex

//in case of user runs multiple checksum calculations
//on same machine there is no point to
//download file per routine
//so first routine takes the leads and downloads file
//then other routines check this variable
//if present they use existing path
var fetchedSize int64 = -1

//fetches file remotely
//unless not timedout sets details and path
func fetchRemote(f *File, timeout int) error {

	mux.Lock()
	defer mux.Unlock()
	filename := path.Base(f.Url)
	if f.StoragePath != "" {
		filename = strings.Join([]string{f.StoragePath, filename}, string(os.PathSeparator))
	}
	if fetchedSize > 0 { //another process already fetched
		f.Size = fetchedSize
		f.path = filename
		f.Status.Type = status.FETCHED
		return nil
	}
	err := validateRemote(f)
	if err != nil {
		return err
	}

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	client := getHttpClient(f, timeout)

	f.path = filename

	quit := make(chan bool)
	defer close(quit)

	resp, err := client.Get(f.Url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	resp.Header.Set("Connection", "Keep-Alive")
	resp.Header.Set("Accept-Language", "en-US")
	resp.Header.Set("User-Agent", "Mozilla/5.0")
	contentLength := resp.Header.Get("Content-Length")
	if contentLength == "" {
		return errors.New("Can not get content length, is this a binary file?")
	}
	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return errors.New("Can not parse content-length, is this binary? " + err.Error())
	}

	f.Size = int64(size)

	go downloadFile(quit, f)

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		quit <- true
		return err
	}

	quit <- true

	f.Status.Type = status.FETCHED
	fetchedSize = f.Size
	return err
}

//only validates file since its already hosted
//see validation for details
func fetchLocal(f *File) error {

	err := validateLocal(f)
	if err != nil {
		f.Status.Type = status.ERROR
		return err
	}
	f.path = f.Url
	f.Status.Type = status.FETCHED
	return nil

}

//removes file
func remove(fileOrDir string) error {
	err := os.Remove(fileOrDir)
	if err != nil {
		return err
	}
	return nil

}

//return http client wrapped by
//proxy settings
//as well as timeout value
func getHttpClient(f *File, timeout int) *http.Client {
	proxyUrl := http.ProxyFromEnvironment
	if f.Proxy != "" {
		proxy, _ := url.Parse(f.Proxy)
		proxyUrl = http.ProxyURL(proxy)
	}
	tr := &http.Transport{
		Proxy:           proxyUrl,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: f.SkipVerify},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(timeout) * time.Second,
	}
	return client
}
