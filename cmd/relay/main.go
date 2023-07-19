package main

import (
	"flag"
	"log"
	"os"
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

	relay, err := newRelay(
		cfg.Port,
		cfg.DatabaseURL,
		cfg.Nip11Pubkey,
		cfg.Nip11Contact,
		cfg.Nip11Description,
		cfg.Nip11Version,
	)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if err := relay.Start(); err != nil {
		log.Printf("relay err: %v\n", err)
		os.Exit(1)
	}
}
