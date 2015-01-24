// Collection data structure for database
package db

import (
	"../dbcore"
	"encoding/json"
	"math/rand"
	"os"
	"path"
	"strconv"
)

const (
	INDEX_FILE = "_idx"
)

type Collection struct {
	db    *Database
	name  string
	parts []*dbcore.Partition
}

func OpenCollection(db *Database, name string) (*Collection, error) {
	collection := &Collection{db: db, name: name}
	return collection, collection.bootstrap()
}

func (col *Collection) bootstrap() error {
	if err := os.MkdirAll(path.Join(col.db.path, col.name), 0700); err != nil {
		return err
	}
	col.parts = make([]*dbcore.Partition, col.db.numParts)

	for i := 0; i < col.db.numParts; i++ {
		var err error
		if col.parts[i], err = dbcore.OpenPartition(
			path.Join(col.db.path, col.name, col.name+"."+strconv.Itoa(i)),
			path.Join(col.db.path, col.name, INDEX_FILE+"."+strconv.Itoa(i))); err != nil {
			return err
		}
	}
	return nil
}

func (col *Collection) close() error {
	for i := 0; i < col.db.numParts; i++ {
		col.parts[i].Lock.Lock()
		col.parts[i].Close()
		col.parts[i].Lock.Unlock()
	}
	return nil
}

func (col *Collection) Count() int {
	col.db.access.RLock()
	defer col.db.access.RUnlock()

	count := 0
	for _, part := range col.parts {
		part.Lock.RLock()
		count += part.ApproxDocCount()
		part.Lock.RUnlock()
	}
	return count
}

// Insert a document into the collection.
func (col *Collection) Insert(doc map[string]interface{}) (id int, err error) {
	docJS, err := json.Marshal(doc)
	if err != nil {
		return
	}
	id = rand.Int()
	partNum := id % col.db.numParts
	col.db.access.RLock()
	part := col.parts[partNum]
	// Put document data into collection
	part.Lock.Lock()
	if _, err = part.Insert(id, []byte(docJS)); err != nil {
		part.Lock.Unlock()
		col.db.access.RUnlock()
		return
	}
	// If another thread is updating the document in the meanwhile, let it take over index maintenance
	if err = part.LockUpdate(id); err != nil {
		part.Lock.Unlock()
		col.db.access.RUnlock()
		return id, nil
	}
	part.UnlockUpdate(id)
	part.Lock.Unlock()
	col.db.access.RUnlock()
	return
}
