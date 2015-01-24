// Collection data structure for database
package db

import (
	"os"
	"path"
)

const (
	INDEX_FILE = "_id.index"
)

type Collection struct {
	db   *Database
	name string
}

func OpenCollection(db *Database, name string) (*Collection, error) {
	collection := &Collection{db: db, name: name}
	return collection, collection.bootstrap()
}

func (col *Collection) bootstrap() error {
	if err := os.MkdirAll(path.Join(col.db.path, col.name), 0700); err != nil {
		return err
	}
	return nil
}

func (col *Collection) close() error {
	return nil
}
