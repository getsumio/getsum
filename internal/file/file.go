package file

import (
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

type IFile interface {
	Path() string
	Data() ([]byte, error)
	Fetch(timeout int) error
	IsRemote() bool
}

type File struct {
	path   string
	data   []byte
	Url    string
	Status *status.Status
	Size   int64
	Proxy  string
}

func (f *File) Path() string {
	return f.path
}

func (f *File) IsRemote() bool {
	return strings.HasPrefix(f.Url, "http") || strings.HasPrefix(f.Url, "ftp")
}

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

var mux sync.Mutex
var fetchedSize int64 = -1

func fetchRemote(f *File, timeout int) error {

	mux.Lock()
	defer mux.Unlock()
	filename := path.Base(f.Url)
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

	header, err := client.Head(f.Url)
	if err != nil {
		return err
	}
	defer header.Body.Close()
	size, err := strconv.Atoi(header.Header.Get("Content-Length"))
	if err != nil {
		return err
	}
	f.Size = int64(size)
	f.path = filename

	quit := make(chan bool)
	defer close(quit)
	go downloadFile(quit, f)

	resp, err := client.Get(f.Url)
	if err != nil {
		quit <- true
		return err
	}
	defer resp.Body.Close()

	resp.Header.Set("Connection", "Keep-Alive")
	resp.Header.Set("Accept-Language", "en-US")
	resp.Header.Set("User-Agent", "Mozilla/5.0")

	_, err = io.Copy(out, resp.Body)
	quit <- true

	f.Status.Type = status.FETCHED
	fetchedSize = f.Size
	return err
}

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

func remove(fileOrDir string) error {
	err := os.Remove(fileOrDir)
	if err != nil {
		return err
	}
	return nil

}

func getHttpClient(f *File, timeout int) *http.Client {
	proxyUrl := http.ProxyFromEnvironment
	if f.Proxy != "" {
		proxy, _ := url.Parse(f.Proxy)
		proxyUrl = http.ProxyURL(proxy)
	}
	tr := &http.Transport{
		Proxy: proxyUrl,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(timeout) * time.Second,
	}
	return client
}
