package systemApiStructs

type HtmlPageInfo struct {
	IsLogin                 bool   `json:"isLogin"`
	Content                 string `json:"content"`
	MessageTemplateFlag     string `json:"messageTemplateFlag"`
	MessageTemplatePosition string `json:"messageTemplatePosition"` // top|bottom
}

type HtmlPageInfoEditReq struct {
	PageName        string `json:"pageName"`
	PageDescription string `json:"pageDescription"`
	HtmlPageInfo
}

type HtmlPageListItem struct {
	PageName        string `json:"pageName"`
	PageDescription string `json:"pageDescription"`
	HtmlPageInfo
}
