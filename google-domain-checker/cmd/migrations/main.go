package main

import (
	"database/sql"
	"flag"
	_ "github.com/IlyushaZ/check-domain/google-domain-checker/migration"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

const dialect = "postgres"
const driverName = "postgres"

type config struct {
	DbString string `yaml:"db_string"`
	Dir      string `yaml:"migrations_dir"`
}

func main() {
	var dir string
	flag.StringVar(&dir, "confdir", "./config/app.yaml", "path to config file")
	flag.Parse()

	file, err := ioutil.ReadFile(dir)
	if err != nil {
		log.Fatalf("error opening config file %s", dir)
	}

	var conf config
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		log.Fatal("error parsing config file")
	}

	args := flag.Args()
	command := args[0]
	switch command {
	case "create":
		if err := goose.Run("create", nil, conf.Dir, args[1:]...); err != nil {
			log.Fatalf("migrate run: %v", err)
		}
		return
	case "fix":
		if err := goose.Run("fix", nil, conf.Dir); err != nil {
			log.Fatalf("migrate run: %v", err)
		}
		return
	}

	db, err := sql.Open(driverName, conf.DbString)
	if err != nil {
		log.Fatalf("%q: %v", conf.DbString, err)
	}
	defer db.Close()

	if err := goose.SetDialect(dialect); err != nil {
		log.Fatal(err)
	}

	if err := goose.Run(command, db, conf.Dir); err != nil {
		log.Fatalf("Goose run error: %v", err)
	}
}
