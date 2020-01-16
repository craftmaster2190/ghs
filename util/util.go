package util

import (
	"fmt"
	"os"
)

func Check(e error) {
	if e != nil {
		fmt.Printf("Error: %s\n", e.Error())
		os.Exit(1)
	}
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
