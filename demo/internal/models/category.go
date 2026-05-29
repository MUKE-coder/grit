package models

import "time"

type Category struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	BusinessID uint      `gorm:"not null;index" json:"business_id"`
	Name       string    `gorm:"not null" json:"name"`
	CreatedAt  time.Time `json:"created_at"`
}
