package response

type Response struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

func Error(msg string, status int) Response {
	return Response{
		Status: status,
		Error:  msg,
	}
}
func StatusOK(msg string, status int) Response {
	return Response{
		Status: status,
		Error:  msg,
	}
}
