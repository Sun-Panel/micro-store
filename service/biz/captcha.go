package biz

import (
	"errors"
	"sun-panel/lib/captcha"

	"github.com/gin-gonic/gin"
)

type CaptchaType struct {
}

func (cpt *CaptchaType) CaptchaVerifyHandle(c *gin.Context, vcode string, clear bool) error {

	var captchaId string
	var err error

	// 获取captchaId
	if captchaId, err = captcha.CaptchaGetIdByCookieHeader(c, "CaptchaId"); err != nil {
		return errors.New("图形验证码错误，请尝试刷新页面重试")
	}

	return cpt.CaptchaVerifyHandleByCaptchaId(captchaId, vcode, clear)
}

func (cpt *CaptchaType) CaptchaVerifyHandleByCaptchaId(captchaId string, vcode string, clear bool) error {
	// 验证码错误
	if !captcha.CaptchaVerifyHandle(captchaId, vcode, clear) {
		captcha.CaptchaVerifyHandle(captchaId, vcode, true) // 此为清除验证码
		return errors.New("图形验证码错误")
	}
	return nil
}
