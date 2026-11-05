// Command docsamples prints the runnable docs blocks as shell scripts.
//
// The docs mark a block with verify="<flow>" to claim the commands in it work.
// This turns those claims into something a CI job can execute: -list names the
// flows, and -flow <id> prints the script for one, in page order.
//
//	go run ./tools/docsamples -list
//	go run ./tools/docsamples -flow migrate-rollback > /tmp/flow.sh
//	go run ./tools/docsamples -check
//
// A flow comes from one of two places: a verify="<flow>" block in a page, or a
// file in docs/content/verified. The second exists because adding a verify=
// block means editing a React component, and executed coverage stayed at four
// blocks while the docs grew to 193 pages (grit#168). This tool does not care
// which a flow came from.
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
	flowDir := flag.String("flows", docscheck.FlowFileDir, "the directory of declared flow files")
	flow := flag.String("flow", "", "print the script for this flow")
	list := flag.Bool("list", false, "list the flows and how many blocks each has")
	check := flag.Bool("check", false, "report a flow file whose page: does not exist, and exit non-zero")
	flag.Parse()

	if *check {
		dangling, err := docscheck.FlowFilePages(*docs, *flowDir)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, problem := range dangling {
			fmt.Fprintln(os.Stderr, problem)
		}
		if len(dangling) > 0 {
			os.Exit(1)
		}
		fmt.Println("every declared flow names a page that exists")
		return
	}

	flows, err := docscheck.AllFlows(*docs, *flowDir)
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
