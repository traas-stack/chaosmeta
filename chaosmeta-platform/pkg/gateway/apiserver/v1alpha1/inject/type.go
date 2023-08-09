package inject

import "chaosmeta-platform/pkg/models/inject/basic"

type ScopesListResponse struct {
	Page     int           `json:"page"`
	PageSize int           `json:"pageSize"`
	Total    int64         `json:"total"`
	Scopes   []basic.Scope `json:"scopes"`
}

type TargetsListResponse struct {
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
	Total    int64          `json:"total"`
	Targets  []basic.Target `json:"targets"`
}

type FaultsListResponse struct {
	Page     int           `json:"page"`
	PageSize int           `json:"pageSize"`
	Total    int64         `json:"total"`
	Faults   []basic.Fault `json:"faults"`
}

type ArgsListResponse struct {
	Page     int          `json:"page"`
	PageSize int          `json:"pageSize"`
	Total    int64        `json:"total"`
	Args     []basic.Args `json:"args"`
}
