package utils

import "os"

func HasDirectory(dirName string) (bool, error) {
	_, err := os.Stat(dirName)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
