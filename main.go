package main

import (
	"flag"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

func main() {
	log.Printf("build info: commit: %v date: %v\n", commit, buildDate)

	configPath := flag.String("config", "config.yml", "location of config file")
	flag.Parse()

	var cfg Config
	log.Printf("loading config from file %q\n", *configPath)
	if err := cfg.Load(*configPath); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Printf("%#v\n", cfg)

	subscriptionsDB, err := sqlx.Connect("postgres", cfg.SubscriptionsDBURL)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	relay, err := newRelay(cfg, subscriptionsDB)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	relay.server.Router().HandleFunc("/admin", adminHandler(cfg, relay.storage))
	relay.server.Router().HandleFunc("/admin/delete", adminDeleteHandler(cfg, relay.storage))

	if err := relay.Start(); err != nil {
		log.Printf("relay err: %v\n", err)
		os.Exit(1)
	}
}
