package microapp

type MicroAppPublic struct {
}

type MicroAppVersionGetVersionListReq struct {
	Page        int  `json:"page"`
	Limit       int  `json:"limit"`
	AppRecordId uint `json:"appRecordId"`
}

type MicroAppVersionGetListReq struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Order      string `json:"order"`
	CategoryId uint   `json:"categoryId"`
	Keyword    string `json:"keyword"`
}
