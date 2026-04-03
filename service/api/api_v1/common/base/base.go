package base

import (
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/biz"

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

type FormFileUploadOptions struct {
	// 表单文件字段名
	FormFileName string
	// 同意的文件扩展名,为空不检查
	AgreeExtNames []string
	// 允许的最大文件大小，单位字节，为0不检查
	MaxSize int64
	// 文件保存目录
	SaveDir string
}

type FormFileUploadFileInfo struct {
	// 文件储存路径
	FileSavePath string
	// 文件原始名称
	FileOriginalName string
	// 文件信息
	FormFile *multipart.FileHeader
	// 文件扩展名 eg: .zip
	Ext string
}

var (
	ErrUploadExtensionNameNotAllowed = errors.New("extension not allowed") // "上传文件类型不允许"
	ErrUploadExceedMaxSize           = errors.New("exceeds maximum size")  // "上传文件超出最大尺寸"
)

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

func GetCurrentDeveloper(c *gin.Context) models.Developer {
	if value, exist := c.Get("developerInfo"); exist {
		if v, ok := value.(models.Developer); ok {
			return v
		}
	}
	return models.Developer{}
}

func GetCurrentUserLang(c *gin.Context) string {
	if value, exist := c.Get("Lang"); exist {
		if v, ok := value.(string); ok {
			return v
		}
	}
	return "en"
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

func FormFileUpload(c *gin.Context, options FormFileUploadOptions) (uploadInfo FormFileUploadFileInfo, err error) {

	f, err := c.FormFile(options.FormFileName)
	if err != nil {
		return
	} else {
		fileExt := strings.ToLower(path.Ext(f.Filename))

		if len(options.AgreeExtNames) != 0 && !cmn.InArray(options.AgreeExtNames, fileExt) {
			err = ErrUploadExtensionNameNotAllowed
			return
		}

		if options.MaxSize != 0 && f.Size > options.MaxSize {
			err = ErrUploadExceedMaxSize
			return
		}

		fileName := cmn.Md5(fmt.Sprintf("%s%s", f.Filename, time.Now().String()))
		fildNamePrefix := fmt.Sprintf("%s/%s-%s-", options.SaveDir, time.Now().Format("20060102150405"), cmn.BuildRandCode(10, cmn.RAND_CODE_MODE2))
		isExist, _ := cmn.PathExists(options.SaveDir)
		if !isExist {
			os.MkdirAll(options.SaveDir, os.ModePerm)
		}
		filepath := fmt.Sprintf("%s%s%s", fildNamePrefix, fileName, fileExt)
		c.SaveUploadedFile(f, filepath)

		uploadInfo = FormFileUploadFileInfo{
			FileSavePath:     filepath,
			FileOriginalName: f.Filename,
			FormFile:         f,
			Ext:              fileExt,
		}
	}

	return
}

func UploadImageFile(c *gin.Context, agreeExtNames []string, maxSize int64) (fileInfo models.File, err error) {
	userInfo, _ := GetCurrentUserInfo(c)
	configUpload := global.Config.GetValueString("base", "source_path")

	// 文件类型
	// fileTypeStr := c.PostForm("type")
	// fileType, _ := strconv.Atoi(fileTypeStr)

	f, err := c.FormFile("imgfile")
	if err != nil {
		return
	} else {
		fileExt := strings.ToLower(path.Ext(f.Filename))

		if !cmn.InArray(agreeExtNames, fileExt) {
			err = ErrUploadExtensionNameNotAllowed
			return
		}

		if maxSize != 0 && f.Size > maxSize {
			err = ErrUploadExceedMaxSize
			return
		}

		fileName := cmn.Md5(fmt.Sprintf("%s%s", f.Filename, time.Now().String()))
		fildDir := fmt.Sprintf("%s/%d/%d/%d/", configUpload, time.Now().Year(), time.Now().Month(), time.Now().Day())
		isExist, _ := cmn.PathExists(fildDir)
		if !isExist {
			os.MkdirAll(fildDir, os.ModePerm)
		}
		filepath := fmt.Sprintf("%s%s%s", fildDir, fileName, fileExt)
		c.SaveUploadedFile(f, filepath)

		pureFilePath := filepath[len(configUpload):]
		// 向数据库添加记录
		fileInfo = models.File{
			UserId:   userInfo.ID,
			FileName: f.Filename,
			Src:      pureFilePath,
			Ext:      fileExt,
		}
	}

	return
}

func HandleBizErrorAndReturn(c *gin.Context, err error) {
	// 业务错误：转换为数字错误码，前端统一处理
	if errCode, ok := biz.IsBizError(err); ok {
		intCode := BizCodeToInt(errCode) // API层负责转换
		apiReturn.ErrorByCode(c, intCode)
		return
	}
	// 其他错误：数据库错误
	apiReturn.ErrorDatabase(c, err.Error())
	c.Abort()
}

// bizCodeToInt 业务错误码转数字错误码（API层负责转换）
func BizCodeToInt(code string) int {
	codeMap := map[string]int{
		biz.ErrCodeAppNotFound:          2000,
		biz.ErrCodeVersionNotFound:      2001,
		biz.ErrCodeVersionExists:        2002,
		biz.ErrCodeVersionCodeExists:    2003,
		biz.ErrCodeStatusNotAllowed:     2004,
		biz.ErrCodeApprovedCannotDelete: 2005,
		biz.ErrCodeNotPendingReview:     2006,
		biz.ErrCodeNoUpdateContent:      2007,
		// 微应用开发者相关 3000-3099
		biz.ErrCodeAppIdExists:         3000,
		biz.ErrCodeNoPermission:        3001,
		biz.ErrCodePendingReviewExists: 3002,
		biz.ErrCodeNoPendingReviewApp:  3003,
		biz.ErrCodeInvalidParam:        3004,
	}

	if intCode, ok := codeMap[code]; ok {
		return intCode
	}
	return -1
}

// serveFile 流式传输文件（支持断点续传）
func ServeFile(c *gin.Context, filePath string) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	// 设置响应头
	fileName := fileInfo.Name()
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")
	c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	// 支持 Range 请求（断点续传）
	http.ServeContent(c.Writer, c.Request, fileName, fileInfo.ModTime(), file)
}

// serveFileNonStreaming 非流式传输文件（一次性读取，可检测完整传输）
func ServeFileNonStreaming(c *gin.Context, filePath string) {
	// 读取整个文件到内存
	data, err := os.ReadFile(filePath)
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	// 设置响应头
	fileName := filepath.Base(filePath)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")
	c.Header("Content-Length", strconv.Itoa(len(data)))

	// 写入响应（完成后可以执行token过期等操作）
	c.Data(http.StatusOK, "application/octet-stream", data)
}
