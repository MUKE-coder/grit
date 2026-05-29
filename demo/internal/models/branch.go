package models

import "time"

type Branch struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	BusinessID uint      `gorm:"not null;index" json:"business_id"`
	Name       string    `gorm:"not null" json:"name"`
	Address    string    `json:"address"`
	IsDefault  bool      `gorm:"default:false" json:"is_default"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
