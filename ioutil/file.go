package ioutil

import (
	"io/ioutil"
	"os"

	"github.com/JREAMLU/j-kit/constant"
)

// ReadAll read all
func ReadAll(path string) (string, error) {
	fi, err := os.Open(path)
	if err != nil {
		return constant.EmptyStr, err
	}

	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	if err != nil {
		return constant.EmptyStr, err
	}

	return string(fd), nil
}

// ReadAllBytes read all bytes
func ReadAllBytes(path string) ([]byte, error) {
	fi, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fi.Close()
	return ioutil.ReadAll(fi)
}

// WriteFile write file
func WriteFile(path, content string, isOverride bool) error {
	var flag int

	if isOverride {
		flag = os.O_CREATE | os.O_RDWR | os.O_TRUNC
	} else {
		flag = os.O_CREATE | os.O_EXCL | os.O_RDWR
	}

	f, err := os.OpenFile(path, flag, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(content)
	return nil
}

// MkdireAll mkdir all
func MkdireAll(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// create controller directory
		if err := os.MkdirAll(path, 0777); err != nil {
			return err
		}
	}

	return nil
}

// CheckFileExists check file exists
func CheckFileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

// CreateFileOfTrunc create file of trunc
func CreateFileOfTrunc(filePath string) (*os.File, error) {
	if CheckFileExists(filePath) {
		return os.OpenFile(filePath, os.O_TRUNC, 0666)
	}

	return os.Create(filePath)
}

// CreateFileOfAppend create file of append
func CreateFileOfAppend(filePath string) (*os.File, error) {
	return os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
}
