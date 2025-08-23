package service

import (
	"abrarvan_challenge/config"
	"abrarvan_challenge/infrastructure/persistance/broker"
	"abrarvan_challenge/infrastructure/persistance/database"
	"abrarvan_challenge/infrastructure/persistance/repository"
	"abrarvan_challenge/logging"
	"abrarvan_challenge/model"
	"abrarvan_challenge/provider"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
)

var logger logging.Logger = logging.NewLogger(config.GetConfig())

type UserService struct {
	provider provider.IProvider
	repo     repository.IUserRepository
}

func NewUserService(provider provider.IProvider, repo repository.IUserRepository) *UserService {
	return &UserService{
		provider: provider,
		repo:     repo,
	}
}

func (userService *UserService) GetUserByPhoneNumber(phoneNumber string) (*model.User, error) {
	usr, err := userService.repo.FindByPhoneNumber(phoneNumber)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (userService *UserService) SendMessage(phoneNumber string, message string) error {
	usr, err := userService.repo.FindByPhoneNumber(phoneNumber)
	if err != nil {
		logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return err
	}
	usrBalance, hasEnoughBalance := userService.CheckBalance(usr)
	if !hasEnoughBalance {
		logger.Error(logging.Validation, logging.Request, "not enough balance"+phoneNumber, map[logging.ExtraKey]interface{}{
			"phoneNumber": phoneNumber,
			"usrBalance":  usrBalance,
		})
		return fmt.Errorf("not enough balance")
	}
	messageBody, err := json.Marshal(map[string]string{
		"phoneNumber": phoneNumber,
		"message":     message,
	})
	if err != nil {
		logger.Error(logging.RabbitMQ, logging.Publish, "Failed to marshal message: "+err.Error(), nil)
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	if usr.UserType == model.UserNormal {
		err = broker.Publish("consumerChannel", "", "my_queue", messageBody, 1)
		if err != nil {
			logger.Error(logging.RabbitMQ, logging.Publish, "Failed to publish message: "+err.Error(), nil)
			return err
		}
	} else {
		err = broker.Publish("consumerChannel", "", "my_queue", messageBody, 5)
		if err != nil {
			logger.Error(logging.RabbitMQ, logging.Publish, "Failed to publish message: "+err.Error(), nil)
			return err
		}
	}

	return nil
}

func (userService *UserService) CheckBalance(user *model.User) (float64, bool) {
	smsCost := userService.providerSmsConst()
	if user.Balance == 0 || user.Balance-float64(smsCost) <= 0 {
		return user.Balance, false
	}
	return user.Balance, true
}

func (userService *UserService) providerSmsConst() uint {
	return 10
}

func (userService *UserService) ConsumerSendSms(phoneNumber, message string) error {
	db := database.GetDb()

	err := db.Transaction(func(tx *gorm.DB) error {
		var user model.User
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("phone_number = ?", phoneNumber).First(&user).Error; err != nil {
			logger.Error(logging.Postgres, logging.Select, fmt.Sprintf("Failed to lock user with phone %s: %v", phoneNumber, err), nil)
			return fmt.Errorf("failed to lock user with phone %s: %w", phoneNumber, err)
		}

		currentUserBalance, hasEnoughBalance := userService.CheckBalance(&user)
		smsCost := userService.providerSmsConst()
		if !hasEnoughBalance {
			logger.Error(logging.Validation, logging.Request, fmt.Sprintf("Insufficient balance for user %s: have %.2f, need %.2f", phoneNumber, user.Balance, smsCost), map[logging.ExtraKey]interface{}{
				"phoneNumber": phoneNumber,
				"balance":     user.Balance,
				"smsCost":     smsCost,
			})
			return fmt.Errorf("insufficient balance: have %.2f, need %.2f", currentUserBalance, smsCost)
		}

		user.Balance -= float64(smsCost)
		if err := tx.Save(&user).Error; err != nil {
			logger.Error(logging.Postgres, logging.Update, fmt.Sprintf("Failed to update balance for user %s: %v", phoneNumber, err), nil)
			return fmt.Errorf("failed to update balance: %w", err)
		}

		if err := userService.provider.SendSMS(phoneNumber, message); err != nil {
			logger.Error(logging.SMSProvider, logging.Request, fmt.Sprintf("Failed to send SMS to %s: %v", phoneNumber, err), nil)
			return fmt.Errorf("failed to send SMS: %w", err)
		}

		logger.Info(logging.SMSProvider, logging.Request, fmt.Sprintf("Successfully sent SMS to %s and deducted balance", phoneNumber), map[logging.ExtraKey]interface{}{
			"phoneNumber": phoneNumber,
			"newBalance":  user.Balance,
		})
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
