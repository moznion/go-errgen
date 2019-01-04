package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/moznion/go-errgen"
)

// auto-fill with `make build`
var (
	version  string
	revision string
)

var typ string
var prefix string
var outputFilePath string
var versionOption bool

func init() {
	flag.StringVar(&typ, "type", "", "[mandatory] struct type name of source of error definition")
	flag.StringVar(&prefix, "prefix", "ERR-", "[optional] prefix of error type")
	flag.StringVar(&outputFilePath, "out-file", "", "[optional] the output destination path of the generated code")
	flag.BoolVar(&versionOption, "version", false, "show version and revision")

	flag.Parse()

	if versionOption {
		fmt.Printf("v%s/%s\n", version, revision)
		os.Exit(0)
	}

	if typ == "" {
		log.Fatal("[ERROR] mandatory parameter `-type` is missing")
	}
}

func main() {
	errgen.Run(typ, prefix, outputFilePath)
}
