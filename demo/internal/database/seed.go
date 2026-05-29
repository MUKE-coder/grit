package database

import (
	"errors"
	"fmt"
	"log"

	"gritdemo/internal/models"

	"gorm.io/gorm"
)

const gritDemoAdminEmail = "admin@grit.demo"
const gritDemoAdminPassword = "password123"
const gritMotorsBusinessName = "Grit Motors"

// Seed makes sure the deployment has the canonical Grit admin user, business,
// and a default branch. Runs on every server boot.
//
// Two paths:
//   - Empty DB: full first-time setup (admin + business + default branch + role).
//   - Non-empty DB (e.g. legacy Nakawa Fashion data still around): only top up
//     what's missing — guarantees admin@grit.demo always works without
//     touching existing rows.
//
// The result is: no matter what state the DB is in, the documented admin
// credentials log in successfully.
func Seed(db *gorm.DB) error {
	// Ensure the canonical admin exists.
	var admin models.User
	err := db.Where("email = ?", gritDemoAdminEmail).First(&admin).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("looking up admin: %w", err)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		admin = models.User{Name: "Admin", Email: gritDemoAdminEmail}
		if err := admin.SetPassword(gritDemoAdminPassword); err != nil {
			return fmt.Errorf("hashing admin password: %w", err)
		}
		if err := db.Create(&admin).Error; err != nil {
			return fmt.Errorf("creating admin: %w", err)
		}
		log.Printf("seed: created %s / %s", gritDemoAdminEmail, gritDemoAdminPassword)
	}

	// Ensure a Business exists for the admin to log into. Prefer one that
	// already lists the admin as owner; otherwise pick any Business; otherwise
	// create the Grit Motors business from scratch.
	var business models.Business
	if err := db.Where("owner_id = ?", admin.ID).First(&business).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("looking up admin's business: %w", err)
		}
		// Try any existing business (legacy Nakawa Fashion case).
		if err := db.First(&business).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("looking up any business: %w", err)
			}
			business = models.Business{Name: gritMotorsBusinessName, OwnerID: admin.ID}
			if err := db.Create(&business).Error; err != nil {
				return fmt.Errorf("creating business: %w", err)
			}
			log.Printf("seed: created business %q", business.Name)
		}
	}

	// Ensure the admin has an admin role on that business.
	var role models.UserBusinessRole
	err = db.Where("user_id = ? AND business_id = ?", admin.ID, business.ID).First(&role).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("looking up admin role: %w", err)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := db.Create(&models.UserBusinessRole{
			UserID: admin.ID, BusinessID: business.ID, Role: models.RoleAdmin,
		}).Error; err != nil {
			return fmt.Errorf("creating admin role: %w", err)
		}
		log.Printf("seed: granted admin role on %q", business.Name)
	}

	// Ensure at least one branch exists for the business.
	var branchCount int64
	db.Model(&models.Branch{}).Where("business_id = ?", business.ID).Count(&branchCount)
	if branchCount == 0 {
		if err := db.Create(&models.Branch{
			BusinessID: business.ID,
			Name:       "Main Branch",
			IsDefault:  true,
		}).Error; err != nil {
			return fmt.Errorf("creating default branch: %w", err)
		}
		log.Println("seed: created default Main Branch")
	}

	return nil
}
