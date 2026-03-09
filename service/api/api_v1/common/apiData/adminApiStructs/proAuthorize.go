package adminApiStructs

type ProAuthorizeUpdateUserExpiredTimeByDayReq struct {
	UserId    uint   `json:"userId"`
	DayNum    int    `json:"dayNum"`
	Note      string `json:"note"`
	AdminNote string `json:"adminNote"`
	// OrderNo string `json:"orderNo"`
}

type ProAuthorizeGetUserProAuthorizeListItemResp struct {
	UserId      uint   `json:"userId"`
	Username    string `json:"username"`
	ExpiredTime string `json:"expiredTime"`
	Name        string `json:"name"`
}
