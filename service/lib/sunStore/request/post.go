package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type RequestParam struct {
	P string
}

type RequestResp struct {
	Code int    `json:"code"`
	Data string `json:"data"`
	Msg  string `json:"msg"`
}

type RequestRespDebug struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

// 发送POST请求
// func SendPostRequest(url string, contentType string, requestData []byte) ([]byte, *http.Response, error) {
// 	// 发起 POST 请求
// 	resp, err := http.Post(url, contentType, bytes.NewBuffer(requestData))
// 	if err != nil {
// 		return nil, resp, fmt.Errorf("POST:Request error: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	// 检查响应状态码
// 	if resp.StatusCode != http.StatusOK {
// 		return nil, resp, fmt.Errorf("POST:Response status code error: %d", resp.StatusCode)
// 	}

// 	// 读取响应体
// 	responseBody, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, resp, fmt.Errorf("POST:Error reading response body: %v", err)
// 	}

// 	return responseBody, resp, nil
// }

// 发送POST请求带有header
func SendPostHeaderRequest(url string, headers map[string]string, requestData []byte) ([]byte, *http.Response, error) {
	// 创建一个新的请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestData))
	if err != nil {
		return nil, nil, fmt.Errorf("POST:Request creation error: %v", err)
	}

	// 添加自定义头部
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 创建一个HTTP客户端
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, resp, fmt.Errorf("POST:Request error: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, resp, fmt.Errorf("POST:Response status code error: %d", resp.StatusCode)
	}

	// 读取响应体
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("POST:Error reading response body: %v", err)
	}

	return responseBody, resp, nil
}

func SendPostJsonRequest(url string, requestData interface{}, responseData interface{}) (int, error, *http.Response) {

	reqByte, err := json.Marshal(requestData)
	// fmt.Println("请求数据", string(reqByte))
	if err != nil {
		return 0, err, nil
	}

	header := map[string]string{}
	header["Content-Type"] = "application/json"

	respContent, httpRes, err := SendPostHeaderRequest(url, header, reqByte)
	if err != nil {
		return 0, err, httpRes
	}
	// fmt.Println("接口返回总数据", string(respContent))

	// 第一轮解析
	requestResp := RequestRespDebug{}
	if err := json.Unmarshal(respContent, &requestResp); err != nil {
		return requestResp.Code, errors.New("requestResp error :" + err.Error()), httpRes
	}

	if requestResp.Code != 0 {
		return requestResp.Code, errors.New(requestResp.Msg), httpRes
	}

	requestRespData, _ := json.Marshal(requestResp.Data)
	// fmt.Println("data数据", string(requestRespData))
	// 第二轮解析
	if err := json.Unmarshal(requestRespData, &responseData); err != nil {
		return requestResp.Code, errors.New("requestResp data error :" + err.Error()), httpRes
	}

	return 0, nil, httpRes
}
