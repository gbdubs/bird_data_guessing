package bird_data_guessing

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gbdubs/bird"
)

func readMemoized(name bird.BirdName) (bool, BirdData, error) {
	mp := memoPath(name)
	asBytes, err := ioutil.ReadFile(mp)
	if errors.Is(err, os.ErrNotExist) {
		return false, BirdData{}, nil
	} else if err == io.EOF {
		return true, BirdData{}, nil
	} else if err != nil {
		return false, BirdData{}, fmt.Errorf("While looking up memoized file at %s: %v", mp, err)
	}
	r := &BirdData{}
	err = xml.Unmarshal(asBytes, r)
	if err != nil {
		return false, BirdData{}, fmt.Errorf("While unmarshalling data from %s: %v", mp, err)
	}
	return true, *r, nil
}

func writeMemoized(bd BirdData) error {
	asBytes, err := xml.MarshalIndent(&bd, "", " ")
	if err != nil {
		return fmt.Errorf("While marshalling bird data for %s: %v", bd.Name.English, err)
	}
	mp := memoPath(bd.Name)
	err = os.MkdirAll(filepath.Dir(mp), 0777)
	if err != nil {
		return fmt.Errorf("While creating parent directory %s: %v", filepath.Dir(mp), err)
	}
	err = ioutil.WriteFile(mp, asBytes, 0777)
	if err != nil {
		return fmt.Errorf("While writing memoized file to %s: %v", mp, err)
	}
	return nil
}

func memoPath(name bird.BirdName) string {
	return fmt.Sprintf("/memo/bird_data_guessing/%s.xml", name.English)
}
