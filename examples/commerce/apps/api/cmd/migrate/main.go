package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/gorm"

	"commerce/apps/api/internal/config"
	"commerce/apps/api/internal/crypto"
	"commerce/apps/api/internal/database"
	"commerce/apps/api/internal/migrate"
	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/paginate"
	"commerce/apps/api/internal/sync"
)

func main() {
	fresh := flag.Bool("fresh", false, "Drop all tables before migrating")
	status := flag.Bool("status", false, "Show what each recorded run changed")
	down := flag.Bool("down", false, "Undo what the most recent run added")
	steps := flag.Int("steps", 1, "How many runs --down undoes")
	dryRun := flag.Bool("dry-run", false, "With --down, print the statements and change nothing")
	yes := flag.Bool("yes", false, "With --down, do not ask before dropping")
	gritVersion := flag.String("grit-version", "", "Recorded against the run")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// The history has to exist before anything reads or writes it, including on a
	// project that was scaffolded before down migrations existed. Its first run
	// there is recorded as a baseline, because the schema it finds was built by
	// runs nobody recorded and dropping all of it is not a rollback.
	if err := migrate.EnsureHistory(db); err != nil {
		log.Fatalf("%v", err)
	}

	switch {
	case *status:
		if err := showStatus(db); err != nil {
			log.Fatalf("%v", err)
		}
		return
	case *down:
		if err := rollback(db, *steps, *dryRun, *yes); err != nil {
			log.Fatalf("%v", err)
		}
		return
	}

	if *fresh {
		fmt.Println("Dropping all tables...")
		if err := database.DropAll(db); err != nil {
			log.Fatalf("Failed to drop tables: %v", err)
		}
		// The history described a schema that no longer exists.
		if err := migrate.Reset(db); err != nil {
			log.Fatalf("%v", err)
		}
		fmt.Println("All tables dropped.")
	}

	before, err := migrate.Snapshot(db)
	if err != nil {
		log.Fatalf("Failed to read the current schema: %v", err)
	}

	// Empty webhook external ids become NULL before the unique index is built.
	//
	// The unique index on (provider, external_id) used to be declared in a
	// method nothing called, so a project from before that was fixed can hold
	// several rows with external_id = "", and no database will build a unique
	// index over repeated empty strings. NULL repeats freely, and an event the
	// provider gave no id was never deduplicable anyway.
	//
	// Here rather than inside models.Migrate because this file is rewritten by
	// grit upgrade and models/user.go is not: that one holds the model registry
	// people add to. A migration fix living in a file upgrades never touch is a
	// fix only new projects get.
	if db.Migrator().HasTable("webhook_events") {
		if err := db.Exec("UPDATE webhook_events SET external_id = NULL WHERE external_id = ''").Error; err != nil {
			log.Printf("migrate: clearing empty webhook external ids (the unique index may refuse them): %v", err)
		}
	}

	fmt.Println("Running migrations...")
	if err := models.Migrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	// Postgres gets a trigram index on every column tagged search:"trigram", the
	// ones list search reads.
	if err := paginate.EnsureSearchIndexes(db, models.Models()...); err != nil {
		log.Printf("Search indexes were not created, so search reads the whole table: %v", err)
	}
	// Rows soft-deleted before a delete also moved updated_at get it now, so
	// sync pull, which pages on updated_at, carries those deletes too.
	if err := sync.BackfillUpdatedAt(db, models.Models()...); err != nil {
		log.Printf("Some soft-deleted rows were not brought up to date for sync: %v", err)
	}
	// Values in encrypted columns written before FIELD_ENCRYPTION_KEY was set,
	// two-factor secrets among them, are encrypted now. Without a key this
	// changes nothing.
	if n, err := crypto.EncryptExisting(db, models.Models()...); err != nil {
		log.Printf("Some values in encrypted columns are still stored in the clear: %v", err)
	} else if n > 0 {
		fmt.Printf("Encrypted %d value(s) stored before field encryption was on.\n", n)
	}

	after, err := migrate.Snapshot(db)
	if err != nil {
		// The migration itself worked. Only the record of it did not, and
		// reporting failure here would read as a failed migration.
		fmt.Printf("Migrations completed, but the schema could not be read back, so this run was not recorded: %v\n", err)
		os.Exit(0)
	}

	changes := migrate.Diff(before, after)
	if len(changes) == 0 {
		fmt.Println("Migrations completed successfully. The schema was already up to date.")
		os.Exit(0)
	}

	// A baseline is the run that found an empty database, whether that is a new
	// project or an old one meeting the history for the first time.
	baseline := len(before) == 0
	run, err := migrate.Record(db, changes, *gritVersion, baseline)
	if err != nil {
		fmt.Printf("Migrations completed, but recording the run failed, so it cannot be rolled back: %v\n", err)
		os.Exit(0)
	}

	fmt.Printf("Migrations completed successfully. Run %s added %d change(s):\n", run.ID, len(changes))
	for _, change := range changes {
		fmt.Println("  + " + change.String())
	}
	if baseline {
		fmt.Println("\nThis is the baseline run, the one that built the schema. It cannot be rolled back: use grit migrate --fresh to start over.")
	} else {
		fmt.Println("\nUndo it with: grit migrate down")
	}
	os.Exit(0)
}

