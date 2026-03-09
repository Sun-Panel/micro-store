package admin

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/biz"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type ClientBlackListIPApi struct {
}

func (b *ClientBlackListIPApi) GetList(c *gin.Context) {
	res := biz.ClientCache.BlacklistIP.GetAll()
	apiReturn.SuccessData(c, res)
}

func (b *ClientBlackListIPApi) Deletes(c *gin.Context) {
	req := []string{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if err := biz.ClientCache.BlacklistIP.Del(req...); err != nil {
		apiReturn.Error(c, err.Error())
		return
	}
	apiReturn.Success(c)
}

type ClientBlackListIPSetReq struct {
	IPs    []string `json:"ips"`
	Second int64    `json:"second"`
}

func (b *ClientBlackListIPApi) Set(c *gin.Context) {
	req := ClientBlackListIPSetReq{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	timestamp := time.Now().Add(time.Duration(req.Second) * time.Second)

	if err := biz.ClientCache.BlacklistIP.Set(timestamp, req.IPs...); err != nil {
		apiReturn.Error(c, err.Error())
		return
	}
	apiReturn.Success(c)
}
