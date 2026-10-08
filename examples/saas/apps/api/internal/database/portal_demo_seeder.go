package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"

	"saas/apps/api/internal/jsontime"
	"saas/apps/api/internal/models"
	"saas/apps/api/internal/money"
)

// metered is one metric on a customer's plan, and how it was recorded.
type metered struct {
	name   string
	days   int
	perDay int
}

// SeedPortalDemo makes the example demonstrable: two customers who can sign in
// to the portal and see their own billing, and nobody else's.
//
// Hand-written, and the one seeder in this project that is. The generated
// seeders fill a table; this one builds the situation the example exists to
// show, which needs rows in five tables to agree with each other: a user, the
// role that lets them through the door, a plan, a subscription pointing at
// both, and the invoices and usage hanging off it.
//
// Two customers rather than one, deliberately. With a single customer the
// portal looks right whether or not the owner scope is doing anything, and the
// whole claim of this example is that Ada cannot see Basil's invoices. You can
// check it by hand: sign in as each and compare.
func SeedPortalDemo(db *gorm.DB) error {
	var existing int64
	db.Model(&models.Subscription{}).Count(&existing)
	if existing > 0 {
		log.Println("Portal demo already seeded, skipping...")
		return nil
	}

	// The role a customer gets. Not ADMIN and not an operator with a smaller
	// screen: a plain USER holding exactly the two permissions the portal's
	// pages read, which is what the API checks before the owner scope ever
	// runs.
	role := models.Role{
		Name:        "Customer",
		Description: "Reads their own subscription, invoices and usage in the portal",
	}
	if err := role.SetGrants([]string{
		"subscriptions.view",
		"invoices.view",
		"usage_records.view",
	}); err != nil {
		return fmt.Errorf("granting the customer role: %w", err)
	}
	if err := db.Where("name = ?", role.Name).FirstOrCreate(&role).Error; err != nil {
		return fmt.Errorf("creating the customer role: %w", err)
	}

	plans := []models.Plan{
		{Name: "Starter", Slug: "starter", Price: money.New(1900, "USD"), Interval: "monthly", Seats: 3,
			Blurb: "For one person and a side project.", Active: true},
		{Name: "Team", Slug: "team", Price: money.New(9900, "USD"), Interval: "monthly", Seats: 15,
			Blurb: "Shared workspaces, roles and audit history.", Active: true},
	}
	for i := range plans {
		if err := db.Where("slug = ?", plans[i].Slug).FirstOrCreate(&plans[i]).Error; err != nil {
			return fmt.Errorf("seeding plan %s: %w", plans[i].Slug, err)
		}
	}

	customers := []struct {
		email  string
		first  string
		last   string
		plan   *models.Plan
		seats  int
		status string
		// One entry per metric: how many days it was recorded on, and roughly
		// how much each day. Several rows per metric rather than one, because
		// the usage page counts rows per value with ?breakdown=metric and a
		// single row per metric makes every card read "1".
		metrics []metered
	}{
		{"ada@example.com", "Ada", "Okello", &plans[1], 9, "active", []metered{
			{"api_calls", 12, 1800},
			{"exports", 4, 2},
			{"seats_used", 2, 9},
		}},
		{"basil@example.com", "Basil", "Mwangi", &plans[0], 2, "trialing", []metered{
			{"api_calls", 3, 240},
			{"exports", 1, 1},
		}},
	}

	now := time.Now()
	for _, c := range customers {
		user := models.User{
			FirstName: c.first,
			LastName:  c.last,
			Email:     c.email,
			Password:  "customer123",
			Role:      models.RoleUser,
			Active:    true,
		}
		if err := db.Where("email = ?", c.email).FirstOrCreate(&user).Error; err != nil {
			return fmt.Errorf("creating %s: %w", c.email, err)
		}
		if err := db.Where(models.UserRole{UserID: user.ID, RoleID: role.ID}).
			FirstOrCreate(&models.UserRole{UserID: user.ID, RoleID: role.ID}).Error; err != nil {
			return fmt.Errorf("giving %s the customer role: %w", c.email, err)
		}

		periodEnd := jsontime.Date{Time: now.AddDate(0, 1, 0)}
		subscription := models.Subscription{
			UserID:           user.ID,
			PlanID:           c.plan.ID,
			Status:           c.status,
			Seats:            c.seats,
			CurrentPeriodEnd: &periodEnd,
		}
		if err := db.Create(&subscription).Error; err != nil {
			return fmt.Errorf("subscribing %s: %w", c.email, err)
		}

		// Three months of history, the oldest paid and the newest still open,
		// so the portal's invoice table has more than one status in it.
		for month := 2; month >= 0; month-- {
			issued := jsontime.Date{Time: now.AddDate(0, -month, 0)}
			invoice := models.Invoice{
				UserID:         user.ID,
				SubscriptionID: subscription.ID,
				Number:         fmt.Sprintf("INV-%s-%03d", now.AddDate(0, -month, 0).Format("200601"), 1),
				Amount:         c.plan.Price,
				Status:         "paid",
				IssuedOn:       &issued,
				PaidOn:         &issued,
			}
			if month == 0 {
				invoice.Status = "open"
				invoice.PaidOn = nil
			}
			if err := db.Create(&invoice).Error; err != nil {
				return fmt.Errorf("invoicing %s: %w", c.email, err)
			}
		}

		for _, m := range c.metrics {
			for day := 0; day < m.days; day++ {
				recorded := jsontime.Date{Time: now.AddDate(0, 0, -day)}
				record := models.UsageRecord{
					UserID:         user.ID,
					SubscriptionID: subscription.ID,
					Metric:         m.name,
					// A little variation per day, so the numbers do not look
					// like a fixture and a chart over them has a shape.
					Quantity:   m.perDay + (day*7)%23,
					RecordedOn: &recorded,
				}
				if err := db.Create(&record).Error; err != nil {
					return fmt.Errorf("recording usage for %s: %w", c.email, err)
				}
			}
		}

		log.Printf("Portal demo: %s / customer123 is on %s", c.email, c.plan.Name)
	}

	return nil
}
