package file

import (
	"errors"
	"fmt"
	"os"
)

func validateLocal(f *File) error {
	if f == nil {
		return errors.New("File not initialized!")
	}
	err := validateUrl(f)
	if err != nil {
		return err
	}
	if f.IsRemote() {
		return errors.New(fmt.Sprintf("Given file has remote url %s", f.Url))
	}
	info, err := os.Stat(f.Url)
	if os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("Given url %s can not be accessed or not understood, supported protocols file, http, https, ftp", f.Url))
	} else if err != nil {
		return err
	}
	if info.IsDir() {
		return errors.New(fmt.Sprintf("Given file %s is directory!", f.Url))
	}

	fInfo, err := os.Lstat(f.Url)
	if err != nil {
		return err
	}
	if fInfo.Mode()&os.ModeSymlink != 0 {
		return errors.New(fmt.Sprintf("Given file %s is a symbolik link!", f.Url))
	}

	return nil
}

func validateRemote(f *File) error {
	return nil

}

func validateUrl(f *File) error {
	if f.Url == "" {
		return errors.New("Empty file url provided!")
	}
	return nil

}
