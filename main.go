package main

import (
	"fmt"
	"os"

	"github.com/phelrine/secretsenv/lib"
)

func main() {
	cli := &lib.CLI{
		Loaders: map[string]lib.SecretLoader{
			"aws": lib.NewAWSLoader(),
		},
	}
	err := cli.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing command line: %v\n", err)
		os.Exit(1)
	}
	err = cli.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}
	err = cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running: %v\n", err)
		os.Exit(1)
	}
}
