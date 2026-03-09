package adminApiStructs

import (
	"sun-panel/api/api_v1/common/apiData/commonApiStructs"
	"sun-panel/models/common"
	"time"
)

type GoodsAddReq struct {
	Title         string                 `json:"title"`         // 商品标题
	Price         float64                `json:"price"`         // 价格
	OriginalPrice float64                `json:"originalPrice"` // 原价
	Discount      string                 `json:"discount"`      // 优惠活动
	Description   string                 `json:"description"`   // 描述
	Param         map[string]interface{} `json:"param"`         // 商品参数
	Num           int                    `json:"num"`
}

type GoodsUpdateReq struct {
	GoodsAddReq
	Id uint `json:"id"`
}

type GoodsUpdateSaleReq struct {
	commonApiStructs.InfoIdReq
	Sort   *int  `json:"sort"`
	Status *bool `json:"status"` // 上下架数据
	Num    *int  `json:"num"`
	Id     uint  `json:"id"`
}

type GoodsGetListItemResp struct {
	common.CommonGoodsStruct
	Sort           int       `json:"sort"`
	Status         int       `json:"status"` // 上下架数据
	Num            int       `json:"num"`
	LastSnapshotId uint      `json:"lastSnapshotId"`
	UserId         uint      `json:"userId"`
	Id             uint      `json:"id"`
	CreateTime     time.Time `json:"createTime"`
	UpdateTime     time.Time `json:"updateTime"`
}
