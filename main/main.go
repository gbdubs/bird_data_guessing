package main

import (
	"log"

	"github.com/gbdubs/bird_data_guessing"
	"github.com/gbdubs/bird_region_rosters"
	"github.com/gbdubs/verbose"
)

func main() {
	brri := bird_region_rosters.Input{
		RegionCodes: []string{"USco"},
		Verbose:     verbose.New(),
	}
	brro, err := brri.Execute()
	if err != nil {
		log.Fatal(err)
	}

	bdgi := bird_data_guessing.Input{
		Names:   brro.Entries,
		Verbose: verbose.New(),
	}
	_, err = bdgi.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
