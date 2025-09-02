package pathtools

import (
	"os"
	"syscall"
)

// DirExists simply returns true if the pathName is a existing directory
func DirExists(pathName string) bool {
	stat, err := os.Stat(pathName)
	return err == nil && stat.IsDir()
}

// CreatePath Create full path, with permissions updated from parent folder.
func CreatePath(path string) error {
	if DirExists(path) {
		return nil
	}
	return createPath(path)
}

func createPath(path string) error {
	oldmask := syscall.Umask(0)
	defer syscall.Umask(oldmask)
	return os.MkdirAll(path, os.ModeSticky|os.ModePerm)
}
