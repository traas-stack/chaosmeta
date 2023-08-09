package kube

type QueryNodeResponse struct {
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	//Nodes    []v1.Node `json:"nodes"`
}

type QueryNamespaceResponse struct {
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	//Namespaces []v1.Namespace `json:"namespaces"`
}

type QueryPodResponse struct {
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	//Pods     []v1.Pod `json:"pods"`
}

type QueryDeploymentResponse struct {
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	//Deployments []appv1.Deployment `json:"deployments"`
}
