// Command errorcodes writes the docs site's error table from the catalogue.
//
// The page at /docs/backend/errors is the documentation a frontend developer
// reads instead of triggering each error to find out what it is. Generating it
// from internal/errorcodes means it cannot fall behind the API: a test compares
// the committed file with what this would write and fails if they differ.
//
//	go run ./tools/errorcodes              # rewrite docs/components/error-codes.ts
//	go run ./tools/errorcodes -check       # fail if it is out of date
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/MUKE-coder/grit/v3/internal/errorcodes"
)

func main() {
	out := flag.String("out", "docs/components/error-codes.ts", "file to write")
	check := flag.Bool("check", false, "compare instead of writing")
	flag.Parse()

	want := errorcodes.DocsTypeScript()

	if *check {
		have, err := os.ReadFile(*out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s cannot be read: %v\n", *out, err)
			os.Exit(1)
		}
		if string(have) != want {
			fmt.Fprintf(os.Stderr, "%s is out of date: run go run ./tools/errorcodes\n", *out)
			os.Exit(1)
		}
		fmt.Printf("%s is current (%d codes)\n", *out, len(errorcodes.All()))
		return
	}

	if err := os.WriteFile(*out, []byte(want), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("wrote %s (%d codes)\n", *out, len(errorcodes.All()))
}
