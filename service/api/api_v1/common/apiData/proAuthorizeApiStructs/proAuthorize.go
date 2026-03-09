package proAuthorizeApiStructs

import "time"

type GetAuthorizeResp struct {
	ExpiredTime string `json:"expiredTime"`
}

type GetAuthorizeHistoryRecordResp struct {
	ChangeTime  time.Time `json:"changeTime"`
	ExpiredTime time.Time `json:"expiredTime"`
	DayNum      int       `json:"dayNum"`
	Note        string    `json:"note"`
	OrderNo     string    `json:"orderNo"`
	AdminNote   string    `json:"adminNote"`
}
