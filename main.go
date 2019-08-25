package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/spy16/bifrost/server"
)

func main() {
	addr := flag.String("addr", ":8083", "Address to start the client app")
	cfgFile := flag.String("config", "config.toml", "Configuration file path")
	flag.Parse()

	cfg, err := readConfigs(*cfgFile)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	srv := server.New(*cfg, "./web/templates", "./web/static")
	log.Printf("starting server on %s...", *addr)
	log.Fatalf("server exited: %v", http.ListenAndServe(*addr, srv))
}

func readConfigs(cfgFile string) (*server.Config, error) {
	fh, err := os.Open(cfgFile)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	cfg := server.Config{}
	if _, err := toml.DecodeReader(fh, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
