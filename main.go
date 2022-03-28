package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/teguhatma/simple-bank/api"
	db "github.com/teguhatma/simple-bank/db/sqlc"
	"github.com/teguhatma/simple-bank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config file", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal("cannot start server", err)
	}
}
