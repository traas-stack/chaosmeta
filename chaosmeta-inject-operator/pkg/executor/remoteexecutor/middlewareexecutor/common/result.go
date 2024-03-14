package common

type TaskResult struct {
	TaskId    string
	Success   bool
	ErrorCode string
	Message   string
	Result    string
	Host      string
}
