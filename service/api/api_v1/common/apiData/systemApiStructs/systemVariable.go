package systemApiStructs

type SystemVariableListItem struct {
	Id          uint   `json:"id"`
	Description string `json:"description"`
	ConfigName  string `json:"configName"`
	ConfigValue string `json:"configValue"`
}

type SystemVariable struct {
	Description string      `json:"description"`
	Name        string      `json:"name"`
	Value       interface{} `json:"value"`
}

type SystemVariableEditReq struct {
	Id          uint   `json:"id"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Value       string `json:"value"`
}
