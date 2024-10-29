package storage

import (
	"io/ioutil"
	"os"
	"strings"
)

type StorageFileIndex struct {
	Prefix  string
	Dir     string
	Current string
	Cursor  string
}

const LEDGER_FILE_SIZE_LIMIT = 268435456

var DATA_ROOT_DIR = os.Getenv("QUANTUM_DATA_DIR")

func (si StorageFileIndex) ListFiles() []string {
	files, _ := ioutil.ReadDir(si.Dir)

	var matchedFiles []string

	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), si.Prefix) {
			matchedFiles = append(matchedFiles, file.Name())
		}
	}

	return matchedFiles
}
