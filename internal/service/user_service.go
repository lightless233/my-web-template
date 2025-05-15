package service

import (
	"crypto/md5"
	"encoding/hex"

	"go.uber.org/zap"
	"my-web-template/internal/entity/vo"
	"my-web-template/internal/repository"
	"my-web-template/internal/result"
)

type UserServiceInterface interface {
	SaveUser(username, email, password string) (*vo.UserVO, result.AppError)
	GetUserByUsername(username string) (*vo.UserVO, result.AppError)
}

type UserService struct {
	userRepository *repository.UserRepository
	logger         *zap.SugaredLogger
}

func NewUserService(userRepository *repository.UserRepository, logger *zap.SugaredLogger) *UserService {
	return &UserService{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (u *UserService) SaveUser(username, email, password string) (*vo.UserVO, result.AppError) {
	sum := md5.Sum([]byte(password))
	newPassword := hex.EncodeToString(sum[:])
	user, err := u.userRepository.SaveUser(username, email, newPassword)
	if err != nil {
		return nil, err
	}

	return user.ToVO(), nil
}

func (u *UserService) GetUserByUsername(username string) (*vo.UserVO, result.AppError) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	return user.ToVO(), nil
}

var _ UserServiceInterface = (*UserService)(nil)
