package systemApiStructs

type MdPageInfo struct {
	IsLogin                 bool   `json:"isLogin"`
	Content                 string `json:"content"`
	MessageTemplateFlag     string `json:"messageTemplateFlag"`
	MessageTemplatePosition string `json:"messageTemplatePosition"` // top|bottom
}

type MdPageInfoEditReq struct {
	MdPageName        string `json:"mdPageName"`
	MdPageDescription string `json:"mdPageDescription"`
	MdPageInfo
}

type MdPageListItem struct {
	MdPageName        string `json:"mdPageName"`
	MdPageDescription string `json:"mdPageDescription"`
	MdPageInfo
}
