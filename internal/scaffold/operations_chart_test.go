package scaffold

import (
	"strings"
	"testing"
)

// The Operations charts are hand-drawn SVG, so their arithmetic is the thing
// that breaks, and it breaks silently.
//
// points() maps a series onto a polyline. The x coordinate divides by
// values.length - 1, which is zero when a fresh server has been polled once,
// and 0/0 is NaN. The browser then rejects the entire polyline with
// "<polyline> attribute points: Expected number, NaN,96.0" and draws nothing:
// no crash, no empty state, just a panel with no line in it and an error in a
// console nobody has open. The y axis already guarded a span of zero. The x
// axis did not, and it took opening the page on a server with one reading to
// see it.
func TestOperationsChartSurvivesASeriesOfOne(t *testing.T) {
	page := adminPerformancePageV2()

	if !strings.Contains(page, `if (values.length === 0) return "";`) {
		t.Error("points() does not handle an empty series, where Math.max returns -Infinity")
	}
	if !strings.Contains(page, "values.length === 1 ? [values[0], values[0]] : values") {
		t.Error("points() still divides by values.length - 1 for a one-point series, which is 0/0: the polyline renders as NaN and the browser drops it")
	}
	// The divisor has to be the padded series, not the original, or the guard
	// above is bypassed by the very line it exists to protect.
	if strings.Contains(page, "(i / (values.length - 1))") {
		t.Error("the x coordinate still divides by the unpadded length")
	}
	if !strings.Contains(page, "(i / (series.length - 1))") {
		t.Error("the x coordinate no longer divides by the padded series")
	}

	// And the flat-line guard it was modelled on is still there.
	if !strings.Contains(page, "span === 0 ? height / 2") {
		t.Error("a series that has not moved is drawn on the floor of the box again, which reads as a drop to nothing")
	}
}
