package model

import (
	"fmt"
	"gorm.io/gorm"
)

func MigrateDatabaseTables(db *gorm.DB) error {
	models := []interface{}{&User{}, &Message{}, &BillingHistory{}, &Provider{}}

	if err := autoMigrateModels(db, models...); err != nil {
		return err
	}
	return nil
}

func autoMigrateModels(db *gorm.DB, models ...interface{}) error {
	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to auto-migrate models: %v", err)
	}
	return nil
}
