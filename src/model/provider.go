package model

import (
	"gorm.io/gorm"
	"time"
)

type Provider struct {
	gorm.Model
	Name              string    `gorm:"size:100;not null"`
	APIUrl            string    `gorm:"size:500;not null"`
	APIKey            string    `gorm:"size:255"`
	APISecret         string    `gorm:"size:255"`
	CostPerSMS        float64   `gorm:"type:decimal(6,4);default:0.0300"`
	Priority          int       `gorm:"default:5"`
	IsActive          bool      `gorm:"default:true"`
	SuccessRate       float64   `gorm:"type:decimal(5,2);default:95.00"`
	AvgDeliveryTime   int       `gorm:"default:30"`
	DailyLimit        int       `gorm:"default:1000000"`
	CurrentDailyUsage int       `gorm:"default:0"`
	LastResetDate     time.Time `gorm:"type:date;default:CURRENT_DATE"`
}
