package pkg

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/leveldbcache"
	"github.com/syndtr/goleveldb/leveldb"
)

var Cache httpcache.Cache

func init() {
	dir := filepath.Join("tmp", "w360")
	_ = os.Mkdir(dir, os.ModePerm)

	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		log.Fatal(err)
	}

	Cache = leveldbcache.NewWithDB(db)
}
