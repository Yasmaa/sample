package registry

import (
	"github.com/go-playground/validator/v10"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"

	"api/internal/delivery/http/handler"
	"api/internal/delivery/http/validator"
	"api/internal/repository"
	"api/internal/usecase"

	"go.uber.org/zap"
)

type interactor struct {
	db        *gorm.DB
	validator *validator.Validate
	logger    *zap.Logger
	schedular *asynq.Scheduler
}

type Interactor interface {
	NewAppHandler() handler.AppHandler
}

func NewInteractor(ps *gorm.DB, v *validator.Validate, lg *zap.Logger, a *asynq.Scheduler) Interactor {
	return &interactor{db: ps, validator: v, logger: lg, schedular: a}
}

func (i *interactor) NewAppHandler() handler.AppHandler {
	return handler.AppHandler{
		AuthHandler:     i.NewAuthHandler(),
		UserHandler:     i.NewUserHandler(),}
}

func (i *interactor) NewCustomValidator() validation.CustomValidator {
	return validation.NewCustomValidator(i.validator)
}

// Auth API
func (i *interactor) NewAuthHandler() handler.AuthHandler {
	return handler.NewAuthHandler(i.NewAuthService(), i.NewCustomValidator())
}

func (i *interactor) NewAuthService() usecase.AuthService {
	return usecase.NewAuthService(i.NewAuthRepository())
}

func (i *interactor) NewAuthRepository() repository.AuthRepository {
	return repository.NewAuthRepository(i.db, i.logger)
}


// USER API
func (i *interactor) NewUserHandler() handler.UserHandler {
	return handler.NewUserHandler(i.NewUserService(), i.NewCustomValidator())
}

func (i *interactor) NewUserService() usecase.UserService {
	return usecase.NewUserService(i.NewUserRepository())
}

func (i *interactor) NewUserRepository() repository.UserRepository {
	return repository.NewUserRepository(i.db, i.logger)
}
