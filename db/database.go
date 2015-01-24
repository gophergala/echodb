// Database API
package db

import (
	"fmt"
	"io"
	"os"
	"path"
	"sync"
	"time"
)

const (
	PARTS_LENGTH_FILE = "_count" // Holds total count of parittions
)

type Database struct {
	path        string
	numParts    int
	collections map[string]*Collection
	access      *sync.RWMutex
}

func OpenDatabase(dbPath string) (*Database, error) {
	db := &Database{path: dbPath, access: new(sync.RWMutex)}
	return db, db.checkLoadError()
}

func (db *Database) checkLoadError() error {
	return nil
}
