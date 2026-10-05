package services

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"commerce/apps/api/internal/models"
)

// The activity feed's create, update and delete rows are written here, in
// batches, rather than inside the request that caused them. Written inline,
// each was a transaction of its own in every create, update and delete, which
// on Postgres is three round trips before the response goes out.
//
// The row is built in the request, while its IP address and user agent can
// still be read, and handed to a writer. A burst larger than the queue writes
// inline rather than dropping rows. FlushActivity, called at shutdown, waits
// for what is still queued.
//
// Sign-ins, security events and other LogActivity calls are still written
// before the request returns.

const (
	activityQueueSize = 4096
	activityBatchSize = 256
	// activityLinger is how long the writer waits for more rows after the first
	// before it writes. Without it, steady traffic that never arrives at the
	// same instant is written one row at a time, which saves the request its
	// wait but not the database its commits.
	activityLinger = 50 * time.Millisecond
)

type activityWriter struct {
	db   *gorm.DB
	rows chan models.UserActivity
}

var (
	activityWriters sync.Map // *gorm.Config -> *activityWriter
	activityPending atomic.Int64
)

// queueActivity builds the row now and writes it soon.
func queueActivity(db *gorm.DB, c *gin.Context, args ActivityArgs) {
	row := activityRow(c, args)
	w := activityWriterFor(db)
	activityPending.Add(1)
	select {
	case w.rows <- row:
	default:
		w.write([]models.UserActivity{row})
	}
}

// activityWriterFor returns the writer for db's connection, starting it on
// first use. Keyed by the connection's config rather than the handle: a handle
// bound to a request's context is a new *gorm.DB every time, and its context
// ends with the request, before the row is written.
func activityWriterFor(db *gorm.DB) *activityWriter {
	if w, ok := activityWriters.Load(db.Config); ok {
		return w.(*activityWriter)
	}
	w := &activityWriter{
		db:   db.Session(&gorm.Session{NewDB: true, Context: context.Background()}),
		rows: make(chan models.UserActivity, activityQueueSize),
	}
	actual, loaded := activityWriters.LoadOrStore(db.Config, w)
	if !loaded {
		go w.run()
	}
	return actual.(*activityWriter)
}

func (w *activityWriter) run() {
	for first := range w.rows {
		batch := []models.UserActivity{first}
		linger := time.NewTimer(activityLinger)
	collect:
		for len(batch) < activityBatchSize {
			select {
			case row := <-w.rows:
				batch = append(batch, row)
			case <-linger.C:
				break collect
			}
		}
		linger.Stop()
		w.write(batch)
	}
}

func (w *activityWriter) write(rows []models.UserActivity) {
	defer activityPending.Add(-int64(len(rows)))
	err := w.db.CreateInBatches(rows, activityBatchSize).Error
	if err == nil {
		return
	}
	if len(rows) == 1 {
		log.Printf("activity: failed to write %s: %v", rows[0].Action, err)
		return
	}
	// One bad row must not take the rest of the batch with it.
	for i := range rows {
		if err := w.db.Create(&rows[i]).Error; err != nil {
			log.Printf("activity: failed to write %s: %v", rows[i].Action, err)
		}
	}
}

// FlushActivity waits until every queued activity row is written, or ctx ends.
// cmd/server/main.go calls it once the server has stopped taking requests.
func FlushActivity(ctx context.Context) {
	for activityPending.Load() > 0 {
		select {
		case <-ctx.Done():
			log.Printf("activity: %d rows were not written before shutdown", activityPending.Load())
			return
		case <-time.After(10 * time.Millisecond):
		}
	}
}
