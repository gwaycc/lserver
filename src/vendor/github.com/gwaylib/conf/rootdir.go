package conf

import (
	"os"
	"path/filepath"
)

func RootDir() string {
	dir := os.Getenv("PRJ_ROOT")
	if len(dir) == 0 {
		// Compatible the old version
		dir = os.Getenv("PJ_ROOT")
	}

	if len(dir) == 0 {
		panic("Need PRJ_ROOT environment for project directory")
	}

	p, err := filepath.Abs(dir)
	if err != nil {
		panic(err)
	}
	return p
}
