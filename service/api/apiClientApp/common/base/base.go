package base

import (
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// const AES_KEY = "wAOhhAqbBxgjNOlcEAGvwawvUZmWDLkN"

const RUN_MODE = "a" // debug /a

const timestampDifference = 172800 // 48小时

type AESPostParam struct {
	P string
}

type ClientBase struct {
	ClientId   string
	Version    string
	MacAddress string
	LanIP      string
	Timestamp  int64
}

type ProStatus struct {
	ExpirationTimestamp int64
}

type Claims struct {
	Username  string     `json:"username"`
	Client    ClientBase `json:"client"`
	ProStatus ProStatus  `json:"proStatus"`

	jwt.RegisteredClaims // 嵌入标准claims字段
}

func GetCurrentClient(c *gin.Context) *ClientBase {
	v, exist := c.Get("SunPanelClientInfo")
	if !exist {
		return nil
	}

	if clientBase, ok := v.(ClientBase); ok {
		return &clientBase
	}
	return nil
}

func GetProStatus(c *gin.Context) *ProStatus {
	v, exist := c.Get("proStatus")
	if !exist {
		return nil
	}

	if proStatus, ok := v.(ProStatus); ok {
		return &proStatus
	}
	return nil
}

func GetCurrentUser(c *gin.Context) (models.User, bool) {
	v, exists := c.Get("currentUser")
	if !exists {
		return models.User{}, false
	}
	user, ok := v.(models.User)
	return user, ok
}

// // Deprecated: Please use AESEncryptWithKey instead.
// func AESEncrypt(content string) (string, error) {
// 	return AES.Encrypt(AES_KEY, content)
// }

// // Deprecated: Please use AESDecryptWithKey instead.
// func AESDecrypt(content string) (string, error) {
// 	return AES.Decrypt(AES_KEY, content)
// }

// func AESEncryptWithKey(key string, content string) (string, error) {
// 	return AES.Encrypt(key, content)
// }

// func AESDecryptWithKey(key string, content string) (string, error) {
// 	return AES.Decrypt(key, content)
// }

// func DecryptParam(c *gin.Context, param interface{}) error {
// 	req := AESPostParam{}
// 	sceretKey := GetVersionSecretKey(c)
// 	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
// 		return errors.New("param wrong format:" + err.Error())
// 	}
// 	if decryptStr, err := AESDecryptWithKey(sceretKey, req.P); err != nil {
// 		return errors.New("decrypt error:" + err.Error())
// 	} else {
// 		if err := json.Unmarshal([]byte(decryptStr), param); err != nil {
// 			return errors.New("JSON encrypted string error:" + err.Error())
// 		}
// 	}
// 	return nil
// }

// func EncryptParam(sceretKey string, param interface{}) (string, error) {

// 	if byt, err := json.Marshal(param); err != nil {
// 		return "", errors.New("json error:" + err.Error())
// 	} else {
// 		if encrypt, err := AESEncryptWithKey(sceretKey, string(byt)); err != nil {
// 			return "", errors.New("encrypt error:" + err.Error())
// 		} else {
// 			return encrypt, nil
// 		}
// 	}
// }

// func GetUserInfoWithCheckByCToken(c *gin.Context, token string) (*models.User, error) {
// 	authServiceClientTokenUser := global.AuthServiceClientTokenUser{}
// 	if btoken, ok := global.CUserAuthServiceClientToken.Get(token); !ok {
// 		global.CUserAuthServiceClientToken.Delete(token)
// 		return nil, errors.New("token expired")
// 	} else {
// 		if v, ok := global.UserAuthServiceClientToken.Get(btoken); !ok {
// 			global.UserAuthServiceClientToken.Delete(btoken)
// 			global.CUserAuthServiceClientToken.Delete(token)
// 			return nil, errors.New("token expired -btoken")
// 		} else {
// 			authServiceClientTokenUser = v
// 			// 此处采用新的验证方式，支持多设备登录，将不支持此验证 - 2024-9-29
// 			// 	if authServiceClientTokenUser.Ctoken != token {
// 			// 		global.CUserAuthServiceClientToken.Delete(token)
// 			// 		return nil, errors.New("token expired -token")
// 			// 	}
// 		}
// 	}
// 	return &authServiceClientTokenUser.User, nil
// }

// func GetRequestParam(c *gin.Context, param interface{}) error {
// 	if RUN_MODE == "debug" {
// 		if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
// 			return errors.New("param wrong format:" + err.Error())
// 		}
// 	} else {
// 		return DecryptParam(c, &param)
// 	}

// 	return nil
// }

// func GetRequestResp(sceretKey string, param interface{}) (any, error) {
// 	if RUN_MODE == "debug" {
// 		return param, nil
// 	} else {
// 		return EncryptParam(sceretKey, &param)
// 	}
// }

// func CheckTimestamp(timestamp int64) bool {
// 	if RUN_MODE == "debug" {
// 		return true
// 	}
// 	currentTimestamp := time.Now().Unix()
// 	TimeDifference := math.Abs(float64(currentTimestamp) - float64(timestamp))
// 	global.Logger.Debugln("接收", timestamp, "当前", currentTimestamp, "时间差", TimeDifference, "接受的时间差", timestampDifference)
// 	return TimeDifference < timestampDifference
// }

// func IsTimeDifferenceGreaterThan(thresholdDays int, timeStr1, timeStr2 time.Time) bool {
// 	// 计算时间差
// 	diff := timeStr2.Sub(timeStr1)

// 	// 将阈值转换为Duration
// 	thresholdDuration := time.Duration(thresholdDays) * 24 * time.Hour

// 	// 比较时间差和阈值
// 	return diff > thresholdDuration
// }

func GetVersionSecretKey(c *gin.Context) string {
	return c.GetString("secretKey")
}
