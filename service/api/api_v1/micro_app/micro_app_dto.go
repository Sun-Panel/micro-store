package microapp

type MicroAppPublic struct {
}

type MicroAppVersionGetListReq struct {
	Page        int  `json:"page"`
	Limit       int  `json:"limit"`
	AppRecordId uint `json:"appRecordId"`
}
