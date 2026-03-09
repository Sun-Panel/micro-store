package openApi

import (
	"errors"
	"sun-panel/lib/sunStore/request"
)

type UserAuthUser struct {
	OpenApi *OpenApi
}

type UserInfoResp struct {
	// ID           string `json:"id"`
	// Status    int    `json:"status"`    // 状态 1.启用 2.停用 3.未激活

	Username   string `json:"username"`  // 账号
	Password   string `json:"password"`  // 密码
	Name       string `json:"name"`      // 名称
	HeadImage  string `json:"headImage"` // 头像地址
	Role       int    `json:"role"`      // 角色 1.管理员 2.普通用户
	Mail       string `json:"mail"`      // 邮箱
	SystemLang string `json:"systemLang"`
	Lang       string `json:"lang"`
	TimeZone   string `json:"timeZone"`
}

func NewUser(o *OpenApi) *UserAuthUser {
	user := UserAuthUser{
		OpenApi: o,
	}
	return &user
}

func (u *UserAuthUser) GetCurrentUserInfo() (UserInfoResp, int, error) {
	url := u.OpenApi.GetHost() + "/openApi/v1/u/user/getCurrentUserInfo"

	userInfo := UserInfoResp{}
	httpCode, _, err := u.OpenApi.Post(url, nil, &userInfo)

	if request.DeadlyError(httpCode) != nil {
		return userInfo, httpCode, errors.New(err.Error())
	}

	return userInfo, 0, err
}
