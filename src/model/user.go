package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type UserTypeLevel string
type MessagePiority string

const (
	UserNormal  UserTypeLevel = "normal"
	UserExpress UserTypeLevel = "express"
)

const (
	MessageNormal  MessagePiority = "normal"
	MessageExpress MessagePiority = "express"
)

type User struct {
	gorm.Model
	APIKey       string        `gorm:"uniqueIndex;size:64;not null"`
	PhoneNumber  string        `gorm:"uniqueIndex;size:11;not null"`
	UserType     UserTypeLevel `gorm:"size:10;default:'normal';check:user_type IN ('normal','express')"`
	Name         string        `gorm:"size:255;not null"`
	Email        string        `gorm:"size:255"`
	Balance      float64       `gorm:"type:decimal(12,4);default:0.0000"`
	RatePerSMS   float64       `gorm:"type:decimal(6,4);default:0.0500"`
	DailyLimit   int           `gorm:"default:10000"`
	MonthlyLimit int           `gorm:"default:300000"`
	IsActive     bool          `gorm:"default:true"`
	Messages     []Message     `gorm:"foreignKey:UserID"`
}

type Message struct {
	gorm.Model
	UserID        uint           `gorm:"not null;index"`
	PhoneNumber   string         `gorm:"size:20;not null"`
	MessageText   string         `gorm:"type:text;not null"`
	MessageType   MessagePiority `gorm:"size:20;default:'normal';check:message_type IN ('normal','express')"`
	Status        string         `gorm:"size:20;default:'pending';check:status IN ('pending','queued','sent','delivered','failed','cancelled')"`
	Cost          float64        `gorm:"type:decimal(6,4);not null"`
	ProviderID    uint           `gorm:"index"`
	ExternalMsgID string         `gorm:"size:255"`
	Priority      int            `gorm:"default:5;check:priority BETWEEN 1 AND 10"`
	Attempts      int            `gorm:"default:0"`
	MaxAttempts   int            `gorm:"default:3"`
	CreatedAt     time.Time      `gorm:"index"`
	QueuedAt      time.Time
	SentAt        time.Time
	DeliveredAt   time.Time
	FailedAt      time.Time
	ErrorMessage  string   `gorm:"type:text"`
	User          User     `gorm:"foreignKey:UserID"`
	Provider      Provider `gorm:"foreignKey:ProviderID"`
}

func SeedUsers(db *gorm.DB) error {
	users := []User{
		{
			APIKey:       "api_key_1234567890abcdef",
			Name:         "John Doe",
			Email:        "john.doe@example.com",
			Balance:      100.0000,
			RatePerSMS:   0.0500,
			DailyLimit:   10000,
			PhoneNumber:  "09999948734",
			UserType:     UserNormal,
			MonthlyLimit: 300000,
			IsActive:     true,
			Messages:     []Message{},
			Model:        gorm.Model{CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
		{
			APIKey:       "api_key_abcdef1234567890",
			Name:         "Jane Smith",
			Email:        "jane.smith@example.com",
			Balance:      50.0000,
			RatePerSMS:   0.0450,
			DailyLimit:   5000,
			PhoneNumber:  "09332823692",
			UserType:     UserExpress,
			MonthlyLimit: 150000,
			IsActive:     true,
			Messages:     []Message{},
			Model:        gorm.Model{CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
	}
	for _, user := range users {
		var fakeUser User
		result := db.Where("phone_number = ?", user.PhoneNumber).First(&fakeUser).Error
		if errors.Is(result, gorm.ErrRecordNotFound) {
			if err := db.Create(&user).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
