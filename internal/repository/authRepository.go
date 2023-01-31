package repository

import (
	"api/internal/domain"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthRepository interface {
	CreateUser(data *domain.User) (*domain.User, error)
	GetUserByEmail(email string) (domain.User, error)
	Confirm(token string) (domain.User, error)
	UpdateUser(id int, data *domain.User) (domain.User, error)
	CreateResetPasswordToken(data *domain.ResetPassword) (*domain.ResetPassword, error)
	GetResetPasswordToken(token string) (domain.ResetPassword, error)
	UpdateResetPasswordToken(token string) (domain.ResetPassword, error)
}

type authRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewAuthRepository(db *gorm.DB, logger *zap.Logger) AuthRepository {
	return &authRepository{db: db, logger: logger}
}

func (ur *authRepository) CreateUser(data *domain.User) (*domain.User, error) {

	err := ur.db.Create(&data).Error

	if err != nil {
		return &domain.User{}, err
	}
	return data, nil
}

func (ur *authRepository) GetUserByEmail(email string) (domain.User, error) {
	user := domain.User{}
	err := ur.db.Where("Email = ?", email).First(&user).Error

	if err != nil {

		ur.logger.Info("error while fetching user", zap.String("Error", err.Error()),
			zap.Duration("backoff", time.Second))

		return user, err
	}
	return user, nil
}

func (ur *authRepository) Confirm(token string) (domain.User, error) {
	user := domain.User{}
	err := ur.db.Where("token = ?", token).First(&user).Error
	if err == nil {
		user.IsVerified = true
		ur.db.Save(user)
	}
	return user, err
}

func (ur *authRepository) UpdateUser(id int, user *domain.User) (domain.User, error) {
	prop := domain.User{}
	user.PassHash()
	err := ur.db.Model(&domain.User{}).Where("ID = ?", id).First(&prop).Updates(user).Error

	return prop, err
}

func (pr *authRepository) CreateResetPasswordToken(data *domain.ResetPassword) (*domain.ResetPassword, error) {

	err := pr.db.Create(&data).Error

	if err != nil {
		return &domain.ResetPassword{}, err
	}
	return data, nil

}

func (pr *authRepository) UpdateResetPasswordToken(token string) (domain.ResetPassword, error) {
	prop := domain.ResetPassword{}
	err := pr.db.Where("token = ?", token).First(&prop).Error
	prop.Expired = true
	pr.db.Save(prop)

	return prop, err
}

func (pr *authRepository) GetResetPasswordToken(token string) (domain.ResetPassword, error) {
	prop := domain.ResetPassword{}
	err := pr.db.Where("token = ?", token).First(&prop).Error
	return prop, err
}
