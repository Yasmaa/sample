package usecase

import (
	"api/config"
	"api/internal/delivery/http/handler/requests"
	"api/internal/domain"
	"api/internal/repository"
	"api/pkg/auth"
	"api/pkg/middlewares"
	"errors"
	"fmt"
	"time"

	"github.com/dgryski/dgoogauth"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService interface {
	Confirm(token string) bool
	Register(data *requests.RegisterRequest) (*domain.User, error)
	Login(data *requests.LoginRequest) (string, bool, error)
	ResetLink(data *requests.ResetLinkRequest) error
	ResetPassword(data *requests.ResetPasswordRequest) error
	VerifyToken(token string) error

	VerifyTF(data *requests.TwoFaToken) (string, bool, error)
}

type authService struct {
	AuthRepository repository.AuthRepository
}

func NewAuthService(ur repository.AuthRepository) AuthService {
	return &authService{AuthRepository: ur}
}

func (us *authService) Confirm(token string) bool {
	_, err := us.AuthRepository.Confirm(token)
	return err == nil
}

func (us *authService) Register(data *requests.RegisterRequest) (*domain.User, error) {

	user := &domain.User{
		Firstname: data.Firstname,
		Lastname:  data.Lastname,
		Username:  data.Username,
		Email:     data.Email,
		Password:  data.Password,
		Token: 	uuid.NewString(),
	}


	u, err := us.AuthRepository.CreateUser(user)
	if err != nil {
		return u, err
	}
	url := fmt.Sprintf("http://%s:%s/auth/confirm/%s", config.C.Link.HOST, config.C.Link.PORT, user.Token)
	go auth.SendMail("Account Confirmation", data.Email, url, "confirm.html")
	return u, err

}

func (us *authService) Login(data *requests.LoginRequest) (string, bool, error) {

	u, err := us.AuthRepository.GetUserByEmail(data.Email)
	if err != nil {
		return "", false, err
	}

	if !u.IsVerified {
		return "", false, errors.New("please confirm your account before login")
	}

	token, err := u.LoginCheck(data.Password)
	fmt.Println("token", token, err)
	if err != nil {
		return "", false, errors.New("wrong email or password")
	}

	return token, u.TwoFa, nil

}

func (us *authService) ResetLink(data *requests.ResetLinkRequest) error {

	result, err := us.AuthRepository.GetUserByEmail(data.Email)
	if err != nil {
		return err
	}

	tokenMaker, err := middlewares.NewPasetoMaker(config.C.Token)
	if err != nil {
		return err
	}
	resetToken, err := tokenMaker.CreateToken(result.ID, result.Email, 24*3*time.Hour)
	if err != nil {
		return err
	}

	us.AuthRepository.CreateResetPasswordToken(&domain.ResetPassword{Token: resetToken, UserID: result.ID, Expired: false})
	url := fmt.Sprintf("http://%s:%s/auth/reset-password/%s", config.C.Link.HOST, config.C.Link.PORT, resetToken)
	go auth.SendMail("Reset Password", result.Email, url, "reset.html")

	return nil
}

func (us *authService) ResetPassword(data *requests.ResetPasswordRequest) error {

	tokenMaker, _ := middlewares.NewPasetoMaker(config.C.Token)
	payload, err := tokenMaker.VerifyToken(data.Token)
	if err != nil {
		return err
	}

	rt, err := us.AuthRepository.GetResetPasswordToken(data.Token)
	if err != nil {
		return err
	}

	if rt.Expired {
		return errors.New("expired token")
	}
	_, err = us.AuthRepository.UpdateUser(int(payload.ID), &domain.User{ Model: gorm.Model{ID: payload.ID} ,Password: data.Password})
	if err != nil {
		return err
	}

	_, err = us.AuthRepository.UpdateResetPasswordToken(data.Token)
	if err != nil {
		return err
	}
	return nil
}

func (us *authService) VerifyToken(token string) error {

	tokenMaker, _ := middlewares.NewPasetoMaker(config.C.Token)
	_, err := tokenMaker.VerifyToken(token)
	if err != nil {
		return err
	}
	return nil

}

func (us *authService) VerifyTF(data *requests.TwoFaToken) (string, bool, error) {

	user, err := us.AuthRepository.GetUserByEmail(data.Email)
	if err != nil {
		return "", false, err
	}

	userSecret := user.Secret
	if userSecret == "" {

		return "", false, err

	}

	otpc := &dgoogauth.OTPConfig{
		Secret:      userSecret,
		WindowSize:  3,
		HotpCounter: 0,
	}

	isVerified, err := otpc.Authenticate(data.Token)

	tokenMaker, _ := middlewares.NewPasetoMaker(config.C.Token)
	token, err := tokenMaker.CreateToken(user.ID, user.Email, 24*3*time.Hour)
	if err != nil {
		return "", false, err
	}

	return token, isVerified, nil
}
