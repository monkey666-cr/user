package handler

import (
	"context"
	"user/domain/model"
	"user/domain/service"

	user "user/proto/user"
)

type User struct {
	UserDataService service.IUserDataService
}

func (u *User) Register(ctx context.Context, in *user.UserRegisterRequest,
	out *user.UserRegisterResponse) error {
	userRegister := &model.User{
		UserName:     in.UserName,
		FirstName:    in.FirstName,
		HashPassword: in.Pwd,
	}
	_, err := u.UserDataService.AddUser(userRegister)
	if err != nil {
		return err
	}

	out.Message = "注册成功"

	return nil
}

func (u *User) Login(ctx context.Context, in *user.UserLoginRequest,
	out *user.UserLoginResponse) error {
	isOk, err := u.UserDataService.CheckPwd(in.UserName, in.Pwd)
	if err != nil {
		return err
	}

	out.IsSuccess = isOk

	return nil
}

func (u *User) GetUserInfo(ctx context.Context, in *user.UserInfoRequest,
	out *user.UserInfoResponse) error {
	userModel, err := u.UserDataService.FindUserByName(in.UserName)
	if err != nil {
		return err
	}

	out.UserName = userModel.UserName
	out.FirstName = userModel.FirstName
	out.UserId = userModel.ID

	return nil
}
