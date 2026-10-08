package handlers

import (
	"context"
	"encoding/csv"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"saas/apps/api/internal/respond"
)

// Import kicks off a BACKGROUND CSV import of usage_records. It streams the upload
// to a temp file (so a large file never sits in memory), starts an ImportJob,
// then hands the file to the service in a goroutine and returns 202
// immediately. Poll GET /imports/:id for progress and the result.
func (h *UsageRecordHandler) Import(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		respond.Fail(c, respond.CodeInvalidFile, "No CSV file provided")
		return
	}
	defer file.Close()

	// Stream the upload to a temp file: never ReadAll a large CSV into memory.
	tmp, err := os.CreateTemp("", "grit-import-*.csv")
	if err != nil {
		respond.Fail(c, respond.CodeTempError, "Could not buffer the upload")
		return
	}
	tmpPath := tmp.Name()
	if _, err := io.Copy(tmp, file); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		respond.Fail(c, respond.CodeInvalidCSV, "Could not read the upload")
		return
	}
	tmp.Close()

	// Count data rows up front (streaming) so the client's progress bar has a
	// denominator without holding the file in memory.
	total, err := countCSVRowsUsageRecord(tmpPath)
	if err != nil || total < 0 {
		os.Remove(tmpPath)
		respond.Fail(c, respond.CodeInvalidCSV, "Could not read the CSV file")
		return
	}

	job, err := h.service().StartImport(h.ctx(c), total)
	if err != nil {
		os.Remove(tmpPath)
		respond.Fail(c, respond.CodeJobError, "Could not start import")
		return
	}

	// In the background, so a large file never blocks the request. The context
	// keeps the caller and the organization and drops the cancellation: the
	// request is over as soon as this returns, and the import has only begun.
	go h.service().ImportCSV(context.WithoutCancel(h.ctx(c)), job.ID, tmpPath)

	c.JSON(http.StatusAccepted, gin.H{
		"data":    gin.H{"job_id": job.ID, "total": total},
		"message": "Import started",
	})
}

// countCSVRowsUsageRecord counts data rows (excluding the header) without holding
// the file in memory.
func countCSVRowsUsageRecord(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return -1, err
	}
	defer f.Close()
	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	reader.ReuseRecord = true
	n := 0
	for i := 0; ; i++ {
		_, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return -1, err
		}
		if i == 0 {
			continue // header row
		}
		n++
	}
	return n, nil
}

// Template returns a ready-to-fill CSV template (header row) for importing usage_records.
// belongs_to columns use the related record's natural key (e.g. "category"), or
// its id column ("<relation>_id") when the related model has no natural key.
func (h *UsageRecordHandler) Template(c *gin.Context) {
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", `attachment; filename="usage-records-template.csv"`)
	c.String(http.StatusOK, "subscription,metric,quantity,user\n")
}
