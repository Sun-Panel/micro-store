package common

type CommonGoodsStruct struct {
	Title         string  `json:"title" gorm:"type:varchar(100)"`          // 商品标题
	Price         float64 `json:"price" gorm:"type:DECIMAL(10,2)"`         // 价格
	OriginalPrice float64 `json:"originalPrice" gorm:"type:DECIMAL(10,2)"` // 原价
	Discount      string  `json:"discount" gorm:"varchar(100)"`            // 优惠活动
	Description   string  `json:"description"`                             // 描述
	Param         string  `json:"param" gorm:"text"`                       // 商品参数,json串
}
