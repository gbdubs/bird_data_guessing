package main

import (
	"errors"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/gbdubs/bird"
	"github.com/gbdubs/bird_data_guessing"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "Bird Data Guessing",
		Usage:   "A CLI for scraping the web to come up with a set of information likely to be true about a given bird.",
		Version: "1.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "latin_name",
				Usage: "the latin name of the bird to look up",
			},
			&cli.StringFlag{
				Name:  "english_name",
				Usage: "the english name of the bird to look up",
			},
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "Whether to print the output or silently succeed, if the command succeeds.",
			},
		},
		Action: func(c *cli.Context) error {
			ln := c.String("latin_name")
			en := c.String("english_name")
			if ln == "" {
				return errors.New("latin_name must be provided")
			}
			if en == "" {
				return errors.New("english_name must be provided")
			}
			v := c.Bool("verbose")
			input := &bird_data_guessing.Input{
				Names: []bird.BirdName{
					bird.BirdName{
						Latin:   ln,
						English: en,
					},
				},
			}
			output, err := input.Execute()
			if err != nil {
				if v {
					spew.Dump(err)
				}
				return err
			}
			if v {
				spew.Dump(output.BirdData)
			}
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
