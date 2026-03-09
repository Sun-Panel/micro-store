package sunStore

import (
	"fmt"
	"sun-panel/api/api_v1/sunStore/apiResp"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang-jwt/jwt/v4"
)

type Api struct {
}

type GoodsBuyQualificationCheckReq struct {
	GoodsId int64 `json:"goods_id"`
}

type JWTClaims struct {
	Sub string `json:"sub"` // 主体
	Iat int    `json:"iat"`
	Exp int    `json:"exp"`
	jwt.RegisteredClaims
}

func ParseJWT(tokenString string, secretKey string) (*JWTClaims, error) {
	// 定义解析函数
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// 类型断言并返回声明
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func checkHeaderJwt(c *gin.Context, secretKey string) (*JWTClaims, error) {
	// 从请求头中获取 JWT token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("authorization header is required")
	}

	// 简单提取 Bearer token (实际项目中可能需要更复杂的处理)
	tokenString := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	} else {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	// 解析和验证 JWT
	claims, err := ParseJWT(tokenString, secretKey)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	// // 验证时间是否在有效期内
	// if claims.ExpiresAt.Time.Before(jwt.TimeFunc()) {
	// 	return nil, fmt.Errorf("token has expired")
	// }

	// claims.exp

	return claims, nil
}

type OrderQualificationApiReq struct {
	Email      string            `json:"email"`
	Sku        string            `json:"sku"`
	Number     int               `json:"number"`
	ExtendData map[string]string `json:"extendData"`
	RequestID  string            `json:"requestId"`
}

type OrderQualificationApiResp struct {
	Qualification bool `json:"qualification"`
}

// OnlyProGoodsBuyQualificationCheck 仅PRO可购买的商品购买资格检查
func (api *Api) OnlyProGoodsBuyQualificationCheck(c *gin.Context) {
	secretKey, err := global.SystemSetting.GetVariableString("openAPI_jwt_secret_key")
	if err != nil {
		global.Logger.Errorln("GetVariableString failed", err)
		apiResp.ErrorResponse(c, apiResp.ErrCodeAuthFailed, err.Error())
		return
	}

	if _, err := checkHeaderJwt(c, secretKey); err != nil {
		global.Logger.Errorln("JWT check failed", err)
		apiResp.ErrorResponse(c, apiResp.ErrCodeAuthFailed, err.Error())
		return
	}

	req := OrderQualificationApiReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		global.Logger.Errorln("ShouldBindBodyWith failed", err)
		apiResp.ErrorResponse(c, apiResp.ErrCodeBadRequest, err.Error())
		return
	}

	// 对比货号
	skus := []string{}
	getSkusErr := global.SystemSetting.GetVariableByInterface("only_pro_qualification_buy_skus", &skus)

	if getSkusErr != nil {
		global.Logger.Errorln("GetVariableString failed", err)
		apiResp.ErrorResponse(c, apiResp.ErrCodeBadRequest, err.Error())
		return
	}

	if !cmn.InArray(skus, req.Sku) {
		apiResp.ErrorResponse(c, apiResp.ErrCodeNoPurchaseQualification, "")
		return
	}

	global.Logger.Infoln("GoodsBuyQualificationCheck", cmn.AnyToJsonStr(req))
	userEmail := req.Email

	var user models.User
	if err := global.Db.Where("username = ?", userEmail).First(&user).Error; err != nil {
		global.Logger.Errorln("User not found", err)
		apiResp.ErrorResponse(c, apiResp.ErrCodeErrUserNotFound, "")
		return
	}

	// 查询用户关联的ProAuthorize信息
	var proAuth models.ProAuthorize
	if err := global.Db.Where("user_id = ?", user.ID).First(&proAuth).Error; err != nil {
		global.Logger.Errorln("ProAuthorize not found", err)
		apiResp.ErrorResponse(c, apiResp.ErrCodeNoPurchaseQualification, "Authorization information not found")
		return
	}

	// 检查过期时间是否小于30天
	expiredTime := proAuth.ExpiredTime
	now := time.Now()
	diff := expiredTime.Sub(now)

	// 如果剩余时间小于30天，则没有购买资格
	if diff < 30*24*time.Hour {
		apiResp.ErrorResponse(c, apiResp.ErrCodeNoPurchaseQualification, "PRO is valid for less than 30 days")
		return
	}

	apiResp.SuccessResponse(c, OrderQualificationApiResp{Qualification: true})
}
