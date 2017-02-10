package main

import (
	"flag"
	"log"

	orcareaper "github.com/robzienert/orca-reaper"
)

func main() {
	c, err := orcareaper.ParseConfig()
	if err != nil {
		flag.Usage()
		log.Fatalf("could not start reaper: %s", err.Error())
	}

	if err := orcareaper.Run(c); err != nil {
		log.Fatalf("encountered fatal error: %s", err.Error())
	}
}
