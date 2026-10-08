package services_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"saas/apps/api/internal/models"
	"saas/apps/api/internal/services"
)

func TestToOCSFAuthentication(t *testing.T) {
	ev := services.ToOCSF(models.UserActivity{
		ID:        "act-1",
		UserID:    "user-9",
		Action:    "auth.login",
		Severity:  "info",
		Summary:   "ada@example.com signed in",
		IPAddress: "203.0.113.7",
		CreatedAt: time.Unix(1_700_000_000, 0),
	}, "test-app")

	assert.Equal(t, 3002, ev["class_uid"], "auth.login -> Authentication class")
	assert.Equal(t, 1, ev["activity_id"], "Logon activity")
	assert.Equal(t, 300201, ev["type_uid"], "class_uid*100 + activity_id")
	assert.Equal(t, 1, ev["status_id"], "success")
	assert.Equal(t, int64(1_700_000_000_000), ev["time"], "epoch millis")

	actor := ev["actor"].(map[string]interface{})
	assert.Equal(t, "user-9", actor["user"].(map[string]interface{})["uid"])
	src := ev["src_endpoint"].(map[string]interface{})
	assert.Equal(t, "203.0.113.7", src["ip"])
}

// A failed sign-in is the one action whose OCSF status must be Failure — that
// is what a SIEM alerts on.
func TestToOCSFFailedLogin(t *testing.T) {
	ev := services.ToOCSF(models.UserActivity{
		Action:   "auth.login_failed",
		Severity: "warn",
	}, "test-app")

	assert.Equal(t, 3002, ev["class_uid"])
	assert.Equal(t, 2, ev["status_id"], "Failure")
	assert.Equal(t, "Failure", ev["status"])
	assert.Equal(t, 3, ev["severity_id"], "warn -> Medium")
}

func TestToOCSFAccountChange(t *testing.T) {
	cases := map[string]int{
		"auth.register":  1, // Create
		"password.reset": 4, // Password Reset
		"user.delete":    6, // Delete
		"user.update":    13,
	}
	for action, wantActivity := range cases {
		ev := services.ToOCSF(models.UserActivity{Action: action}, "app")
		assert.Equal(t, 3001, ev["class_uid"], action+" -> Account Change")
		assert.Equal(t, wantActivity, ev["activity_id"], action)
	}
}

// An unknown action must still map — dropping it would blind the SIEM to
// exactly the novel event worth seeing.
func TestToOCSFUnknownActionStillMaps(t *testing.T) {
	ev := services.ToOCSF(models.UserActivity{Action: "invoice.approve"}, "app")
	assert.Equal(t, 6003, ev["class_uid"], "falls back to API Activity")
	assert.Equal(t, 0, ev["activity_id"], "Unknown activity, never dropped")

	ev2 := services.ToOCSF(models.UserActivity{Action: "invoice.create"}, "app")
	assert.Equal(t, 6003, ev2["class_uid"])
	assert.Equal(t, 1, ev2["activity_id"], "verb create -> API Activity Create")
}

// System events have no actor; emitting a blank uid would read as a real
// account to a SIEM.
func TestToOCSFSystemEventHasNoActor(t *testing.T) {
	ev := services.ToOCSF(models.UserActivity{Action: "cron.cleanup", Summary: "nightly purge"}, "app")
	_, hasActor := ev["actor"]
	assert.False(t, hasActor, "no actor key when UserID is empty")
}

func TestToOCSFPreservesNativeAction(t *testing.T) {
	ev := services.ToOCSF(models.UserActivity{
		Action:       "ticket.close",
		ResourceType: "ticket",
		ResourceID:   "t-42",
	}, "app")
	um := ev["unmapped"].(map[string]interface{})
	assert.Equal(t, "ticket.close", um["grit_action"])
	assert.Equal(t, "ticket", um["grit_resource_type"])

	md := ev["metadata"].(map[string]interface{})
	require.Equal(t, services.OCSFVersion, md["version"])
}
