package main

import (
	"flag"
	"github.com/coccodrillo/rethinkdb-migrate/base"
	"strings"
)

type Up struct {
	*base.BaseMigration
}

func (u *Up) Synopsis() string {
	return "Migrates the database to the most recent version available"
}

func (u *Up) Help() string {
	return strings.TrimSpace(`
	Usage: rethink-migrate up [flags] ...
	  Runs all migrations
	Options:
	  -config=config.yml    Config file with connection
	  -env="development"        Set environment
	  -strict=true          Abort migrations on first error
	  -limit=0              Limit migrations to run
	  -check                Just list migrations to be applied
	`)
}

func (u *Up) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("up", flag.ContinueOnError)
	cmdFlags.Usage = func() {
		ui.Output(u.Help())
	}
	cmdFlags.BoolVar(&u.Strict, "strict", true, "Abort migrations on first error")
	cmdFlags.IntVar(&u.Limit, "limit", 0, "Limit migrations to run")
	cmdFlags.BoolVar(&u.Check, "check", false, "Dry run")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}
	return u.Runner(true)
}

type Down struct {
	*base.BaseMigration
}

func (d *Down) Synopsis() string {
	return "Migrates the database down to undo changes"
}

func (d *Down) Help() string {
	return strings.TrimSpace(`
	Usage: rethink-migrate down [flags] ...
	  Runs all migrations
	Options:
	  -config=config.yml    Config file with connection
	  -env="development"    Env
	  -limit=1              Limit migrations to run
	  -strict=true          Abort migrations on first error
	  -check                Just list migrations to be applied
	`)
}

func (d *Down) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("up", flag.ContinueOnError)
	cmdFlags.Usage = func() {
		ui.Output(d.Help())
	}
	cmdFlags.BoolVar(&d.Strict, "strict", true, "Abort migrations on first error")
	cmdFlags.IntVar(&d.Limit, "limit", 1, "Limit migrations to run")
	cmdFlags.BoolVar(&d.Check, "check", false, "Dry run")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}
	return d.Runner(false)
}
