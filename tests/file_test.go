package file_test

import (
	"testing"

	"github.com/getsumio/getsum/internal/file"
	"github.com/getsumio/getsum/internal/status"
)

func TestNotExist(t *testing.T) {

	f := &file.File{Url: "./someFile", Status: &status.Status{}}
	err := f.Fetch(60)
	if err == nil {
		panic("Expected an error while file doesnt exist")
	}

}

func TestExist(t *testing.T) {

	f := &file.File{Url: "./emptyFile", Status: &status.Status{}}
	err := f.Fetch(60)
	if err != nil {
		panic("UnExpected error while file exist " + err.Error())
	}

}

func TestSymlink(t *testing.T) {

	f := &file.File{Url: "./symLink", Status: &status.Status{}}
	err := f.Fetch(60)
	if err == nil {
		panic("For symlinks there should be an error")
	}
}

func TestDir(t *testing.T) {
	f := &file.File{Url: "./testDir", Status: &status.Status{}}
	err := f.Fetch(60)
	if err == nil {
		panic("For dirs there should be an error")
	}

}

func TestContentLength(t *testing.T) {
	f := &file.File{Url: "https://www.github.com", Status: &status.Status{}}
	err := f.Fetch(60)
	if err == nil {
		panic("If content length not provided there should be an error")
	}

}

func TestEmpty(t *testing.T) {
	f := &file.File{Url: "", Status: &status.Status{}}
	err := f.Fetch(60)
	if err == nil {
		panic("Empty file not allowed")
	}

}
