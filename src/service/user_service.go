package service

import (
	"abrarvan_challenge/config"
	"abrarvan_challenge/infrastructure/persistance/broker"
	"abrarvan_challenge/infrastructure/persistance/repository"
	"abrarvan_challenge/logging"
	"abrarvan_challenge/model"
	"abrarvan_challenge/provider"
	"encoding/json"
	"fmt"
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
	err = broker.Publish("consumerChannel", "", "my_queue", messageBody)
	if err != nil {
		logger.Error(logging.RabbitMQ, logging.Publish, "Failed to publish message: "+err.Error(), nil)
		return err
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
