package bird_data_guessing

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getDocumentFromUrl(url string, missingKey string) (*goquery.Document, error) {
	r := &goquery.Document{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return r, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return r, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		if resp.StatusCode == 404 {
			markMissing(missingKey)
			return r, missingError(missingKey)
		}
		return r, errors.New(fmt.Sprintf("Request failed: %d %s", resp.StatusCode, resp.Status))
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return doc, err
	}
	return doc, nil
}

func getDocumentFromString(html string) *goquery.Document {
	res, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		panic(err)
	}
	return res
}
