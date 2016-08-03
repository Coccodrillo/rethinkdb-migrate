# rethinkdb-migrate

## Database migration manager for Rethinkdb, written in Go and using Dan's [Gorethink](https://github.com/dancannon/gorethink) library

### Prerequisites:

 - Glide, database migration tool.

    Install:
      - On Mac OS X you can install the latest release via [Homebrew](https://github.com/Homebrew/homebrew):

      ```
        $ brew install glide
      ```

      - On 12.04 and newer install from the PPA:

      ```
      sudo add-apt-repository ppa:masterminds/glide && sudo apt-get update
      sudo apt-get install glide
      ```

### Instalation

 ```
 $ go get github.com/Coccodrillo/rethinkdb-migrate
 ```

 or as a library

 ```
 $ go get github.com/Coccodrillo/rethinkdb-migrate/base
 ```

 where you need to pass Gorethink Session to

```
session, err = r.Connect(r.ConnectOpts{
	Address:   c.Address,
	Database:  c.Database,
	Username:  c.Username,
	Password:  c.Password,
	TLSConfig: t,
})
if err != nil {
	log.Fatalf("error: %v", err)
}
b := base.NewBaseMigration(session, "migrations")

```
### Configuration
Into a file config.yaml, pass connection details

### Migrations
At the moment, the queries are hardcoded to subpackage migrations. Write a method with a receiver for Migration structs. It has to be exported and ending in underscore and number (i.e. MIgration_1) which will count as migration. It has to return a gorethink.Term structure that will then be queried. It attempts to create migrations table and write the applied migrations into it. An example migration:

 ```
 func (m *Migration) Is_the_table_settings_ready_1(up bool) (term r.Term) {
	if up {
		term = r.TableList().Do(func(result r.Term) r.Term {
			return r.Branch(result.Contains("settings"),
				nil,
				r.TableCreate("settings"),
			)
		})
	} else {
		term = r.TableList().Do(func(result r.Term) r.Term {
			return r.Branch(result.Contains("settings"),
				r.TableDrop("settings"),
				nil,
			)
		})
	}
	return term
}
 ```

### Parameters

 -config=config.yml   Config file with connection (missing implementation)

 -env="development"   Env (missing implementation)

 -limit=1             Limit migrations to run

 -strict=true         Abort migrations on first error

 -check               Just list migrations to be applied

### Commands

Up - Migrates the database to the most recent version available [defaults to limit-0, runs all migrations]

Down - Migrates the database down to undo changes [defaults to limit-1, reverts one migration]

### Example Usage

```
 $ rethinkdb-migrate up -limit 1
 ```

### Todo
- [ ] Implementing missing configuration options
- [ ] More configuration options
- [x] Usage as a library
= [ ] Lexer for native queries
- [ ] Status and Info commands
- [ ] Better logging
- [ ] Adding Godoc

This was done as an "scratch my itch" project, done within an afternoon, I am sure it could be improved a lot, but I couldn't find anything, so this looks fine as a starting point. Any contributions in the above-mentioned areas (and a lot of other or that matter) are welcome and appreciated.
