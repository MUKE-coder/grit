package database

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ForRequest is the handle a read should use.
//
// The primary when this person wrote in the last few seconds, a replica
// otherwise. On a project with no replicas both are the same connection and
// this costs a map lookup.
//
// Generated list and detail handlers call it. Hand-written ones should too:
//
//	var posts []models.Post
//	database.ForRequest(c, h.DB).Find(&posts)
//
// Not for a read that decides a write. Those take database.Primary, or better,
// happen inside the transaction that does the writing.
func ForRequest(c *gin.Context, db *gorm.DB) *gorm.DB {
	if c == nil {
		return db
	}
	if pinned, ok := c.Get("read_your_writes"); ok && pinned == true {
		return Primary(db)
	}
	return db
}
