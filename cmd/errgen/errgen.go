package main

import (
	"flag"
	"log"

	"github.com/moznion/go-errgen"
)

var typ string
var prefix string

// TODO out filename

func init() {
	flag.StringVar(&typ, "type", "", "[mandatory] struct type name of source of error definition")
	flag.StringVar(&prefix, "prefix", "ERR-", "[optional] prefix of error type")

	flag.Parse()

	if typ == "" {
		log.Fatal("[ERROR] mandatory parameter `-type` is missing")
	}
}

func main() {
	errgen.Run(typ, prefix)
}
