package clientApi

import (
	"sun-panel/lib/sunStore/request"
)

type GoodsOrder struct {
	ClientAPI *ClientAPI
}

type GoodsOrderUpdateAtatusReq struct {
	OrderNo string `json:"orderNo"`
	Status  string `json:"status"`
}

func NewGoodsOrder(c *ClientAPI) *GoodsOrder {
	user := GoodsOrder{
		ClientAPI: c,
	}
	return &user
}

func (g *GoodsOrder) UpdateStatus(req GoodsOrderUpdateAtatusReq) error {
	url := g.ClientAPI.GetHost() + "/openApi/v1/c/goodsOrder/updateStatus"

	httpCode, _, err := g.ClientAPI.Post(url, req, nil)
	if request.DeadlyError(httpCode) != nil {
		return err
	}
	return nil
}
