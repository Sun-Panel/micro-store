package adminApiStructs

type CommonEmail struct {
	MailTo []string `json:"mailTo"`
	// SendName         string            `json:"send_name"`
	Title          string            `json:"title"`
	Body           string            `json:"body"`
	TemplateArg    map[string]string `json:"templateArg"`    // 模板变量
	IsSeparateSend bool              `json:"isSeparateSend"` // 是否单独发送
}

type EmailSendReq struct {
	CommonEmail
	ReplaceArg bool `json:"replaceArg"` // 替换参数
}

type EmailSendByTemplateIdReq struct {
	CommonEmail
	EmailTemplateId uint `json:"emailTemplateId"` // 模板id
}
