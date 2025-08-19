package model

import (
	"gorm.io/gorm"
)

type BillingHistory struct {
	gorm.Model
	UserID          uint    `gorm:"not null;index"`
	TransactionType string  `gorm:"size:20;not null;check:transaction_type IN ('credit','debit','refund','bonus')"`
	Amount          float64 `gorm:"type:decimal(12,4);not null"`
	BalanceBefore   float64 `gorm:"type:decimal(12,4);not null"`
	BalanceAfter    float64 `gorm:"type:decimal(12,4);not null"`
	Description     string  `gorm:"type:text"`
	PaymentMethod   string  `gorm:"size:50"`
	ReferenceID     string  `gorm:"size:255"`
	MessageID       uint
	User            User `gorm:"foreignKey:UserID"`
}
