package openApi

import (
	"encoding/json"
	"errors"
	"strings"
	"sun-panel/lib/sunStore/request"
)

type OpenApi struct {
	AccessToken string
	Host        string // 请求地址
}

type Host interface {
	GetHost() string
}

func NewOpenApi(Host, accessToken string) *OpenApi {
	openApi := OpenApi{}
	openApi.AccessToken = accessToken
	openApi.Host = Host
	return &openApi
}

func (o *OpenApi) Post(url string, requestData interface{}, responseData interface{}) (httpCode int, apiCode int, err error) {
	reqByte, err := json.Marshal(requestData)
	// fmt.Println("请求数据", string(reqByte))
	if err != nil {
		return httpCode, apiCode, err
	}

	header := map[string]string{
		"Authorization": "Bearer " + o.AccessToken,
	}

	respContent, httpResp, err := request.SendPostHeaderRequest(url, header, reqByte)
	if err != nil {
		return httpCode, apiCode, err
	}
	// fmt.Println("接口返回总数据", string(respContent))

	httpCode = httpResp.StatusCode

	if httpCode == 500 {
		return httpCode, httpCode, errors.New("server error")
	}

	// 第一轮解析
	requestResp := request.RequestRespDebug{}
	if err := json.Unmarshal(respContent, &requestResp); err != nil {
		return httpCode, requestResp.Code, errors.New("requestResp error :" + err.Error())
	}

	if requestResp.Code != 0 {
		return httpCode, requestResp.Code, errors.New(requestResp.Msg)
	}

	requestRespData, _ := json.Marshal(requestResp.Data)
	// fmt.Println("data数据", string(requestRespData))
	// 第二轮解析
	if err := json.Unmarshal(requestRespData, &responseData); err != nil {
		return httpCode, requestResp.Code, errors.New("requestResp data error :" + err.Error())
	}

	return httpCode, 0, nil
}

func (o *OpenApi) GetHost() string {
	return strings.TrimRight(o.Host, "/")
}
