// Collection data structure for database
package db

import (
	"os"
)

const (
	INDEX_FILE_SUFFIX = ".index"
	INDEX_SEP         = "!"
)

type Collection struct {
	db   *Database
	name string
}

func OpenCollection(db *Database, name string) (*Collection, error) {
	collection := &Collection{db: db, name: name}
	return collection, collection.checkLoadError()
}

func (col *Collection) checkLoadErrro() error {
	return nil
}
