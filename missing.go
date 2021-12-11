package bird_data_guessing

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func isKnownMissing(key string) bool {
	_, err := os.Stat(missingPath(key))
	if err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	}
	panic(err)
}

func markMissing(key string) {
	path := missingPath(key)
	err := os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(path, []byte("404"), 0777)
	if err != nil {
		panic(err)
	}
}

func missingError(key string) error {
	return errors.New(fmt.Sprintf("404: Permanently Missing %s", key))
}

func missingPath(key string) string {
	return fmt.Sprintf("/tmp/bird_data_guessing/missing/%s.txt", key)
}
