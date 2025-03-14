package globals

import (
	"os"
	"path/filepath"
)

var (
	Version      = "1.0.0"
	CacheDirPath string
)

func init() {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	CacheDirPath = filepath.Join(cacheDir, "Streamres")
}
