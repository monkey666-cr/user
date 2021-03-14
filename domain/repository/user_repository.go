package repository

import (
	"gorm.io/gorm"
	"user/domain/model"
)

type IUserRepository interface {
	InitTable() error

	FindUserByName(string) (*model.User, error)

	FindUserByID(int64) (*model.User, error)

	CreateUser(*model.User) (int64, error)

	DeleteUserByID(int64) error

	UpdateUser(*model.User) error

	FindAll() ([]model.User, error)
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{mysqlDb: db}
}

type UserRepository struct {
	mysqlDb *gorm.DB
}

func (u *UserRepository) InitTable() error {
	return u.mysqlDb.AutoMigrate(&model.User{})
}

func (u *UserRepository) FindUserByName(name string) (user *model.User, err error) {
	user = &model.User{}
	return user, u.mysqlDb.Where("user_name = ?", name).First(user).Error
}

func (u *UserRepository) FindUserByID(userID int64) (user *model.User, err error) {
	user = &model.User{}
	return user, u.mysqlDb.First(user, userID).Error
}

func (u *UserRepository) CreateUser(user *model.User) (userID int64, err error) {
	return user.ID, u.mysqlDb.Create(user).Error
}

func (u *UserRepository) DeleteUserByID(userID int64) (err error) {
	return u.mysqlDb.Where("id = ?", userID).Delete(&model.User{}).Error
}

func (u *UserRepository) UpdateUser(user *model.User) error {
	return u.mysqlDb.Model(user).Updates(&user).Error
}

func (u *UserRepository) FindAll() (userAll []model.User, err error) {
	return userAll, u.mysqlDb.Find(&userAll).Error
}
