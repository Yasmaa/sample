package repository

import (
	"api/internal/domain"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(data *domain.User) (*domain.User, error)
	GetAllUsers() (domain.Users, error)
	GetUserByEmail(email string) (domain.User, error)
	UpdateUser(id int, data *domain.User) (domain.User, error)
	DeleteUser(id int) (bool, error)
	DeleteUsers(ids *[]int) (bool, error)

	CreateResetPasswordToken(data *domain.ResetPassword) (*domain.ResetPassword, error)

	GetResetPasswordToken(token string) (domain.ResetPassword, error)
	UpdateResetPasswordToken(token string) (domain.ResetPassword, error)

	DisableTF(email string) error
}

type userRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewUserRepository(db *gorm.DB, logger *zap.Logger) UserRepository {
	return &userRepository{db: db, logger: logger}
}

func (ur *userRepository) CreateUser(data *domain.User) (*domain.User, error) {

	err := ur.db.Create(&data).Error

	if err != nil {
		return &domain.User{}, err
	}
	return data, nil
}

func (ur *userRepository) GetAllUsers() (domain.Users, error) {

	users := domain.Users{}

	err := ur.db.Find(&users).Error
	if err != nil {

		ur.logger.Info("error while fetching users", zap.String("Error", err.Error()),
			zap.Duration("backoff", time.Second))

		return nil, err
	}

	return users, nil
}

func (ur *userRepository) GetUserByEmail(email string) (domain.User, error) {
	user := domain.User{}
	err := ur.db.Where("Email = ?", email).First(&user).Error

	if err != nil {

		ur.logger.Info("error while fetching user", zap.String("Error", err.Error()),
			zap.Duration("backoff", time.Second))

		return user, err
	}
	return user, nil
}

func (ur *userRepository) UpdateUser(id int, user *domain.User) (domain.User, error) {
	prop := domain.User{}
	err := user.PassHash()
	if err != nil {

		return prop, err
	}
	fmt.Println(user.TwoFa)
	err = ur.db.Model(&domain.User{}).Where("ID = ?", id).First(&prop).Updates(user).Error

	return prop, err
}

func (ur *userRepository) DeleteUser(id int) (bool, error) {

	err := ur.db.Unscoped().Delete(&domain.User{}, id).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (pr *userRepository) DeleteUsers(ids *[]int) (bool, error) {

	err := pr.db.Unscoped().Where("ID IN (?)", *ids).Delete(&domain.Users{}).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (pr *userRepository) CreateResetPasswordToken(data *domain.ResetPassword) (*domain.ResetPassword, error) {

	err := pr.db.Create(&data).Error

	if err != nil {
		return &domain.ResetPassword{}, err
	}
	return data, nil
}

func (pr *userRepository) UpdateResetPasswordToken(token string) (domain.ResetPassword, error) {
	prop := domain.ResetPassword{}
	err := pr.db.Where("token = ?", token).First(&prop).Error
	prop.Expired = true
	pr.db.Save(prop)

	return prop, err
}

func (pr *userRepository) GetResetPasswordToken(token string) (domain.ResetPassword, error) {
	prop := domain.ResetPassword{}
	err := pr.db.Where("token = ?", token).First(&prop).Error
	return prop, err
}

func (ur *userRepository) DisableTF(email string) error {

	err := ur.db.Model(&domain.User{}).Where("email = ?", email).Update("two_fa", false).Error

	return err
}
