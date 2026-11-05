package scaffold

import "strings"

// Cutting the Blog example out of a template that still contains it.
//
// The landing page is one long raw string with four Blog regions in it, each
// already delimited by a grit:home:blog-* marker because `grit remove resource
// Blog` cuts them out of the file on disk. Splitting the string into nine
// pieces to make them conditional would make the template unreadable for no
// gain, so the template keeps its markers and the writer cuts the same regions
// the removal would (grit#173).
//
// The point of doing it here rather than there: a file the removal edits reads
// as edited for ever after, and `grit upgrade` leaves an edited file alone. So
// the blank starter was choosing between a project with no 404s and a project
// that still gets framework fixes to its landing page. Cutting at write time
// means the file on disk is what the template produces, the manifest hash
// matches, and both properties hold.

// withoutRegion removes everything from the line holding start to the line
// holding end, inclusive. A missing marker leaves src alone, because a template
// whose shape has changed should not be half-cut.
func withoutRegion(src, start, end string) string {
	lines := strings.Split(src, "\n")
	from, to := -1, -1
	for i, line := range lines {
		if from < 0 && strings.Contains(line, start) {
			from = i
			continue
		}
		if from >= 0 && strings.Contains(line, end) {
			to = i
			break
		}
	}
	if from < 0 || to < 0 {
		return src
	}
	return strings.Join(append(lines[:from:from], lines[to+1:]...), "\n")
}

// withoutLinesContaining drops every line holding any of the given substrings.
func withoutLinesContaining(src string, substrings ...string) string {
	lines := strings.Split(src, "\n")
	kept := make([]string, 0, len(lines))
	for _, line := range lines {
		drop := false
		for _, s := range substrings {
			if strings.Contains(line, s) {
				drop = true
				break
			}
		}
		if !drop {
			kept = append(kept, line)
		}
	}
	return strings.Join(kept, "\n")
}

// withoutBlogOnTheLandingPage cuts the four Blog regions out of the web landing
// page: the import, the data hook, the hero's call to action and the Recent
// Posts section.
//
// The same four regions `grit remove resource Blog` cuts, by the same markers,
// so a project that removes the Blog later and one scaffolded without it end up
// with the same file.
func withoutBlogOnTheLandingPage(src string) string {
	src = withoutRegion(src, "grit:home:blog-start", "grit:home:blog-end")
	src = withoutRegion(src, "grit:home:blog-hook-start", "grit:home:blog-hook-end")
	// The hero's primary call to action is "Read the Blog", which without a
	// blog is a link to a 404. Marked separately because it sits in the middle
	// of the hero rather than in a section of its own.
	src = withoutRegion(src, "grit:home:blog-cta-start", "grit:home:blog-cta-end")
	src = withoutLinesContaining(src,
		"// grit:home:blog-import",
		`from "@/lib/blog-api"`,
		`from "@/hooks/use-blogs"`,
	)
	return src
}
