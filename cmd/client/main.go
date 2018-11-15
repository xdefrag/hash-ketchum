package main

import (
	"flag"
	"log"
	"os"

	"github.com/xdefrag/hash-ketchum/api"
	"github.com/xdefrag/hash-ketchum/pkg/tools/config"
)

var reqs = flag.Int("n", 10000, "Number of request")

func main() {
	flag.Parse()

	cfg := api.ClientConfig{
		Host: config.WithDefaultString(os.Getenv("SERVER_HOST"), "0.0.0.0"),
		Port: config.WithDefaultInt(config.Atoi(os.Getenv("SERVER_PORT")), 8080),
	}

	cred := api.Credentials{
		Login: config.WithDefaultString(os.Getenv("SERVER_LOGIN"), "user"),
	}

	logger := log.New(os.Stdout, "CLIENT: ", 0)

	client := api.NewClient(cfg, cred, logger)

	if err := client.Run(*reqs); err != nil {
		log.Fatal(err)
	}
}
