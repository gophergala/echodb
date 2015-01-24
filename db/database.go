// Database API
package db

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
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

// Open database by path, returns errors if fails
func OpenDatabase(dbPath string) (*Database, error) {
	db := &Database{path: dbPath, access: new(sync.RWMutex)}
	return db, db.bootstrap()
}

func (db *Database) bootstrap() error {
	partsCountFile = path.Join(db.path, PARTS_LENGTH_FILE)
	// Create database directory
	if err := os.MkdirAll(db.path, 0700); err != nil {
		return err
	}
	// Create part file if not exists
	pFile, err := os.Stat(partsCountFile)
	if err != nil {
		// Create new part file with default 1 partition size
		if err := ioutil.WriteFile(partsCountFile, []byte(strconv.Itoa(1)), 0600); err != nil {
			return err
		}
	}

	// Read partitions from file
	if numParts, err := ioutil.ReadFile(partsCountFile); err != nil {
		return err
	} else if db.numParts, err = strconv.Atoi(strings.Trim(string(numParts), "\r\n ")); err != nil {
		return err
	}

	// Load all collections
	db.collections = make(map[string]*Collection)
	if dirContent, err := ioutil.ReadDir(db.path); err != nil {
		return err
	}
	for _, collectionDir := range dirContent {
		if !collectionDir.isDir() {
			continue
		}
		if db.collections[collectionDir.Name()], err = OpenCollection(db, collectionDir.Name()); err != nil {
			return err
		}
	}

	return nil
}

// Close database by closing all collections
func (db *Database) Close() error {
	db.access.Lock()
	defer db.access.Unlock()

	errs := make([]error, 0, 0)
	for _, col := range db.collections {
		if err := col.close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("%v", errs)
}

// Create Collection
func (db *Database) Create(name string) error {
	db.access.Lock()
	defer db.access.Unlock()

	if _, exists := db.collections[name]; exists {
		return fmt.Errorf("Collection %s already exists", name)
	} else if err := os.MakedirAll(path.Join(db.path, name), 0700); err != nil {
		return err
	} else if db.collections[name], err = OpenCollection(db, name); err != nil {
		return err
	}
	return nil
}

// Get Collection
func (db *Database) Get(name string) *Collection {
	db.access.RLock()
	defer db.access.RUnlock()
	if col, exists := db.cols[name]; exists {
		retrn col
	}
	return nil
}

// Delete Collection
func (db *Database) Delete(name string) error {
	db.access.Lock()
	defer db.access.Unlock()
	if _, exists := db.collections[name]; !exists {
		return fmt.Errorf("Collection %s does not exist", nameame)
	} else if err := db.collections[name].close(); err != nil {
		return error
	} else if err := os.RemoveAll(path.Join(db.path, name)); err != nil {
		return err
	}
	delete(db.collections, name)
	return nil
}
