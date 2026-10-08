package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"saas/apps/api/internal/backup"
	"saas/apps/api/internal/config"
	"saas/apps/api/internal/database"
	"saas/apps/api/internal/storage"
)

// Backs up every registered model to a ZIP (CSV per table + dump.sql +
// metadata.json). By default it uploads to the storage the API uses, or to the
// "backups" disk when STORAGE_DISKS defines one, and records the row;
// --output writes a local file instead and touches nothing else.
func main() {
	out := flag.String("output", "", "Write the archive to this local path instead of uploading")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	svc := &backup.Service{DB: db}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	if *out != "" {
		f, err := os.Create(*out)
		if err != nil {
			log.Fatalf("Failed to create %s: %v", *out, err)
		}
		man, err := svc.ArchiveTo(ctx, f)
		if err != nil {
			f.Close()
			log.Fatalf("Backup failed: %v", err)
		}
		if err := f.Close(); err != nil {
			log.Fatalf("Failed to write %s: %v", *out, err)
		}
		var sizeKB float64
		if fi, err := os.Stat(*out); err == nil {
			sizeKB = float64(fi.Size()) / 1024
		}
		fmt.Printf("Backup written to %s: %d tables, %d rows, %.1f KB\n",
			*out, len(man.Tables), man.TotalRows, sizeKB)
		return
	}

	// The stores the API opens, from the same settings, so the archive lands
	// where the Data & Backup page looks for it: STORAGE_DRIVER=local included,
	// and the backups disk when there is one.
	storage.SetPublicPrefixes(cfg.StoragePublicPrefixes)
	if st, err := storage.Open(cfg.StorageDriver, cfg.Storage, storage.LocalConfig{
		Root:      cfg.Storage.LocalRoot,
		PublicURL: cfg.Storage.PublicURL,
		Secret:    cfg.Storage.URLSecret,
	}); err != nil {
		log.Printf("Default storage unavailable: %v", err)
	} else {
		svc.Storage = st
	}
	for _, named := range cfg.StorageDisks {
		if named.Name != backup.DiskName {
			continue
		}
		st, err := storage.Open(named.Driver, named.Storage, storage.LocalConfig{
			Root:      named.Storage.LocalRoot,
			PublicURL: named.Storage.PublicURL,
			Secret:    named.Storage.URLSecret,
		})
		if err != nil {
			log.Fatalf("The %q storage disk is not available: %v", named.Name, err)
		}
		storage.Disks.Add(named.Name, st)
	}
	store := svc.Store()
	if store == nil {
		log.Fatal("Storage is not configured\n(use --output <file> to write a local archive)")
	}

	rec, err := svc.Generate(ctx, "CLI")
	if err != nil {
		log.Fatalf("Backup failed: %v", err)
	}
	fmt.Printf("Backup %s uploaded to %s: %d tables, %d rows, %.1f KB\n",
		rec.ID, store.Describe(), rec.TableCount, rec.RowCount, float64(rec.SizeBytes)/1024)
}
