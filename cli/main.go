package main

import (
	"fmt"
	"log"
	"os"
	"salami/common/constants"
	"salami/common/driver"
	"strings"

	"github.com/urfave/cli/v2"
)

type SalamiMultiError struct {
	errors []error
}

func (m *SalamiMultiError) Error() string {
	msgs := []string{}
	for _, err := range m.errors {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, ", ")
}

func (m *SalamiMultiError) Errors() []error {
	return m.errors
}

func main() {
	app := &cli.App{
		Name:        "salami",
		HelpName:    "Salami",
		Version:     constants.SalamiVersion,
		Usage:       "a declarative DSL for cloud infrastructure based on natural language descriptions",
		UsageText:   "salami [global options] [command] [command options]",
		HideVersion: true,
		Suggest:     true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Enable verbose mode",
			},
		},
		Commands: []*cli.Command{
			{
				Name:      "compile",
				Usage:     "Runs the compilation end-to-end",
				UsageText: "salami [global options] compile [command options]",
				Action: func(cCtx *cli.Context) error {
					verbose := cCtx.Bool("verbose")
					errors := driver.Run(verbose)

					if len(errors) == 1 {
						return errors[0]
					} else if len(errors) > 1 {
						return &SalamiMultiError{errors: errors}
					}

					return nil
				},
			},
			{
				Name:  "version",
				Usage: "Prints the version",
				Action: func(cCtx *cli.Context) error {
					fmt.Println("Salami version " + constants.SalamiVersion)
					return nil
				},
			},
		},
		Authors: []*cli.Author{{
			Name:  "Petr Gazarov",
			Email: "petrgazarov@gmail.com",
		}},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
