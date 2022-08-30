package response

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Success(data any) Response {
	return Response{
		Message: "success",
		Data:    data,
	}
}

func Error(err error, data any) Response {
	return Response{
		Code:    101,
		Message: err.Error(),
		Data:    data,
	}
}
