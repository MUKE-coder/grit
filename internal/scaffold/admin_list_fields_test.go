package scaffold

import (
	"regexp"
	"strings"
	"testing"
)

// A list in a response is [] or it is a crash.
//
// AccessReview.Items was tagged json:"items,omitempty", so a review with
// nothing to certify sent no items key at all. The page did
// detail.items.filter(...) on undefined and died inside its error boundary.
// Opening the first review on a project with no role assignments was enough to
// do it, which is the first thing anybody does on that screen, and tsc could
// not see it: the Go tag and the TypeScript interface are two files nobody
// compares, and the interface said items was always there.
//
// omitempty is right for a scalar that means "unset". It is wrong for a list a
// client iterates, where empty is a real answer rather than an absence.
func TestAnEmptyAccessReviewDoesNotBreakItsPage(t *testing.T) {
	model := apiAccessReviewModelGo()
	service := apiAccessReviewServiceGo()
	page := adminAccessReviewPage()

	if strings.Contains(model, `json:"items,omitempty"`) {
		t.Error("AccessReview.Items is omitempty again: a review with nothing to certify " +
			"sends no items key at all, and the page reads it as a list")
	}
	if !strings.Contains(model, `json:"items"`) {
		t.Error("AccessReview.Items has lost its json tag")
	}

	// A nil slice marshals to null, which breaks a client exactly as hard as
	// a missing key, so the service hands back an empty one.
	if !strings.Contains(service, "review.Items = []models.AccessReviewItem{}") {
		t.Error("the detail service no longer normalises a nil Items slice, so a review " +
			"with no rows sends null instead of []")
	}

	// And the page does not reach into it blind. The API promising something
	// is not the same as the browser having received it: a deploy skew, a
	// cached response, a proxy that drops a field.
	if strings.Contains(page, "detail.items.") {
		t.Error("the access review page reads detail.items without a guard; a page should " +
			"not die of a shape")
	}
	if !strings.Contains(page, "const items = detail?.items ?? [];") {
		t.Error("the access review page no longer guards its items list")
	}
}

// Nothing iterates a list the API is allowed to leave out.
//
// Ticket.Replies carries the same omitempty, and the ticket thread is only
// safe because it happens to check `ticket.replies &&` first. That guard is
// one refactor away from being tidied into the map call it protects, and the
// day it is, a ticket nobody has answered yet stops rendering.
func TestATicketWithNoRepliesStillRenders(t *testing.T) {
	model := ticketModelGo()
	page := adminTicketThreadPage(tripleOptions())

	if !strings.Contains(model, `json:"replies,omitempty"`) {
		// If somebody drops the omitempty, the guard below stops being load
		// bearing and this test has nothing left to protect.
		t.Skip("Ticket.Replies is no longer omitempty, so the page cannot be handed undefined")
	}

	// Every read of .replies has to be guarded, because the field may be
	// absent entirely.
	reads := regexp.MustCompile(`\.replies\.(map|filter|forEach|some|reduce|slice|length)\b`)
	for _, hit := range reads.FindAllStringIndex(page, -1) {
		window := page[listGuardStart(hit[0], 200):hit[0]]
		guarded := strings.Contains(window, ".replies &&") ||
			strings.Contains(page[listGuardStart(hit[0], 2):hit[0]], "?.")
		if !guarded {
			t.Errorf("the ticket thread reads %q with no guard, and Ticket.Replies is "+
				"omitempty, so a ticket nobody has answered sends no replies key",
				page[hit[0]:hit[1]])
		}
	}
}

// listGuardStart is the index n characters before i, clamped at the string start.
func listGuardStart(i, n int) int {
	if i-n < 0 {
		return 0
	}
	return i - n
}
