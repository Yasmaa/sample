package usecase

import (
	"api/config"
	"api/internal/delivery/http/handler/requests"
	"api/internal/domain"
	"api/internal/repository"
	"api/pkg/auth"
	"api/pkg/middlewares"
	"api/pkg/utils"
	"encoding/base32"
	"errors"
	"fmt"

	"github.com/dgryski/dgoogauth"
)

type UserService interface {
	CreateUser(data *domain.User) (*domain.User, error)
	GetAllUsers() (domain.Users, error)
	GetUserByEmail(email string) (domain.User, error)
	UpdateUser(id int, data *domain.User) (domain.User, error)
	DeleteUser(id int) (bool, error)
	DeleteUsers(ids *[]int) (bool, error)

	UpdateCurrentUser(token string, data *domain.User) (*domain.User, error)
	UpdateCurrentUserPassword(token string, data *requests.UpdateUserPasswordRequest) (*domain.User, error)
	GetCurrentUser(token string) (*domain.User, error)

	GetQr(token string) (string, string, error)
	VerifyTF(token string, code string) (bool, error)
	DisableTF(token string) (error)

}

type userService struct {
	UserRepository repository.UserRepository
}

func NewUserService(ur repository.UserRepository) UserService {
	return &userService{UserRepository: ur}
}

func (us *userService) CreateUser(data *domain.User) (*domain.User, error) {
	return us.UserRepository.CreateUser(data)

}

func (us *userService) GetAllUsers() (domain.Users, error) {
	return us.UserRepository.GetAllUsers()

}

func (us *userService) GetUserByEmail(email string) (domain.User, error) {
	return us.UserRepository.GetUserByEmail(email)
}

func (us *userService) UpdateUser(id int, data *domain.User) (domain.User, error) {

	return us.UserRepository.UpdateUser(id, data)

}

func (us *userService) DeleteUser(id int) (bool, error) {

	return us.UserRepository.DeleteUser(id)

}

func (us *userService) DeleteUsers(ids *[]int) (bool, error) {

	return us.UserRepository.DeleteUsers(ids)

}

func (us *userService) UpdateCurrentUser(token string, data *domain.User) (*domain.User, error) {

	tokenMaker, _ := middlewares.NewPasetoMaker(config.C.Token)
	payload, err := tokenMaker.VerifyToken(token)
	if err != nil {

		return &domain.User{}, err
	}

	r, err := us.UserRepository.UpdateUser(int(payload.ID), data)

	if err != nil {
		return &r, err
	} else {
		r.Password = ""
		return &r, nil
	}

}

func (us *userService) UpdateCurrentUserPassword(token string, data *requests.UpdateUserPasswordRequest) (*domain.User, error) {

	tokenMaker, _ := middlewares.NewPasetoMaker(config.C.Token)

	payload, err := tokenMaker.VerifyToken(token)
	if err != nil {
		return &domain.User{}, err
	}

	u, err := us.UserRepository.GetUserByEmail(payload.Email)
	if err != nil {

		return &u, err
	}

	err = auth.ComparePassword(u.Password, data.OldPassword)

	if err != nil {

		return &u, err
	}

	if data.OldPassword == data.NewPassword {

		return &u, errors.New("password must be new")
	}

	if data.ConfirmPassword != data.NewPassword {

		return &u, errors.New("passwords do not match")
	}

	r, err := us.UserRepository.UpdateUser(int(payload.ID), &domain.User{Password: data.NewPassword})
	if err != nil {
		return &r, err
	} else {
		r.Password = ""
		return &r, nil
	}

}

func (us *userService) GetCurrentUser(token string) (*domain.User, error) {

	tokenMaker, _ := middlewares.NewPasetoMaker(config.C.Token)
	payload, err := tokenMaker.VerifyToken(token)
	if err != nil {
		return &domain.User{}, err
	}

	u, err := us.UserRepository.GetUserByEmail(payload.Email)
	if err != nil {
		return &u, err
	} else {
		u.Password = ""
		return &u, nil
	}

}

func (us *userService) GetQr(token string) (string, string, error) {

	tokenMaker, _ := middlewares.NewPasetoMaker(config.C.Token)
	payload, err := tokenMaker.VerifyToken(token)
	if err != nil {
		return "", "", err
	}

	user, err := us.UserRepository.GetUserByEmail(payload.Email)
	if err != nil {
		return "", "", err
	}

	userSecret := user.Secret
	if userSecret == "" {

		randomStr := utils.RandStr(6, "alphanum")
		secret := base32.StdEncoding.EncodeToString([]byte(randomStr))
		userSecret = secret
	}

	_, err = us.UserRepository.UpdateUser(int(user.ID), &domain.User{Secret: userSecret})

	if err != nil {
		return "", "", err
	}

	authLink := fmt.Sprintf("otpauth://totp/GPM:%s?secret=%s&issuer=GPM", user.Username, userSecret)

	return userSecret, authLink, nil
}

func (us *userService) VerifyTF(token string, code string) (bool, error) {

	tokenMaker, _ := middlewares.NewPasetoMaker(config.C.Token)
	payload, err := tokenMaker.VerifyToken(token)
	if err != nil {
		return false, err
	}

	user, err := us.UserRepository.GetUserByEmail(payload.Email)
	if err != nil {
		return false, err
	}

	userSecret := user.Secret
	if userSecret == "" {

		randomStr := utils.RandStr(6, "alphanum")
		secret := base32.StdEncoding.EncodeToString([]byte(randomStr))
		userSecret = secret
	}

	_, err = us.UserRepository.UpdateUser(int(user.ID), &domain.User{Secret: userSecret, TwoFa: true})

	if err != nil {
		return false, err
	}

	otpc := &dgoogauth.OTPConfig{
		Secret:      userSecret,
		WindowSize:  3,
		HotpCounter: 0,
	}

	isVerified, err := otpc.Authenticate(code)

	return isVerified, nil
}




func (us *userService) DisableTF(token string) (error) {

	tokenMaker, _ := middlewares.NewPasetoMaker(config.C.Token)
	payload, err := tokenMaker.VerifyToken(token)
	if err != nil {
		return err
	}

	
	return us.UserRepository.DisableTF(payload.Email)

	

}