// showStatus prints the recorded runs, newest first, with what each one added.
func showStatus(db *gorm.DB) error {
	runs, err := migrate.History(db, 20)
	if err != nil {
		return err
	}
	if len(runs) == 0 {
		fmt.Println("No migration runs recorded yet. Run grit migrate.")
		return nil
	}
	for _, run := range runs {
		changes, err := migrate.ChangesOf(db, run.ID)
		if err != nil {
			return err
		}
		state := "applied"
		if run.Baseline {
			state = "baseline, cannot be rolled back"
		}
		if run.RolledBackAt != nil {
			state = "rolled back " + run.RolledBackAt.Format("2006-01-02 15:04")
		}
		version := run.GritVersion
		if version == "" {
			version = "unknown"
		}
		fmt.Printf("%s  %s  grit %s  %d change(s)  [%s]\n",
			run.ID, run.AppliedAt.Format("2006-01-02 15:04"), version, len(changes), state)
		for _, change := range changes {
			fmt.Println("    " + change.String())
		}
	}
	return nil
}

// rollback undoes the most recent runs that have not been undone already.
//
// It prints every statement before running any of them, and stops at the
// baseline rather than dropping the schema out from under the application.
func rollback(db *gorm.DB, steps int, dryRun, yes bool) error {
	if steps < 1 {
		steps = 1
	}
	runs, err := migrate.Undone(db, steps)
	if err != nil {
		return err
	}
	if len(runs) == 0 {
		fmt.Println("Nothing to roll back: every recorded run has already been undone.")
		return nil
	}

	type plan struct {
		run        migrate.Run
		changes    []migrate.Change
		statements []string
	}
	var plans []plan
	for _, run := range runs {
		if run.Baseline {
			if len(plans) == 0 {
				return fmt.Errorf("run %s is the baseline, the run that built this schema: "+
					"rolling it back would drop the database. Use grit migrate --fresh to start over", run.ID)
			}
			fmt.Printf("Stopping at run %s, the baseline that built the schema.\n", run.ID)
			break
		}
		changes, err := migrate.ChangesOf(db, run.ID)
		if err != nil {
			return err
		}
		plans = append(plans, plan{run: run, changes: changes, statements: migrate.Statements(db, changes)})
	}

	fmt.Printf("Rolling back %d run(s):\n", len(plans))
	for _, item := range plans {
		fmt.Printf("\n  %s (%s, %d change(s))\n",
			item.run.ID, item.run.AppliedAt.Format("2006-01-02 15:04"), len(item.changes))
		for _, statement := range item.statements {
			fmt.Println("    " + statement + ";")
		}
	}
	if dryRun {
		fmt.Println("\n--dry-run: nothing was changed.")
		return nil
	}

	fmt.Println("\nDropping a table or a column does not keep the data in it. This is not reversible.")
	if !yes && !confirmed() {
		fmt.Println("Cancelled. Nothing was changed.")
		return nil
	}

	for _, item := range plans {
		applied, err := migrate.Rollback(db, item.run, item.changes, false)
		for _, statement := range applied {
			fmt.Println("  - " + statement)
		}
		if err != nil {
			return fmt.Errorf("rolling back %s: %w", item.run.ID, err)
		}
		fmt.Printf("Run %s rolled back.\n", item.run.ID)
	}
	return nil
}

// confirmed asks before dropping. Anything other than yes cancels, and so does a
// closed stdin, which is what a CI job or a pipe looks like: there --yes is the
// way to mean it.
func confirmed() bool {
	fmt.Print("Type yes to continue: ")
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil && strings.TrimSpace(line) == "" {
		fmt.Println()
		return false
	}
	return strings.TrimSpace(strings.ToLower(line)) == "yes"
}
