package common

type ListQueryFormatResponse struct {
	Current  int         `json:"current"`
	PageSize int         `json:"pageSize"`
	Total    int64       `json:"total"`
	List     interface{} `json:"list"`
}

type UniCodeRequest struct {
	UniCode string `json:"uniCode"`
}
