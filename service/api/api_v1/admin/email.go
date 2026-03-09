package admin

import (
	"sun-panel/api/api_v1/common/apiData/adminApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/global"
	"sun-panel/lib/cmn/systemSetting"
	"sun-panel/lib/mail"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type EmailApi struct {
}

func (a *EmailApi) Send(c *gin.Context) {

	req := adminApiStructs.EmailSendReq{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	emailInfoConfig := systemSetting.Email{}
	global.SystemSetting.GetValueByInterface("system_email", &emailInfoConfig)
	emailInfo := mail.EmailInfo{
		Username: emailInfoConfig.Mail,
		Password: emailInfoConfig.Password,
		Host:     emailInfoConfig.Host,
		Port:     emailInfoConfig.Port,
	}

	var sendError error
	contentInfo := mail.ContentInfo{}
	contentInfo.SendName = "Sun-Panel"
	contentInfo.Title = req.Title
	contentInfo.Body = req.Body

	if req.ReplaceArg {
		contentInfo.RespaceArg = req.TemplateArg
	}

	// 单独发
	if req.IsSeparateSend {
		for _, v := range req.MailTo {
			contentInfo.MailTo = []string{v}
			sendError = mail.SendMail(mail.NewEmailer(emailInfo), contentInfo)
		}
	} else {
		contentInfo.MailTo = req.MailTo
		sendError = mail.SendMail(mail.NewEmailer(emailInfo), contentInfo)
	}

	if sendError != nil {
		apiReturn.Error(c, sendError.Error())
	} else {
		apiReturn.Success(c)
	}

}

// // 基于模板ID发送
// func (a *EmailApi) SendByTemplateId(c *gin.Context) {

// 	req := adminApiStructs.EmailSendByTemplateIdReq{}

// 	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
// 		apiReturn.ErrorParamFomat(c, err.Error())
// 		return
// 	}

// 	emailInfoConfig := systemSetting.Email{}
// 	global.SystemSetting.GetValueByInterface("system_email", &emailInfoConfig)
// 	emailInfo := mail.EmailInfo{
// 		Username: emailInfoConfig.Mail,
// 		Password: emailInfoConfig.Password,
// 		Host:     emailInfoConfig.Host,
// 		Port:     emailInfoConfig.Port,
// 	}

// 	tmpInfo := models.EmailTemplate{}
// 	if err := global.Db.First(&tmpInfo).Error; err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			apiReturn.ErrorByCodeAndMsg(c, -1, "没有找到邮件模板")
// 		} else {
// 			apiReturn.ErrorDatabase(c, err.Error())
// 		}

// 		return
// 	}

// 	sendName := "Sun-Panel"
// 	sendError := mail.SendMailAndArg(mail.NewEmailer(emailInfo), req.MailTo, sendName, tmpInfo.Title, tmpInfo.Content, req.TemplateArg, req.IsSeparateSend)

// 	if sendError != nil {
// 		apiReturn.ErrorByCode(c, -1)
// 	} else {
// 		apiReturn.Success(c)
// 	}
// }
