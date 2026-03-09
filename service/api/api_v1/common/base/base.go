package base

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/global"
	"sun-panel/lib/captcha"
	"sun-panel/lib/cmn"
	"sun-panel/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

type PageLimitVerify struct {
	Page  int64
	Limit int64
}

const (
	VISIT_MODE_LOGIN = iota
	VISIT_MODE_PUBLIC
)

const (
	GIN_GET_VISIT_MODE = "VISIT_MODE"
)

// 验证输入是否有效并返回错误
func validateInputStruct(params interface{}) (errMsg string, err error) {
	var validate = validator.New()
	//通过label标签返回自定义错误内容
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		if label == "" {
			return field.Name
		}
		return label
	})

	// 自定义验证规则，使用 strings.TrimSpace 函数删除前后空格
	validate.RegisterValidation("trimmedRequired", func(fl validator.FieldLevel) bool {
		return strings.TrimSpace(fl.Field().String()) != ""
	})

	if err = validate.Struct(params); err != nil {
		trans := validateTransInit(validate)
		verrs := err.(validator.ValidationErrors)
		// errs := make(map[string]string)
		for _, value := range verrs.Translate(trans) {
			// errs[key[strings.Index(key, ".")+1:]] = value
			errMsg += " " + value
		}
		// fmt.Println(errs)
	}
	return
}

// 验证输入是否有效并返回错误
func ValidateInputStruct(params interface{}) (errMsg string, err error) {
	return validateInputStruct(params)
}

// 数据验证翻译器
func validateTransInit(validate *validator.Validate) ut.Translator {
	// 万能翻译器，保存所有的语言环境和翻译数据
	uni := ut.New(zh.New())
	// 翻译器
	trans, _ := uni.GetTranslator("zh")
	//验证器注册翻译器
	err := zhTranslations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		fmt.Println(err)
	}
	return trans
}

func GetCurrentUserInfo(c *gin.Context) (userInfo models.User, exist bool) {
	if value, exist := c.Get("userInfo"); exist {
		if v, ok := value.(models.User); ok {
			return v, exist
		}
	}
	return
}

func GetUserTimezoneLocation(c *gin.Context) *time.Location {
	timezone := c.Request.Header.Get("Timezone")
	location, err := time.LoadLocation(timezone)

	if err != nil {
		location = time.Now().Location()
	}
	return location
}

// 获取当前访问模式
func GetCurrentVisitMode(c *gin.Context) (visitMode int) {
	if value, exist := c.Get(GIN_GET_VISIT_MODE); exist {
		if v, ok := value.(int); ok {
			return v
		}
	}
	return
}

// 获取缓存会话的token Ctoken
func GetToken(c *gin.Context) string {
	cToken := c.GetHeader("token")
	return cToken
}

// 验证器验证
func VerificationCheck(verificationId, vCode string) (errCode int, verificationIdRes string) {

	// 需要进一步验证并返回验证信息
	if verificationId == "" || vCode == "" {
		verificationIdRes = cmn.BuildRandCode(16, cmn.RAND_CODE_MODE1)
		errCode = apiReturn.ERROR_CODE_VERIFICATION_MUST
		return
	}

	// 验证码错误
	if !captcha.CaptchaVerifyHandle(verificationId, vCode, true) {
		errCode = apiReturn.ERROR_CODE_VERIFICATION_FAIL
		return
	}
	errCode = apiReturn.ERROR_CODE_SUCCESS
	return
}

// Deprecated: user base.ConvertSystemTimeToUserTime
// 时间转用户时区时间
func ConvertTimeToUserTime(c *gin.Context, t time.Time) time.Time {
	// 获得用户时区
	userTimezone := c.Request.Header.Get("Timezone")

	// 获取当前时间的时区
	currentLocation := t.Location()

	// 尝试解析提供的时区字符串，出问题返回默认的时区
	newLocation, err := time.LoadLocation(userTimezone)
	if err != nil {
		global.Logger.Errorln("timezone conver fail. user timezone:", userTimezone, "system timezone:", currentLocation)
		return t
	}

	// 检查时区是否相同
	if currentLocation.String() == newLocation.String() {
		// 时区相同，不进行转换
		return t
	}

	// 时区不同，进行转换
	return t.In(newLocation)
}

func ConvertSystemTimeToUserTime(userLocaltion *time.Location, t time.Time) time.Time {
	// 时区不同，进行转换
	return t.In(userLocaltion)
}

// ConvertUserTimeToSystemTime 将用户时区时间转换为系统时区时间
// 通常用于将用户前端传来的时间转换为系统时间，用于数据库查询
func ConvertUserTimeToSystemTime(userLocation *time.Location, userTimeZoneTime time.Time) time.Time {

	// 获取当前时间的时区
	currentLocation := time.Local

	// 时区不同，进行转换
	return userTimeZoneTime.In(currentLocation)
}

func ConvertSQLNullTimeUserTimeToString(userLocaltion *time.Location, sqlNullTime sql.NullTime, format string) string {
	// global.Logger.Debugln("nulltime", sqlNullTime.Time, "是否有效", sqlNullTime.Valid)
	if sqlNullTime.Valid {
		return ConvertSystemTimeToUserTime(userLocaltion, sqlNullTime.Time).Format(format)
	}
	return ""
}
