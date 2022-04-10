package util

import "os"

// Exists check file existence
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
