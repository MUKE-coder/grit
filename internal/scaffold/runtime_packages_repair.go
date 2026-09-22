package scaffold

import "strings"

// runtimePackagesMarker is where a plugin adds packages to the API's runtime
// image: ffmpeg, for grit plugin add video. A Dockerfile from before it gets the
// marker on upgrade, after the last package line Grit wrote.
const runtimePackagesMarker = "# Packages a plugin needs at runtime go here, such as ffmpeg for grit plugin\n" +
	"# add video.\n" +
	"# grit:runtime-packages\n"

// runtimePackagesAnchors are the lines the marker follows, most specific
// first: the vips line in the API image, then the base packages of either one.
var runtimePackagesAnchors = []string{
	"RUN if [ \"$IMAGE_BACKEND\" = \"vips\" ]; then apk --no-cache add vips; fi\n",
	"RUN apk --no-cache add ca-certificates tzdata\n",
	"RUN apk add --no-cache ca-certificates tzdata\n",
}

func repairRuntimePackagesMarkerSource(src string) (string, []string, []string) {
	if strings.Contains(src, "# grit:runtime-packages") {
		return src, nil, nil
	}
	for _, anchor := range runtimePackagesAnchors {
		if strings.Count(src, anchor) == 1 {
			return strings.Replace(src, anchor, anchor+"\n"+runtimePackagesMarker, 1),
				[]string{"the Dockerfile has a place for plugins' runtime packages (# grit:runtime-packages)"}, nil
		}
	}
	return src, nil, []string{"the Dockerfile is not the one Grit wrote: add a line \"# grit:runtime-packages\" in the runtime stage, after its apk add, so plugins such as video can add their packages"}
}

// thenRepair runs two source repairs in turn and joins what they report.
func thenRepair(first, second func(string) (string, []string, []string)) func(string) (string, []string, []string) {
	return func(src string) (string, []string, []string) {
		out, changes, warnings := first(src)
		out, more, moreWarnings := second(out)
		return out, append(changes, more...), append(warnings, moreWarnings...)
	}
}
