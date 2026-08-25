package main

import (
	"flag"
	"log"

	"buscharge/internal/console"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	dataDir := flag.String("data", "./buscharge-data", "data directory")
	webDir := flag.String("web", "./web", "web asset directory")
	seed := flag.Bool("seed", false, "seed demo data before serving")
	flag.Parse()

	server, err := console.NewServer(console.Config{
		Addr:    *addr,
		DataDir: *dataDir,
		WebDir:  *webDir,
	})
	if err != nil {
		log.Fatal(err)
	}
	if *seed {
		if err := server.SeedDemo(); err != nil {
			log.Fatal(err)
		}
	}
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
