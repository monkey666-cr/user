package service

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"user/domain/model"
	"user/domain/repository"
)

type IUserDataService interface {
	AddUser(user *model.User) (int64, error)
	DeleteUser(int64) error
	UpdateUser(*model.User, bool) error
	FindUserByName(string) (*model.User, error)
	CheckPwd(userName string, pwd string) (isOk bool, err error)
}

func NewUserDataService(userRepository repository.IUserRepository) IUserDataService {
	return &UserDataService{UserRepository: userRepository}
}

type UserDataService struct {
	UserRepository repository.IUserRepository
}

func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

func ValidatePassword(userPassword string, hashed string) (isOk bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, errors.New("密码对比错误")
	}
	return true, nil
}

func (u *UserDataService) AddUser(user *model.User) (userID int64, err error) {
	pwdByte, err := GeneratePassword(user.HashPassword)
	if err != nil {
		return user.ID, err
	}
	user.HashPassword = string(pwdByte)

	return u.UserRepository.CreateUser(user)
}

func (u *UserDataService) DeleteUser(userID int64) error {
	return u.UserRepository.DeleteUserByID(userID)
}

func (u *UserDataService) UpdateUser(user *model.User, isChangePwd bool) error {
	if isChangePwd {
		pwdByte, err := GeneratePassword(user.HashPassword)
		if err != nil {
			return err
		}
		user.HashPassword = string(pwdByte)
	}
	return u.UserRepository.UpdateUser(user)
}

func (u *UserDataService) FindUserByName(name string) (user *model.User, err error) {
	return u.UserRepository.FindUserByName(name)
}

func (u *UserDataService) CheckPwd(userName, pwd string) (isOk bool, err error) {
	user, _ := u.UserRepository.FindUserByName(userName)
	return ValidatePassword(pwd, user.HashPassword)
}
