package system

import (
	"sun-panel/api/api_v1/common/apiData/commonApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/lib/cmn/systemSetting"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type registerInfo struct {
	Email        string                               `json:"email"`
	UserName     string                               `json:"userName"`
	Password     string                               `json:"password"`
	Vcode        string                               `json:"vcode"`
	EmailVCode   string                               `json:"emailVCode"`
	Verification commonApiStructs.VerificationRequest `json:"verification"`
	ReferralCode string                               `json:"referralCode"`
}

const EmailCodeCapacity = 1000

type RegisterApi struct{}

// 注册提交（开始注册）
func (l *RegisterApi) Commit(c *gin.Context) {
	req := registerInfo{}
	err := c.ShouldBindJSON(&req)
	req.Email = req.UserName
	if err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	errMsg, err := base.ValidateInputStruct(req)
	if err != nil {
		apiReturn.Error(c, errMsg)
		return
	}

	// 图形验证码验证
	if err := biz.Captcha.CaptchaVerifyHandle(c, req.Vcode, true); err != nil {
		apiReturn.Error(c, err.Error())
		return
	}

	// 验证是否开启注册和后缀格式是否正确
	{
		systemSettingInfo := systemSetting.ApplicationSetting{}
		if err := global.SystemSetting.GetValueByInterface("system_application", &systemSettingInfo); err != nil || !systemSettingInfo.Register.OpenRegister {
			apiReturn.Error(c, global.Lang.Get("register.unopened_register"))
			return
		}

		if systemSettingInfo.Register.EmailSuffix != "" && !cmn.VerifyFormat("^.*"+systemSettingInfo.Register.EmailSuffix+"$", req.Email) {
			apiReturn.Error(c, global.Lang.GetWithFields("register.emailSuffix_error", map[string]string{"EmailSuffix": systemSettingInfo.Register.EmailSuffix}))
			return
		}
	}

	// 验证邮箱是否被注册
	{
		userCheck := &models.User{Mail: req.UserName}
		if _, err := userCheck.GetUserInfoByUsername(req.UserName); err == nil && err != gorm.ErrRecordNotFound {
			apiReturn.Error(c, global.Lang.Get("register.mail_exist"))
			return
		}
	}

	// 验证码验证
	{
		vCode, ok := global.VerifyCodeCachePool.Get(req.Email)
		if !ok {
			apiReturn.Error(c, global.Lang.Get("common.captcha_code_error"))
			//验证码不存在
			return
		}
		if vCode != req.EmailVCode {
			apiReturn.Error(c, global.Lang.Get("common.captcha_code_error"))
			return
			//验证码有误
		}
	}

	// 自动生成用户昵称
	name := "用户" + cmn.BuildRandCode(4, cmn.RAND_CODE_MODE3)

	//验证通过，注册
	user := &models.User{
		Mail:     req.UserName,
		Name:     name,
		Username: req.UserName,
		Password: cmn.PasswordEncryption(req.Password),
		Status:   1,
		Role:     2,
	}
	_, err = user.CreateOne()
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	//删除旧的验证码
	global.VerifyCodeCachePool.Delete(req.Email)

	// 填写推荐人信息
	// if req.ReferralCode != "" {
	// 	referrerInfo := models.User{}
	// 	if err := global.Db.Find(&referrerInfo).Error; err == nil {
	// 		newReferralReward := models.ReferralReward{
	// 			ReferrerId: referrerInfo.ID,
	// 			RefereeId:  user.ID,
	// 		}
	// 		global.Db.Create(&newReferralReward)
	// 	}

	// }

	apiReturn.Success(c)
}
