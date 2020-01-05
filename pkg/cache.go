package pkg

import (
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/leveldbcache"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"os"
)

var Cache httpcache.Cache

func init() {
	dir := "cache.w360"
	os.Mkdir(dir, os.ModePerm)

	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		log.Fatal(err)
	}

	Cache = leveldbcache.NewWithDB(db)
}
