package rethinkdb

import (
	"log"
	"strings"

	r "github.com/dancannon/gorethink"
)

type Migrator interface {
	CreateDb(client *r.Session, dbName string) error
	DropDb(client *r.Session, dbName string) error
	CreateTable(tableName string) error
	DropTable(tableName string) error
	AddIndex(tableName string, fields []string, opts map[string]interface{}) error
}

type RethinkDbMigration struct {
	Migrator
}

func (*RethinkDbMigration) CreateDb(client *r.Session, dbName string) error {
	if _, err := r.DBCreate(dbName).Run(client); err != nil {
		//log.Fatalln(err.Error())
		log.Println(err.Error())
		return err
	}
	return nil
}

func (*RethinkDbMigration) DropDb(client *r.Session, dbName string) error {
	if _, err := r.DBDrop(dbName).Run(client); err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (ar *RethinkDbMigration) CreateTable(client *r.Session, dbName string, tableName string) error {
	if _, err := r.DB(dbName).TableCreate(tableName).RunWrite(client); err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (ar *RethinkDbMigration) DropTable(client *r.Session, dbName string, tableName string) error {
	if _, err := r.DB(dbName).TableDrop(tableName).RunWrite(client); err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (ar *RethinkDbMigration) AddIndex(client *r.Session, dbName string, tableName string, fields []string, opts map[string]interface{}) error {
	if len(fields) == 1 {
		if _, err := r.DB(dbName).Table(tableName).IndexCreate(fields[0], r.IndexCreateOpts{Multi: true}).Run(client); err != nil {
			log.Println(err.Error())
			return err
		}
	} else {
		indexName := strings.Join(fields, "_")
		if _, err := r.DB(dbName).Table(tableName).IndexCreateFunc(indexName, func(row r.Term) interface{} {
			fieldSlice := []r.Term{}
			for _, element := range fields {
				fieldSlice = append(fieldSlice, row.Field(element))
			}

			return []interface{}{fieldSlice}
		}).RunWrite(client); err != nil {
			log.Println(err.Error())
			return err
		}
	}

	return nil
}
