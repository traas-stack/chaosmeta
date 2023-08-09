package experiment

import (
	"chaosmeta-platform/pkg/service/experiment"
)

type CreateExperimentResponse struct {
	Uuid string `json:"uuid"`
}

type GetExperimentResponse struct {
	Experiment experiment.Experiment `json:"experiments"`
}

type ExperimentListResponse struct {
	Page        int                     `json:"page"`
	PageSize    int                     `json:"pageSize"`
	Total       int64                   `json:"total"`
	Experiments []experiment.Experiment `json:"experiments"`
}
