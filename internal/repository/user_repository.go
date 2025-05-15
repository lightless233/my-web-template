package repository

import (
	"go.uber.org/zap"
	"my-web-template/internal/constant"
	"my-web-template/internal/model"
	"my-web-template/internal/result"
	"xorm.io/xorm"
)

type UserRepositoryInterface interface {
	SaveUser(username, email, password string) (*model.AppUserModel, result.AppError)
	GetUserByUsername(username string) (*model.AppUserModel, result.AppError)
}

type UserRepository struct {
	db     *xorm.Engine
	logger *zap.SugaredLogger
}

func NewUserRepository(db *xorm.Engine, logger *zap.SugaredLogger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (u *UserRepository) SaveUser(username, email, password string) (*model.AppUserModel, result.AppError) {
	example := &model.AppUserModel{
		Username: username,
		Email:    email,
		Password: password,
		State:    constant.UserStatusActive,
	}
	_, err := u.db.Insert(example)
	if err != nil {
		return nil, result.NewAppErrorFromError(constant.CodeDBError, err, true)
	}

	return example, nil
}

func (u *UserRepository) GetUserByUsername(username string) (*model.AppUserModel, result.AppError) {
	user := &model.AppUserModel{}

	exists, err := u.db.Where("username = ? AND deleted = false", username).Get(user)
	if err != nil {
		return nil, result.NewAppErrorFromError(constant.CodeDBError, err, true)
	}
	if !exists {
		return nil, nil
	}

	return user, nil
}

// 确保接口正确实现，如果 UserRepository 没有实现 UserRepositoryInterface，那么这里会报错
var _ UserRepositoryInterface = (*UserRepository)(nil)
