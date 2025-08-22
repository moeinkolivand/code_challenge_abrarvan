package repository

import (
	"abrarvan_challenge/config"
	"abrarvan_challenge/logging"
	"abrarvan_challenge/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var logger logging.Logger = logging.NewLogger(config.GetConfig())

type IUserRepository interface {
	Create(user *model.User) error
	FindByID(id uint) (*model.User, error)
	FindByPhoneNumber(phoneNumber string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(user *model.User) error {
	logger.Info(logging.Postgres, logging.Insert, "Creating user", map[logging.ExtraKey]interface{}{
		"Email":       user.Email,
		"PhoneNumber": user.PhoneNumber,
	})
	if err := r.db.Create(user).Error; err != nil {
		logger.Error(logging.Postgres, logging.Insert, fmt.Sprintf("Failed to create user: %v", err), nil)
		return fmt.Errorf("failed to create user: %w", err)
	}
	logger.Info(logging.Postgres, logging.Insert, "User created successfully", map[logging.ExtraKey]interface{}{
		"user_id": user.ID,
	})
	return nil
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	logger.Info(logging.Postgres, logging.Select, "Finding user by ID", map[logging.ExtraKey]interface{}{
		"user_id": id,
	})
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn(logging.Postgres, logging.Select, "User not found", map[logging.ExtraKey]interface{}{
				"user_id": id,
			})
			return nil, fmt.Errorf("user not found: id %d", id)
		}
		logger.Error(logging.Postgres, logging.Select, fmt.Sprintf("Failed to find user: %v", err), nil)
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return &user, nil
}

func (r *userRepository) FindByPhoneNumber(phoneNumber string) (*model.User, error) {
	logger.Info(logging.Postgres, logging.Select, "Finding user by username", map[logging.ExtraKey]interface{}{
		"phoneNumber": phoneNumber,
	})
	var user model.User
	if err := r.db.Where("phone_number = ?", phoneNumber).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn(logging.Postgres, logging.Select, "User not found", map[logging.ExtraKey]interface{}{
				"phoneNumber": phoneNumber,
			})
			return nil, fmt.Errorf("user not found: phoneNumber %s", phoneNumber)
		}
		logger.Error(logging.Postgres, logging.Select, fmt.Sprintf("Failed to find user: %v", err), nil)
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return &user, nil
}

func (r *userRepository) Update(user *model.User) error {
	logger.Info(logging.Postgres, logging.Update, "Updating user", map[logging.ExtraKey]interface{}{
		"user_id":      user.ID,
		"username":     user.Email,
		"phone_number": user.PhoneNumber,
	})
	if err := r.db.Save(user).Error; err != nil {
		logger.Error(logging.Postgres, logging.Update, fmt.Sprintf("Failed to update user: %v", err), nil)
		return fmt.Errorf("failed to update user: %w", err)
	}
	logger.Info(logging.Postgres, logging.Update, "User updated successfully", map[logging.ExtraKey]interface{}{
		"user_id": user.ID,
	})
	return nil
}

func (r *userRepository) Delete(id uint) error {
	logger.Info(logging.Postgres, logging.Delete, "Deleting user", map[logging.ExtraKey]interface{}{
		"user_id": id,
	})
	if err := r.db.Delete(&model.User{}, id).Error; err != nil {
		logger.Error(logging.Postgres, logging.Delete, fmt.Sprintf("Failed to delete user: %v", err), nil)
		return fmt.Errorf("failed to delete user: %w", err)
	}
	logger.Info(logging.Postgres, logging.Delete, "User deleted successfully", map[logging.ExtraKey]interface{}{
		"user_id": id,
	})
	return nil
}
