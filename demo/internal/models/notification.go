package models

import "time"

// Notification is an in-app message surfaced through the admin bell.
// Source distinguishes Sentinel security findings from Pulse perf
// findings from manual operator messages.
type Notification struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"index" json:"user_id"` // 0 = visible to all admins
	Source    string     `gorm:"size:16;index" json:"source"`   // sentinel | pulse | system
	Severity  string     `gorm:"size:16;index" json:"severity"`  // critical | high | medium | low | info
	Title     string     `gorm:"size:200" json:"title"`
	Body      string     `gorm:"type:text" json:"body"`
	Link      string     `gorm:"size:500" json:"link"`
	Dedup     string     `gorm:"size:128;uniqueIndex" json:"-"`
	Count     int        `gorm:"default:1" json:"count"`
	ReadAt    *time.Time `json:"read_at"`
	CreatedAt time.Time  `gorm:"index" json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
