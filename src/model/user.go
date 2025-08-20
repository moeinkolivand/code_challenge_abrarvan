package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	APIKey       string    `gorm:"uniqueIndex;size:64;not null"`
	Name         string    `gorm:"size:255;not null"`
	Email        string    `gorm:"size:255"`
	Balance      float64   `gorm:"type:decimal(12,4);default:0.0000"`
	RatePerSMS   float64   `gorm:"type:decimal(6,4);default:0.0500"`
	DailyLimit   int       `gorm:"default:10000"`
	MonthlyLimit int       `gorm:"default:300000"`
	IsActive     bool      `gorm:"default:true"`
	Messages     []Message `gorm:"foreignKey:UserID"`
}

type Message struct {
	gorm.Model
	UserID        uint      `gorm:"not null;index"`
	PhoneNumber   string    `gorm:"size:20;not null"`
	MessageText   string    `gorm:"type:text;not null"`
	MessageType   string    `gorm:"size:20;default:'normal';check:message_type IN ('normal','express')"`
	Status        string    `gorm:"size:20;default:'pending';check:status IN ('pending','queued','sent','delivered','failed','cancelled')"`
	Cost          float64   `gorm:"type:decimal(6,4);not null"`
	ProviderID    uint      `gorm:"index"`
	ExternalMsgID string    `gorm:"size:255"`
	Priority      int       `gorm:"default:5;check:priority BETWEEN 1 AND 10"`
	Attempts      int       `gorm:"default:0"`
	MaxAttempts   int       `gorm:"default:3"`
	CreatedAt     time.Time `gorm:"index"`
	QueuedAt      time.Time
	SentAt        time.Time
	DeliveredAt   time.Time
	FailedAt      time.Time
	ErrorMessage  string   `gorm:"type:text"`
	User          User     `gorm:"foreignKey:UserID"`
	Provider      Provider `gorm:"foreignKey:ProviderID"`
}
