package scaffold

import (
	"strings"
	"testing"
)

// The template and the repair have to produce the same logger.
//
// They are two copies of one block, which is the shape of bug this project
// keeps finding: a test over each side passes while the pair disagrees. So
// one test reads both.
func TestTheDBLoggerRepairAgreesWithTheTemplate(t *testing.T) {
	fresh := apiDatabaseGo()
	if !strings.Contains(fresh, "IgnoreRecordNotFoundError: true") {
		t.Fatal("a new project's database.go does not set IgnoreRecordNotFoundError, " +
			"so seeding still prints a red block for a lookup meant to miss")
	}
	if !strings.Contains(fresh, dbLoggerBlockNew) {
		t.Error("the repair writes a logger block the template does not contain, " +
			"so an upgraded project and a new one would differ")
	}
	if strings.Contains(fresh, dbLoggerLineOld) {
		t.Error("the template still carries the old one-line logger")
	}
}

// An old project gets the new block, and the result is what a new project has.
func TestTheDBLoggerRepairUpgradesAnOldFile(t *testing.T) {
	// The file as it was before the change: the current template with the new
	// block put back to the old line. Built from the real template rather than
	// hand-written, so this is a shape some project actually has.
	old := strings.Replace(apiDatabaseGo(), dbLoggerBlockNew, dbLoggerLineOld, 1)
	if !strings.Contains(old, dbLoggerLineOld) {
		t.Fatal("could not rebuild the pre-repair file")
	}

	out, fixed, warnings := repairDBLoggerSource(old)
	if len(warnings) != 0 {
		t.Fatalf("the repair warned instead of running: %v", warnings)
	}
	if len(fixed) != 1 {
		t.Fatalf("the repair reported %d changes, want 1", len(fixed))
	}
	if out != apiDatabaseGo() {
		t.Error("the repaired file is not byte-for-byte what a new project gets")
	}

	// Twice is once. An upgrade runs over projects that have already been
	// upgraded, so a repair that applies again doubles the block.
	again, fixedAgain, warnAgain := repairDBLoggerSource(out)
	if again != out || len(fixedAgain) != 0 || len(warnAgain) != 0 {
		t.Errorf("running the repair twice changed the file again (%d changes, %v)",
			len(fixedAgain), warnAgain)
	}
}

// A project whose database.go somebody has edited past recognition is warned
// about rather than mangled.
func TestTheDBLoggerRepairWarnsOnAnUnrecognisedFile(t *testing.T) {
	src := "package database\n\nvar x = 1\n"
	out, fixed, warnings := repairDBLoggerSource(src)
	if out != src || len(fixed) != 0 {
		t.Error("the repair edited a file it does not recognise")
	}
	if len(warnings) != 1 {
		t.Fatalf("got %d warnings, want 1 saying what to do by hand", len(warnings))
	}
	if !strings.Contains(warnings[0], "IgnoreRecordNotFoundError") {
		t.Errorf("the warning does not name the setting to add: %q", warnings[0])
	}
}

// The guard is on the setting, not on the prose that explains it.
//
// The comment the repair inserts contains the words "record not found", and a
// guard written over those words would match its own explanation and no-op on
// every project, silently. This proves the guard keys off something only the
// code has.
func TestTheDBLoggerGuardIsNotMatchingItsOwnComment(t *testing.T) {
	commentOnly := "// ErrRecordNotFound is never worth a log line, record not found\n" +
		dbLoggerLineOld + "\n" +
		`import ("log"; "os"; "time")` + "\n"
	out, fixed, _ := repairDBLoggerSource(commentOnly)
	if out == commentOnly || len(fixed) != 1 {
		t.Error("a file carrying the explanation but not the setting was skipped, " +
			"which means the guard is matching prose")
	}
}
