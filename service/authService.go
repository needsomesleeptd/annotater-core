package service

import (
	"errors"

	"github.com/needsomesleeptd/annotater-core/models"
	repository "github.com/needsomesleeptd/annotater-core/repositoryPorts"
	auth_utils "github.com/needsomesleeptd/annotater-core/utilsPorts/authUtils"
	err_wr "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	ErrNoLogin       = models.NewUserErr("login cannot be empty")
	ErrNoPasswd      = models.NewUserErr("password cannot be empty")
	ErrWrongLogin    = models.NewUserErr("wrong login")
	ErrWrongPassword = models.NewUserErr("wrong password")
)

var ERR_LOGIN_STRF = "auth svc - error for user with login %v"

var (
	ErrCreatingUser    = errors.New("error in creating user")
	ErrGeneratingToken = errors.New("error in generating token for user")
	ErrGeneratingHash  = errors.New("error in generating passwdHash for user")
)

const SECRET = "secret"

type IAuthService interface {
	SignIn(candidate *models.User) (tokenStr string, err error)
	SignUp(candidate *models.User) error
}

type AuthService struct {
	logger         *logrus.Logger
	userRepo       repository.IUserRepository
	passwordHasher auth_utils.IPasswordHasher
	tokenizer      auth_utils.ITokenHandler
	key            string
}

func NewAuthService(loggerSrc *logrus.Logger, repo repository.IUserRepository, hasher auth_utils.IPasswordHasher, token auth_utils.ITokenHandler, k string) IAuthService {
	return &AuthService{
		logger:         loggerSrc,
		userRepo:       repo,
		passwordHasher: hasher,
		tokenizer:      token,
		key:            k,
	}
}

func (serv *AuthService) SignUp(candidate *models.User) error {
	var passHash string
	var err error
	if candidate.Login == "" {
		err = err_wr.Wrapf(ErrNoLogin, ERR_LOGIN_STRF, candidate.Login)
		serv.logger.Info(err)
		return err
	}

	if candidate.Password == "" {
		err = err_wr.Wrapf(ErrNoPasswd, ERR_LOGIN_STRF, candidate.Login)
		serv.logger.Info(err)
		return err
	}

	passHash, err = serv.passwordHasher.GenerateHash(candidate.Password)
	if err != nil {
		err = err_wr.Wrapf(err, "error for user with login %v:%v", candidate.Login, ErrGeneratingHash)

		if errors.Is(err, models.ErrDatabaseConnection) {
			serv.logger.Error(err)
		} else {
			serv.logger.Warn(err)
		}

		return err
	}
	candidateHashedPasswd := *candidate
	candidateHashedPasswd.Password = passHash

	err = serv.userRepo.CreateUser(&candidateHashedPasswd)
	if err != nil {
		err = err_wr.Wrapf(err, "error for user with login %v:%v", candidate.Login, ErrCreatingUser)
		if errors.Is(err, models.ErrDatabaseConnection) {
			serv.logger.Error(err)
		} else {
			serv.logger.Warn(err)
		}
		return err
	}
	serv.logger.Infof("auth svc - successfully signed up as user with login %v", candidate.Login)
	return nil
}

func (serv *AuthService) SignIn(candidate *models.User) (string, error) {
	var user *models.User
	var err error
	var tokenStr string
	if candidate.Login == "" {
		err = err_wr.Wrapf(ErrNoLogin, ERR_LOGIN_STRF, candidate.Login)
		serv.logger.Warn(err)
		return "", err
	}

	if candidate.Password == "" {
		err = err_wr.Wrapf(ErrNoPasswd, ERR_LOGIN_STRF, candidate.Login)
		serv.logger.Warn(err)
		return "", err
	}
	user, err = serv.userRepo.GetUserByLogin(candidate.Login)

	if err != nil {
		err = err_wr.Wrapf(err, ERR_LOGIN_STRF+":%v", candidate.Login, ErrWrongLogin)
		if errors.Is(err, models.ErrDatabaseConnection) {
			serv.logger.Error(err)
		} else {
			serv.logger.Warn(err)
		}
		return "", err
	}
	err = serv.passwordHasher.ComparePasswordhash(candidate.Password, user.Password)
	if err != nil {
		err = err_wr.Wrapf(err, ERR_LOGIN_STRF+":%v", candidate.Login, ErrWrongPassword)
		serv.logger.Warn(err)
		return "", err
	}
	tokenStr, err = serv.tokenizer.GenerateToken(*user, serv.key)
	if err != nil {
		err = err_wr.Wrapf(err, ERR_LOGIN_STRF+":%v", candidate.Login, ErrGeneratingToken)
		serv.logger.Warn(err)
		return "", err
	}
	serv.logger.Infof("auth svc - successfully signed in as user with login %v", candidate.Login)
	return tokenStr, nil
}
