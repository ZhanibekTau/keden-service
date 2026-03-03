package app

import (
	"keden-service/back/internal/configs/structures"
	"keden-service/back/internal/models"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Company{},
		&models.Subscription{},
		&models.Document{},
		&models.RefreshToken{},
	)
	if err != nil {
		return err
	}

	// Fix: drop stale company_id column from refresh_tokens if it exists.
	// It was added by a previous AutoMigrate when the model had a CompanyID field,
	// then the field was removed but AutoMigrate never drops columns.
	if err := db.Exec("ALTER TABLE refresh_tokens DROP COLUMN IF EXISTS company_id").Error; err != nil {
		return err
	}

	logrus.Info("Database migrations completed successfully")
	return nil
}

func SeedRoles(db *gorm.DB) error {
	roles := []string{"admin", "client"}
	for _, name := range roles {
		var existing models.Role
		result := db.Where("name = ?", name).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&models.Role{Name: name}).Error; err != nil {
				return err
			}
			logrus.Infof("Role '%s' created", name)
		}
	}
	return nil
}

func SeedAdmin(db *gorm.DB, adminCfg *structures.AdminConfig) error {
	if err := SeedRoles(db); err != nil {
		return err
	}

	if adminCfg.AdminEmail == "" || adminCfg.AdminPassword == "" {
		logrus.Warn("Admin credentials not provided, skipping admin seed")
		return nil
	}

	var existing models.User
	result := db.Where("email = ?", adminCfg.AdminEmail).First(&existing)
	if result.Error == nil {
		logrus.Info("Admin account already exists")
		return nil
	}

	var adminRole models.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(adminCfg.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := models.User{
		Email:        adminCfg.AdminEmail,
		PasswordHash: string(hash),
		FirstName:    "Admin",
		LastName:     "System",
		Phone:        "+70000000000",
		RoleID:       adminRole.ID,
		AccountType:  "individual",
		IsActive:     true,
	}

	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	logrus.Info("Admin account created successfully")
	return nil
}
