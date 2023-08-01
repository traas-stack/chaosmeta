package cluster

type CreateClusterRequest struct {
	Name       string `json:"name"`
	Kubeconfig string `json:"kubeconfig"`
}

type CreateClusterResponse struct {
	ID int64 `json:"id"`
}

type ClusterData struct {
	Id         interface{} `json:"id"`
	Name       string      `json:"name"`
	Kubeconfig string      `json:"kubeconfig"`
}

type ListClusterResponse struct {
	Page     int           `json:"page"`
	PageSize int           `json:"pageSize"`
	Total    int64         `json:"total"`
	Clusters []ClusterData `json:"clusters"`
}

type UpdateClusterRequest struct {
	Name       string `json:"name"`
	Kubeconfig string `json:"kubeconfig"`
}
