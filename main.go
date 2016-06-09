package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
)

var ui cli.Ui

func main() {
	os.Exit(func() int {
		b := NewBaseMigration()
		b.SetUp()
		ui = &cli.BasicUi{Writer: os.Stdout}

		cli := &cli.CLI{
			Args: os.Args[1:],
			Commands: map[string]cli.CommandFactory{
				"up": func() (cli.Command, error) {
					return &Up{b}, nil
				},
				"down": func() (cli.Command, error) {
					return &Down{b}, nil
				},
			},
			HelpFunc: cli.BasicHelpFunc("rethinkdb-migrate"),
			Version:  "0.0.1",
		}

		exitCode, err := cli.Run()
		if err != nil {
			log.Printf("Error while executing - %s\n", err.Error())
			return 1
		}
		return exitCode
	}())
}
