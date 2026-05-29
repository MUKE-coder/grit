package models

import "time"

type Business struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	OwnerID   uint      `gorm:"not null;index" json:"owner_id"`
	Owner     User      `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Branches  []Branch  `gorm:"foreignKey:BusinessID" json:"branches,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
