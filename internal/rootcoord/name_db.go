package rootcoord

import "github.com/milvus-io/milvus/internal/util/typeutil"

type nameDb struct {
	db2Name2ID map[string]map[string]typeutil.UniqueID // database -> collection name -> collection id
}

func (db *nameDb) insert(dbName string, collectionName string, collectionID UniqueID) {
	if _, ok := db.db2Name2ID[dbName]; !ok {
		db.db2Name2ID[dbName] = make(map[string]typeutil.UniqueID)
	}
	db.db2Name2ID[dbName][collectionName] = collectionID
}

func newNameDb() *nameDb {
	return &nameDb{
		db2Name2ID: make(map[string]map[string]typeutil.UniqueID),
	}
}
