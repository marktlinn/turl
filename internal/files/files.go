package files

import (
	"log"
	"os"
)
func ReadFile(path string) *[]byte {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("couldn't read file %s - %v", path, err)
	}

	return &file
}
