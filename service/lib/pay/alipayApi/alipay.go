package alipayApi

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/go-pay/xlog"
)

type PagePayParam struct {
	Appid      string
	PrivateKey string  // 私钥
	ReturnUrl  string  // 支付完成返回的地址
	NotifyUrl  string  // 异步通知URL
	Title      string  // 支付显示的标题
	OrderNo    string  // 订单编号
	Amount     float64 // 订单金额
}

type PayNotifyResultAsyncResult struct {
	TradeNo     string `json:"trade_no"`
	OutTradeNo  string `json:"out_trade_no"`
	TradeStatus string `json:"trade_status"`
}

const (
	WAIT_BUYER_PAY = 1
	TRADE_CLOSED   = 4
	TRADE_SUCCESS  = 2
	TRADE_FINISHED = 2
)

func PagePay(param PagePayParam) (string, error) {
	client, err := alipay.NewClient(param.Appid, param.PrivateKey, true)
	if err != nil {
		return "", err
	}
	client.DebugSwitch = gopay.DebugOn

	// 设置需要的参数
	client.SetLocation(alipay.LocationShanghai).SetReturnUrl(param.ReturnUrl).SetNotifyUrl(param.NotifyUrl)

	// // 设置支付宝请求 公共参数
	// //    注意：具体设置哪些参数，根据不同的方法而不同，此处列举出所有设置参数
	// client.SetLocation(alipay.LocationShanghai). // 设置时区，不设置或出错均为默认服务器时间
	// 	// 						SetCharset(alipay.UTF8).             // 设置字符编码，不设置默认 utf-8
	// 	// 						SetSignType(alipay.RSA2).            // 设置签名类型，不设置默认 RSA2
	// 	SetReturnUrl("https://baidu.com").                                     // 设置返回URL
	// 	SetNotifyUrl("http://120.46.223.120:3912/api/goodsOrder/aliPayResult") // 设置异步通知URL
	// 	// SetAppAuthToken()                    // 设置第三方应用授权

	// 开启支付部分
	bm := make(gopay.BodyMap)

	bm.Set("subject", param.Title).
		// Set("scene", "bar_code").
		// Set("auth_code", "286248566432274952").
		Set("out_trade_no", param.OrderNo).
		Set("total_amount", param.Amount).
		Set("timeout_express", "15m")

	payUrl, err := client.TradePagePay(context.Background(), bm)
	if err != nil {
		if bizErr, ok := alipay.IsBizError(err); ok {
			xlog.Errorf("%+v", bizErr)
			// do something
			return "", err
		}
		xlog.Errorf("client.TradePay(%+v),err:%+v", bm, err)
		return "", err
	}

	return payUrl, nil
}

// （支付宝）异步调用支付结果
// （需要先设置支付宝公钥 SetAliPayPublicKey。如果成功会直接返回成功给支付宝）
func PayNotifyResultAsync(c *gin.Context, aliPayPublicKey string, result interface{}) (err error) {
	notifyReq, err := alipay.ParseNotifyToBodyMap(c.Request) // c.Request 是 gin 框架的写法
	if err != nil {
		xlog.Error(err)
		return
	}
	// 支付宝异步通知验签（公钥模式）
	if ok, err := alipay.VerifySign(aliPayPublicKey, notifyReq); !ok || err != nil {
		return err
	}

	if err := notifyReq.Unmarshal(result); err != nil {
		return err
	}

	// 返回给支付宝success
	c.String(http.StatusOK, "%s", "success")
	return nil
}

// （支付宝）查询支付结果
// 相关文档：https://opendocs.alipay.com/open/02e7gm
func QueryPayResult(c *gin.Context, appid, privateKey, order_no string) (aliRsp *alipay.TradeQueryResponse, err error) {
	client, err := alipay.NewClient(appid, privateKey, true)
	if err != nil {
		return aliRsp, err
	}
	bm := gopay.BodyMap{
		"out_trade_no": order_no,
		// "trade_no":     "",
	}
	return client.TradeQuery(c, bm)
}
