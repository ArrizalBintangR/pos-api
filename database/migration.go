package database

import (
	"log"

	"interview-user/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&models.User{},
		&models.SaleOrder{},
		&models.SaleOrderItem{},
	)

	if err != nil {
		return err
	}

	log.Println("Database migrations completed")
	return nil
}

func Seed(db *gorm.DB) error {
	log.Println("Seeding database...")

	// Check if owner already exists
	var count int64
	db.Model(&models.User{}).Where("role = ?", models.RoleOwner).Count(&count)
	if count > 0 {
		log.Println("Owner user already exists, skipping seed")
		return nil
	}

	// Create default owner
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("owner123"), bcrypt.DefaultCost)
	owner := models.User{
		Username: "owner",
		Password: string(hashedPassword),
		Name:     "System Owner",
		Role:     models.RoleOwner,
		IsActive: true,
	}

	if err := db.Create(&owner).Error; err != nil {
		return err
	}

	// Create default cashier
	hashedPassword, _ = bcrypt.GenerateFromPassword([]byte("cashier123"), bcrypt.DefaultCost)
	cashier := models.User{
		Username: "cashier",
		Password: string(hashedPassword),
		Name:     "Default Cashier",
		Role:     models.RoleCashier,
		IsActive: true,
	}

	if err := db.Create(&cashier).Error; err != nil {
		return err
	}

	log.Println("Database seeding completed")
	log.Println("Default users created:")
	log.Println("  Owner: username=owner, password=owner123")
	log.Println("  Cashier: username=cashier, password=cashier123")

	return nil
}
