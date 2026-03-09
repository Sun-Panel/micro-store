package system

type ApiSystem struct {
	About    About
	LoginApi LoginApi
	UserApi  UserApi
	FileApi  FileApi
	// NoticeApi         NoticeApi
	// ModuleConfigApi   ModuleConfigApi
	MonitorApi MonitorApi
	// GoodsOrderApi     GoodsOrder
	// GoodsApi          Goods
	RegisterApi RegisterApi
	// MessageApi        MessageApi
	SystemVariableApi SystemVariableApi
	CaptchaApi        CaptchaApi
	MdPageApi         MdPageApi
}
