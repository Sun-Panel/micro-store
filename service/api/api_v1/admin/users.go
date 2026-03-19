package admin

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UsersApi struct {
}

// type ParamUserInfo struct {
// 	UserId    uint   `json:"userId"`
// 	Username  string `json:"username" validate:"required,email"`
// 	Password  string `json:"password" validate:"required"`
// 	Name      string `json:"name" `
// 	HeadImage string `json:"headImage" `
// 	Status    int    `json:"status" `
// 	Role      int    `json:"role" `
// 	Mail      string `json:"mail" `
// }

func (a UsersApi) Create(c *gin.Context) {
	param := models.User{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	param.Password = "-"
	if errMsg, err := base.ValidateInputStruct(param); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	mUser := models.User{
		Username:  param.Username,
		Password:  cmn.PasswordEncryption(param.Password),
		Name:      param.Name,
		HeadImage: param.HeadImage,
		Status:    param.Status,
		Role:      param.Role,
		Mail:      param.Username,
	}

	// 验证账号是否存在
	if _, err := mUser.CheckUsernameExist(param.Username); err != nil {
		apiReturn.Error(c, global.Lang.Get("register.mail_exist"))
		return
	}

	userInfo, err := mUser.CreateOne()

	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, gin.H{"userId": userInfo.ID})
}

func (a UsersApi) Deletes(c *gin.Context) {
	type UserIds struct {
		UserIds []uint
	}
	param := UserIds{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.GetAndInsert("common.api_error_param_format", "[", err.Error(), "]"))
		c.Abort()
		return
	}

	if err := global.Db.Delete(&models.User{}, &param.UserIds).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.Success(c)
}

// UpdateUserRequest 用户更新请求结构体（排除时间字段）
type UpdateUserRequest struct {
	ID        uint   `json:"id" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Name      string `json:"name"`
	HeadImage string `json:"headImage"`
	Status    int    `json:"status"`
	Role      int    `json:"role"`
	Mail      string `json:"mail"`
}

func (a UsersApi) Update(c *gin.Context) {
	// 使用专用结构体，避免时间字段解析问题
	var param UpdateUserRequest
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.GetAndInsert("common.api_error_param_format", "[", err.Error(), "]"))
		c.Abort()
		return
	}

	// 转换为 User 模型用于验证和更新
	userModel := models.User{
		Username:  param.Username,
		Name:      param.Name,
		HeadImage: param.HeadImage,
		Status:    param.Status,
		Role:      param.Role,
		Mail:      param.Username,
	}
	// 通过BaseModel设置ID
	userModel.BaseModel.ID = param.ID
	userModel.Password = "-" // 修改不允许修改密码，为了验证通过

	if errMsg, err := base.ValidateInputStruct(userModel); err != nil {
		apiReturn.ErrorParamFomat(c, errMsg)
		return
	}

	allowField := []string{"Username", "Name", "HeadImage", "Status", "Role", "Mail", "Token"}
	mUser := models.User{}

	// 获取当前用户的原始信息（用于清理token缓存）
	var currentUser models.User
	if err := global.Db.Where("id = ?", userModel.ID).First(&currentUser).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	// 验证账号是否存在（检查其他用户是否使用了这个用户名）
	if existingUser, err := mUser.CheckUsernameExist(userModel.Username); err != nil {
		// 用户名已存在
		if existingUser.ID != userModel.ID {
			apiReturn.Error(c, global.Lang.Get("register.mail_exist"))
			return
		}
	}

	userModel.Token = "" // 修改资料就重置token

	if err := global.Db.Select(allowField).Where("id=?", userModel.ID).Updates(&userModel).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	
	global.UserToken.Delete(currentUser.Token) // 更新用户信息
	
	// 返回token等基本信息
	apiReturn.SuccessData(c, userModel)
}

func (a UsersApi) GetList(c *gin.Context) {

	type ParamsStruct struct {
		models.User
		Limit   int
		Page    int
		Keyword string `json:"keyword"`
	}

	param := ParamsStruct{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.GetAndInsert("common.api_error_param_format", "[", err.Error(), "]"))
		c.Abort()
		return
	}

	var (
		list  []models.User
		count int64
	)
	db := global.Db

	// 查询条件
	if param.Keyword != "" {
		db = db.Where("name LIKE ? OR username LIKE ?", "%"+param.Keyword+"%", "%"+param.Keyword+"%")
	}

	if err := db.Limit(param.Limit).Offset((param.Page - 1) * param.Limit).Find(&list).Limit(-1).Offset(-1).Count(&count).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	// 构建返回数据，增加角色标签列表
	type UserResponse struct {
		ID         uint              `json:"id"`
		Username   string            `json:"username"`
		Name       string            `json:"name"`
		HeadImage  string            `json:"headImage"`
		Status     int               `json:"status"`
		Role       int               `json:"role"`
		Mail       string            `json:"mail"`
		CreateTime string            `json:"createTime"`
		Roles      []models.RoleInfo `json:"roles"` // 角色标签列表
	}

	resList := make([]UserResponse, len(list))
	for i, v := range list {
		resList[i] = UserResponse{
			ID:         v.ID,
			Username:   v.Username,
			Name:       v.Name,
			HeadImage:  v.HeadImage,
			Status:     v.Status,
			Role:       v.Role,
			Mail:       v.Mail,
			CreateTime: v.CreatedAt.Format("2006-01-02 15:04:05"),
			Roles:      models.GetRoleList(v.Role),
		}
	}

	apiReturn.SuccessListData(c, resList, count)
}

// UpdatePassword 管理员修改用户密码
func (a UsersApi) UpdatePassword(c *gin.Context) {
	type UpdatePasswordParam struct {
		UserID      uint   `json:"userId" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required,min=6"`
	}

	param := UpdatePasswordParam{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 查询用户是否存在
	mUser := models.User{}
	userInfo, err := mUser.GetUserInfoByUid(param.UserID)
	if err != nil {
		apiReturn.Error(c, "用户不存在")
		return
	}

	// 更新密码
	updateData := map[string]interface{}{
		"password": cmn.PasswordEncryption(param.NewPassword),
		"token":    "", // 修改密码后重置token，强制重新登录
	}

	if err := global.Db.Model(&models.User{}).Where("id=?", param.UserID).Updates(updateData).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	// 删除用户token缓存
	global.UserToken.Delete(userInfo.Token)

	apiReturn.SuccessData(c, gin.H{"userId": param.UserID})
}
