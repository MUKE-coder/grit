// Command docsamples prints the runnable docs blocks as shell scripts.
//
// The docs mark a block with verify="<flow>" to claim the commands in it work.
// This turns those claims into something a CI job can execute: -list names the
// flows, and -flow <id> prints the script for one, in page order.
//
//	go run ./tools/docsamples -list
//	go run ./tools/docsamples -flow migrate-rollback > /tmp/flow.sh
//
// It is a tool rather than part of the grit binary because it is about this
// repository's docs, not about anybody's project.
package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/MUKE-coder/grit/v3/internal/docscheck"
)

func main() {
	docs := flag.String("docs", "docs/app", "the docs site's app directory")
	flow := flag.String("flow", "", "print the script for this flow")
	list := flag.Bool("list", false, "list the flows and how many blocks each has")
	flag.Parse()

	flows, err := docscheck.Flows(*docs)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	names := make([]string, 0, len(flows))
	for name := range flows {
		names = append(names, name)
	}
	sort.Strings(names)

	switch {
	case *list:
		for _, name := range names {
			blocks := flows[name]
			fmt.Printf("%s\t%d block(s)\t%s\n", name, len(blocks), blocks[0].File)
		}
	case *flow != "":
		blocks, ok := flows[*flow]
		if !ok {
			fmt.Fprintf(os.Stderr, "no flow %q in the docs; known flows: %v\n", *flow, names)
			os.Exit(1)
		}
		fmt.Print(docscheck.Script(*flow, blocks))
	default:
		fmt.Fprintln(os.Stderr, "pass -list or -flow <id>")
		os.Exit(2)
	}
}
