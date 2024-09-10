package database

import (
	"log"
	"sync"

	"github.com/jmoiron/sqlx"
)

var store sync.Map

func GetSqlxClient(key string) *sqlx.DB {
	client, ok := store.Load(key)
	if !ok {
		log.Fatal("please open db first", key)
	}

	return client.(*sqlx.DB)
}
