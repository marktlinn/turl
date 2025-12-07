package files

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
)

// Used to search for Username/.config/turl by default i.e. = $HOME/.config/turl/
const DEFAULT_TURL_ROOT_DIR = ".config/turl/"

func GetRouteFiles() {
	turlRootDirPath := os.Getenv("TURL_ROOT_DIR")
	if turlRootDirPath == "" {
		usr, err := user.Current()
		if err != nil {
			log.Fatalf("couldn't get current user %v", err)
		}
		turlRootDirPath = filepath.Join(usr.HomeDir, DEFAULT_TURL_ROOT_DIR)
	}

	dir, err := os.ReadDir(turlRootDirPath)
	if err != nil {
		log.Fatalf("couldn't read dir %s - %v", turlRootDirPath, err)
	}

	for i, file := range dir {
		fileExt := filepath.Ext(file.Name())
		if fileExt == ".yaml" || fileExt == ".yml" {
			log.Printf("OK => item(%d) - Directory(%s)\n", i, file)
		}
	}
}

func ReadFile(path string) *[]byte {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("couldn't read file %s - %v", path, err)
	}

	return &file
}
