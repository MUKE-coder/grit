package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/MUKE-coder/grit/v3/internal/doctor"
)

// runEnvDoctor prints `grit doctor --env`: is this machine ready to run a
// Grit project?
//
// The audience is somebody who has not run grit new yet, so the output says
// what to do rather than what is wrong. Every failure carries the address of
// a page that explains it, because a terminal is a bad place to read four
// paragraphs and a good place to be given a link.
func runEnvDoctor(ctx context.Context, asJSON, skipNode, skipDocker bool) error {
	if ctx == nil {
		ctx = context.Background()
	}
	report := doctor.CheckEnv(ctx, doctor.EnvOptions{
		SkipNode:   skipNode,
		SkipDocker: skipDocker,
	})

	if asJSON {
		out, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(out))
		if report.Failures() > 0 {
			os.Exit(1)
		}
		return nil
	}

	red := color.New(color.FgHiRed)
	yellow := color.New(color.FgHiYellow)
	green := color.New(color.FgHiGreen)
	gray := color.New(color.FgHiBlack)
	cyan := color.New(color.FgHiCyan)

	fmt.Println()
	cyan.Printf("  %d checks on this machine\n", len(report.Findings))
	fmt.Println()

	for _, f := range report.Findings {
		switch f.Level {
		case doctor.EnvOK:
			green.Printf("  ✓ %-16s", f.Name)
			gray.Printf("%s\n", f.Detail)
		case doctor.EnvWarn:
			yellow.Printf("  ⚠ %-16s", f.Name)
			fmt.Printf("%s\n", f.Detail)
		default:
			red.Printf("  ✗ %-16s", f.Name)
			fmt.Printf("%s\n", f.Detail)
		}
		if f.Problem != nil {
			gray.Printf("      %s\n", f.Problem.Fix)
			gray.Printf("      %s\n", f.Problem.URL())
		}
	}

	fmt.Println()
	switch {
	case report.Failures() > 0:
		gray.Printf("  %d blocking, %d worth knowing\n\n", report.Failures(), report.Warnings())
		os.Exit(1)
	case report.Warnings() > 0:
		gray.Printf("  Nothing blocking. %d worth knowing about above.\n\n", report.Warnings())
	default:
		green.Println("  Ready. Run: grit new myapp")
		fmt.Println()
	}
	return nil
}
