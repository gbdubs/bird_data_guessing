package main

import (
	"errors"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
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
				Name:    "latin_name",
				Aliases: []string{"l"},
				Usage:   "the latin name of the bird to look up",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Whether or not to include debugging information, like where scores came from and which regexes they matched.",
			},
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "Whether to print the output or silently succeed, if the command succeeds.",
			},
		},
		Action: func(c *cli.Context) error {
			ln := c.String("latin_name")
			if ln == "" {
				return errors.New("latin_name must be provided")
			}
			v := c.Bool("verbose")
			d := c.Bool("debug")
			input := &bird_data_guessing.Input{
				LatinName: ln,
				Debug:     d,
			}
			output, err := input.Execute()
			if err != nil {
				return err
			}
			if v {
				spew.Dump(output.Data)
				spew.Dump(output.Attributions)
				if d {
					spew.Dump(output.DebugDatas)
				}
			}
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
