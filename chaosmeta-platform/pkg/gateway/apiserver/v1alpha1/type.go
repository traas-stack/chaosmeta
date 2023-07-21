package v1alpha1

type ResponseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TraceID string `json:"trace_id"`
}

func NewResponseData(code int, message string, traceID string) ResponseData {
	return ResponseData{
		Code:    code,
		Message: message,
		TraceID: traceID,
	}
}
