package main

import (
	"log"

	"github.com/atharvakadlag/splitfree/cmd/api"
	"github.com/atharvakadlag/splitfree/config"
	"github.com/atharvakadlag/splitfree/db"
	"github.com/go-sql-driver/mysql"
)

func main() {
	log.Println(config.Envs.DBPasswd, config.Envs.DBUser)
	db, err := db.NewMySQLStorage(*&mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPasswd,
		Addr:                 config.Envs.DBAddr,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	server := api.CreateAPI(":8001", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
