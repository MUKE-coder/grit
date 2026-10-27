package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/MUKE-coder/grit/v3/internal/dbindex"
)

func indexCmd() *cobra.Command {
	var allSortable bool
	var write string
	var only []string

	cmd := &cobra.Command{
		Use:   "index",
		Short: "The indexes this project's list queries need, as SQL",
		Long: `Works out which indexes the generated lists need, and writes the SQL.

Every list is WHERE <scope> ORDER BY <column>, id. Without an index in that
shape the database sorts the whole table on disk for every page. It is
invisible on a seeded table and the slowest thing in the product on a real one:
measured at 2.3 seconds for one page of 500,000 rows, and 16.6 seconds at page
20,000.

The index has three parts, and the project already says what they are. The
scope first, because it is an equality the planner can seek to: org_id for a
tenant-owned model, the owner column for an --owned-by one, nothing otherwise.
Then the sort column. Then id, because every list orders by the column and then
by id to break ties, and an index without that trailing column cannot serve the
whole ORDER BY.

By default only the default sort gets an index, because that is the one every
list request uses when nobody has clicked a column. The rest cost write
throughput on tables that are mostly written, so you choose them with a
measurement in hand.

It writes SQL and creates nothing. The statements are CONCURRENTLY, which does
not lock the table and therefore cannot run inside a transaction, so apply the
file with psql rather than through a runner that wraps everything in BEGIN.

  grit index
  grit index --all-sortable
  grit index --only deal=title,status --only contact=name
  grit index --write db/indexes.sql`,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := findProjectRoot()
			if err != nil {
				return err
			}
			apiRoot, err := findAPIDir()
			if err != nil {
				return err
			}

			resources, err := dbindex.Scan(apiRoot)
			if err != nil {
				return err
			}
			if len(resources) == 0 {
				cMuted.Println("\n  No list-bearing resource found. Generate one first.")
				return nil
			}

			picked := map[string][]string{}
			for _, o := range only {
				name, cols, found := strings.Cut(o, "=")
				if !found {
					return fmt.Errorf("--only wants resource=col,col, got %q", o)
				}
				picked[strings.ToLower(strings.TrimSpace(name))] =
					strings.Split(strings.ReplaceAll(cols, " ", ""), ",")
			}

			indexes := dbindex.Plan(resources, allSortable, picked)
			sql := dbindex.SQL(indexes, "v"+version)

			if write != "" {
				path := filepath.Join(root, filepath.FromSlash(write))
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					return err
				}
				if err := os.WriteFile(path, []byte(sql), 0o644); err != nil {
					return err
				}
				cOK.Printf("\n  ✓ %s\n", write)
				cMuted.Printf("  %d indexes across %d resources.\n\n", len(indexes), len(resources))
				cMuted.Println("  Apply it:")
				cCmd.Printf("    psql \"$DATABASE_URL\" -f %s\n\n", write)
				return nil
			}

			fmt.Println()
			fmt.Print(sql)
			cMuted.Printf("  %d indexes across %d resources. --write <path> to save this.\n\n",
				len(indexes), len(resources))
			return nil
		},
	}

	cmd.Flags().BoolVar(&allSortable, "all-sortable", false,
		"An index for every sortable column, not just the default sort")
	cmd.Flags().StringVar(&write, "write", "", "Write the SQL to this path instead of printing it")
	cmd.Flags().StringArrayVar(&only, "only", nil,
		"Pick the columns for one resource: --only deal=title,status (repeatable)")

	return cmd
}
