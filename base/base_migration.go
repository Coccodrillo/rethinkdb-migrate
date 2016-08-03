package base

import (
	"github.com/coccodrillo/rethinkdb-migrate/migrations"
	r "gopkg.in/dancannon/gorethink.v2"

	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type BaseMigration struct {
	session     *r.Session
	Limit       int
	Check       bool
	Strict      bool
	packageName string
}

func NewBaseMigration(session *r.Session, packageName string) *BaseMigration {
	b := &BaseMigration{Strict: true}
	b.session = session
	b.packageName = packageName
	if err := b.setUp(); err != nil {
		log.Fatalf("error: %v", err)
	}
	return b
}

func (b *BaseMigration) setUp() (err error) {
	_, err = r.TableList().Do(func(result r.Term) r.Term {
		return r.Branch(result.Contains(b.packageName),
			nil,
			r.TableCreate(b.packageName),
		)
	}).Run(b.session)
	return err
}

func (b *BaseMigration) WriteMigration(migrationId int, migrationName string) (err error) {
	_, err = r.Table(b.packageName).Insert(map[string]interface{}{
		"id":         migrationId,
		"name":       migrationName,
		"created_at": r.Now(),
	}).Run(b.session)
	return err
}

func (b *BaseMigration) RemoveMigration(migrationId int) (err error) {
	_, err = r.Table(b.packageName).Get(migrationId).Delete().Run(b.session)
	return err
}

func (b *BaseMigration) GetLastMigration() (lastId int) {
	res, err := r.Table(b.packageName).OrderBy(r.Desc("id")).Limit(1).Run(b.session)
	if err == nil {
		var row map[string]interface{}
		err := res.One(&row)
		if val, ok := row["id"]; ok && err == nil {
			lastId = int(val.(float64))
		}
	}
	return lastId
}

func (b *BaseMigration) Runner(up bool) int {
	migrationList := reflect.TypeOf(&migrations.Migration{})
	lastMigrationId := b.GetLastMigration()
	m := make(map[int]*reflect.Method)
	numMethod := migrationList.NumMethod()
	var keys []int
	for i := 0; i < numMethod; i++ {
		method := migrationList.Method(i)
		toInt, err := strconv.Atoi(strings.Split(method.Name, "_")[len(strings.Split(method.Name, "_"))-1])
		if err == nil && ((up && toInt > lastMigrationId) || (!up && toInt <= lastMigrationId)) {
			m[toInt] = &method
			keys = append(keys, toInt)
		}
	}
	if up {
		sort.Ints(keys)
	} else {
		sort.Sort(sort.Reverse(sort.IntSlice(keys)))
	}
	if b.Limit > len(keys) || b.Limit == 0 {
		b.Limit = len(keys)
	}
	direction := "up"
	if !up {
		direction = "down"
	}
	if b.Limit > 0 {
		log.Printf("Migrating %d migrations %s", b.Limit, direction)
	} else {
		log.Println("All migrations up to date")
	}
	for i := 0; i < b.Limit; i++ {
		method := m[keys[i]]
		log.Printf("Migrating %s", method.Name)
		m := &migrations.Migration{}
		term := method.Func.Interface().(func(*migrations.Migration, bool) r.Term)(m, up)
		if term.String() != "" {
			if b.Check {
				log.Printf("Query to execute: %s", term.String())
			} else {
				res, err := term.Run(b.session)
				defer res.Close()
				if err != nil {
					log.Printf("Error while appying migration: %v", err)
					if b.Strict {
						log.Println("Aborting migrations")
						return 1
					}
				} else {
					log.Printf("Migration applied succesfully")
					if up {
						b.WriteMigration(keys[i], method.Name)
					} else {
						log.Printf("remove %d", keys[i])
						b.RemoveMigration(keys[i])
					}
				}
			}
		} else {
			log.Println("Empty migration, skipping...")
		}
	}
	return 0
}
