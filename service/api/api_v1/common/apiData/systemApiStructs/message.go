package systemApiStructs

import (
	"sun-panel/models"
	"time"
)

type MessageSendReq struct {
	Title        string            `json:"title"`
	Content      string            `json:"content"`
	ToUsernames  []string          `json:"toUsernames"` // 收消息的人(账号)
	Appendixs    []string          `json:"appendixs"`   // 附件（图片，多个使用英文逗号）
	TemplateArg  map[string]string `json:"templateArg"`
	TemplateFlag string            `json:"templateFlag"` // 模板标识
}

type MessageSendFailListResp struct {
	Email        string `json:"email"`
	ErrorMessage string `json:"errorMessage"`
}

type MessageGetReceiveMessageListReq struct {
	ReadStatus []int `json:"readStatus"`
}

type MessageGetReceiveMessageListResp struct {
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	ToUsernames []string `json:"toUsernames"` // 收消息的人(账号)
	Appendixs   []string `json:"appendixs"`   // 附件（图片，多个使用英文逗号）
}

type MessageUpdateReadStatustReq struct {
	MessageId  uint `json:"messageId"`
	ReadStatus bool `json:"readStatus"`
}

type MessageUpdateDeleteReq struct {
	MessageId uint `json:"messageId"`
	IsSend    bool `json:"isSend"` // 是否为发信方
}

type MessageInfo struct {
	MessageId    uint         `json:"messageId"`
	Title        string       `json:"title"`
	Content      string       `json:"content"`
	FromUserId   uint         `json:"fromUserId"`
	ToUserId     uint         `json:"toUserId"`
	Appendix     []string     `json:"appendix"`
	HaveRead     int          `json:"haveRead"`
	TopicId      string       `json:"topicId"`
	TemplateFlag string       `json:"templateFlag"`
	FromUser     *models.User `json:"fromUser"`
	ToUser       *models.User `json:"toUser"`
	CreateTime   time.Time    `json:"createTime"`
}
