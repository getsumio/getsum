package file_test

import (
	"testing"

	"github.com/getsumio/getsum/internal/file"
	"github.com/getsumio/getsum/internal/status"
)

func TestNotExist(t *testing.T) {

	f := &file.File{Url: "./someFile", Status: &status.Status{}}
	err := f.Fetch(60, false)
	if err == nil {
		t.Errorf("Expected an error while file doesnt exist")
	}

	t.Logf("An expected error received: %s", err.Error())
}

func TestExist(t *testing.T) {

	f := &file.File{Url: "./emptyFile", Status: &status.Status{}}
	err := f.Fetch(60, false)
	if err != nil {
		t.Errorf("UnExpected error while file exist " + err.Error())
	}

}

func TestSymlink(t *testing.T) {

	f := &file.File{Url: "./symLink", Status: &status.Status{}}
	err := f.Fetch(60, false)
	if err == nil {
		t.Errorf("For symlinks there should be an error")
	}
	t.Logf("An expected error received: %s", err.Error())
}

func TestDir(t *testing.T) {
	f := &file.File{Url: "./testDir", Status: &status.Status{}}
	err := f.Fetch(60, false)
	if err == nil {
		t.Errorf("For dirs there should be an error")
	}

	t.Logf("An expected error received: %s", err.Error())
}

func TestContentLength(t *testing.T) {
	f := &file.File{Url: "https://www.github.com", Status: &status.Status{}}
	err := f.Fetch(60, false)
	if err == nil {
		t.Errorf("If content length not provided there should be an error")
	}
	t.Logf("An expected error received: %s", err.Error())

}

func TestEmpty(t *testing.T) {
	f := &file.File{Url: "", Status: &status.Status{}}
	err := f.Fetch(60, false)
	if err == nil {
		t.Errorf("Empty file not allowed")
	}

	t.Logf("An expected error received: %s", err.Error())
}
