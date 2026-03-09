package webhook

type ParamMap map[string]interface{}

type CommonRequest struct {
	Event   string `json:"event"`
	EventID string `json:"eventId"`
}

type CommonGoodsStruct struct {
	Title             string                 `json:"title" gorm:"type:varchar(100)"`          // 商品标题
	Price             float64                `json:"price" gorm:"type:DECIMAL(10,2)"`         // 价格
	OriginalPrice     float64                `json:"originalPrice" gorm:"type:DECIMAL(10,2)"` // 原价
	Discount          string                 `json:"discount" gorm:"varchar(100)"`            // 优惠活动
	Description       string                 `json:"description"`                             // 描述
	CurrencyCode      string                 `json:"currencyCode" gorm:"varcha(20)"`          // 货币代码 CNY USD ...
	LimitBuyFrequency GoodsLimitBuyFrequency `json:"limitBuyFrequency" gorm:"json"`
	CustomData        ParamMap               `json:"customData" gorm:"json"`
}

type GoodsLimitBuyFrequency struct {
	BuyFrequency int `json:"buyFrequency"`
	Day          int `json:"day"`
}

type GoodsSnapshot struct {
	ID      uint `json:"id"`
	GoodsId uint `json:"goodsId"`
	CommonGoodsStruct
}

type User struct {
	ID         uint   `json:"id"`
	Username   string `json:"username"`
	Mail       string `json:"mail"`
	Name       string `json:"name"`
	HeadImage  string `json:"headImage"`
	Password   string `json:"password"`
	SystemLang string `json:"systemLang"`
	Lang       string `json:"lang"`
	TimeZone   string `json:"timeZone"`
}

type EventGoodsOrderData struct {
	OrderNo       string        `json:"orderNo"`
	Status        string        `json:"status"`
	CountPrice    float64       `json:"countPrice"`
	CurrencyCode  string        `json:"currencyCode"`
	PayPlatform   string        `json:"payPlatform"`
	GoodsId       uint          `json:"goodsId"`
	User          User          `json:"user"`
	Number        int           `json:"number"`
	GoodsSnapshot GoodsSnapshot `json:"goodsSnapshot"`
}

type EventGoodsOrderDataReq struct {
	CommonRequest
	Data EventGoodsOrderData `json:"data"`
}

type EventUserDataReq struct {
	CommonRequest
	Data User `json:"data"`
}
