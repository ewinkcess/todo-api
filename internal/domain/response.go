package domain

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func NewSuccessResponse(message string, data interface{}) Response {
	return Response{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

func NewErrorResponse(message string, errors interface{}) Response {
	return Response{
		Status:  "error",
		Message: message,
		Errors:  errors,
	}
}
