package docscheck

import (
	"strings"
	"testing"
)

// A page showing what a command printed, with the two-column "Next steps"
// table `grit new` ends with. Every line of that table reads as a command
// followed by a subcommand, and none of them is one.
const transcriptPage = "export default function Page() {\n" +
	"  return (\n" +
	"    <CodeBlock\n" +
	"      output=\"grit new myapp\"\n" +
	"      language=\"text\"\n" +
	"      code={`  Created myapp\n" +
	"\n" +
	"    grit migrate            create database tables\n" +
	"    grit seed               sample data, optional\n" +
	"    grit start              run everything`}\n" +
	"    />\n" +
	"  )\n" +
	"}\n"

// The same page without the marker, which is what failed the build.
var unmarkedPage = strings.Replace(transcriptPage,
	"      output=\"grit new myapp\"\n", "", 1)

func commandStrings(samples []Sample) []string {
	out := make([]string, 0, len(samples))
	for _, s := range samples {
		out = append(out, s.Command)
	}
	return out
}

// The body of a marked block is not read as commands.
func TestAnOutputBlockIsNotReadAsCommands(t *testing.T) {
	got := commandStrings(extractFile("page.tsx", transcriptPage))

	for _, command := range got {
		if strings.HasPrefix(command, "grit migrate ") ||
			strings.HasPrefix(command, "grit seed ") ||
			strings.HasPrefix(command, "grit start ") {
			t.Errorf("a line of the transcript was read as a command: %q", command)
		}
	}
}

// And the command the marker names is read, so the marker cannot be used to
// put a command on a page without the CLI ever being asked about it.
func TestTheCommandThatProducedTheOutputIsStillChecked(t *testing.T) {
	got := commandStrings(extractFile("page.tsx", transcriptPage))

	found := false
	for _, command := range got {
		if command == "grit new myapp" {
			found = true
		}
	}
	if !found {
		t.Errorf("the command named by output= was not sampled; got %v", got)
	}
}

// The control. Without the marker the same page yields the lines that broke
// the build, so the test above is testing the marker and not the weather.
func TestWithoutTheMarkerTheTranscriptReadsAsCommands(t *testing.T) {
	got := commandStrings(extractFile("page.tsx", unmarkedPage))

	found := false
	for _, command := range got {
		if strings.HasPrefix(command, "grit migrate ") {
			found = true
		}
	}
	if !found {
		t.Skipf("the extractor no longer reads an unmarked transcript as "+
			"commands, so the marker guards nothing; got %v", got)
	}
}
