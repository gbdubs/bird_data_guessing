package main

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/gbdubs/bird_data_guessing"
)

func main() {
	input := bird_data_guessing.Input{
		LatinName:   "Branta hutchinsii",
		EnglishName: "Western tanager",
	}
	output, err := input.Execute()
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(output)
}
