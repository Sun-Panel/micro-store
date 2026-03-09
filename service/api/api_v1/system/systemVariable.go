package system

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/global"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type SystemVariableApi struct {
}

// func (a *SystemVariableApi) Set(c *gin.Context) {
// 	req := systemApiStructs.SystemVariable{}

// 	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
// 		apiReturn.ErrorParamFomat(c, err.Error())
// 		return
// 	}

// 	if err := global.SystemSetting.SetVariable(req.Name, req.Value); err != nil {
// 		apiReturn.Error(c, err.Error())
// 	} else {
// 		apiReturn.Success(c)
// 	}
// }

func (a *SystemVariableApi) GetMultiple(c *gin.Context) {
	req := []string{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	reap, _ := global.SystemSetting.GetMultipleVariableString(req...)

	apiReturn.SuccessData(c, reap)
}

type SystemVariableApiGet struct {
	Name string
}

func (a *SystemVariableApi) Get(c *gin.Context) {
	req := SystemVariableApiGet{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	reap := ""
	if v, err := global.SystemSetting.GetVariableString(req.Name); err != nil {
		apiReturn.Error(c, err.Error())
		return
	} else {
		reap = v
	}

	apiReturn.SuccessData(c, reap)
}
