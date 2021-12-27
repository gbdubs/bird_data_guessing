package bird_data_guessing

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gbdubs/bird"
)

var missingFilePattern = "/memo/missing/%s.txt"
var missingLock = sync.RWMutex{}
var missingContents = make(map[string]string)

func isMissing(site string, name bird.BirdName) bool {
	missingLock.RLock()
	if _, ok := missingContents[site]; !ok {
		missingLock.RUnlock()
		missingLock.Lock()
		doRead(site)
		missingLock.Unlock()
		missingLock.RLock()
	}
	result := strings.Contains(missingContents[site], birdRow(name))
	missingLock.RUnlock()
	return result
}

func recordMissing(site string, name bird.BirdName) {
	if isMissing(site, name) {
		return
	}
	missingLock.Lock()
	doWrite(site, name)
	missingLock.Unlock()
}

func birdRow(name bird.BirdName) string {
	return fmt.Sprintf("%s - %s", name.English, name.Latin)
}

func doWrite(site string, name bird.BirdName) {
	fp := fmt.Sprintf(missingFilePattern, site)
	s := fmt.Sprintf("%s\n%s", missingContents[site], birdRow(name))
	missingContents[site] = s
	err := ioutil.WriteFile(fp, []byte(s), 0777)
	if err != nil {
		panic(err)
	}
}

func doRead(site string) {
	fp := fmt.Sprintf(missingFilePattern, site)
	b, err := ioutil.ReadFile(fp)
	if errors.Is(err, os.ErrNotExist) {
		b = []byte("")
		err = os.MkdirAll(filepath.Dir(fp), 0777)
	}
	if err != nil {
		panic(err)
	}
	missingContents[site] = string(b)
}
