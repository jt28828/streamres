package globals

import (
	"os"
	"path/filepath"
)

var VERSION = "1.0.0"

var CacheDirPath string

func init() {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	CacheDirPath = filepath.Join(cacheDir, "Streamres")
}
