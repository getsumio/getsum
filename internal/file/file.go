package file

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type IFile interface {
	Path() string
	Data() ([]byte, error)
	Fetch(timeout int) error
	IsRemote() bool
}

type File struct {
	path        string
	data        []byte
	Url         string
	Status      string
	StatusValue string
	Size        int64
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
	f.Status = "ALLOCATED"
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

func fetchRemote(f *File, timeout int) error {

	filename := path.Base(f.Url)
	_, err := os.Stat(filename)
	if os.IsExist(err) {
		return fetchLocal(f)
	}
	err = validateRemote(f)
	if err != nil {
		return err
	}

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	header, err := http.Head(f.Url)
	if err != nil {
		return err
	}
	defer header.Body.Close()
	size, err := strconv.Atoi(header.Header.Get("Content-Length"))
	if err != nil {
		f.Status = "ERROR"
		return err
	}
	f.Size = int64(size)
	f.path = filename

	quit := make(chan bool)
	go downloadFile(quit, f)

	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	resp, err := client.Get(f.Url)
	if err != nil {
		quit <- true
		f.Status = "ERROR"
		return err
	}
	defer resp.Body.Close()

	resp.Header.Set("Connection", "Keep-Alive")
	resp.Header.Set("Accept-Language", "en-US")
	resp.Header.Set("User-Agent", "Mozilla/5.0")

	_, err = io.Copy(out, resp.Body)
	quit <- true

	f.Status = "FETCHED"
	return err
}

func fetchLocal(f *File) error {

	err := validateLocal(f)
	if err != nil {
		f.Status = "ERROR"
		return err
	}
	f.path = f.Url
	f.Status = "FETCHED"
	return nil

}

func remove(fileOrDir string) error {
	err := os.Remove(fileOrDir)
	if err != nil {
		return err
	}
	return nil

}
