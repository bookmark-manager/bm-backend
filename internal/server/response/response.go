package response

type Response struct {
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func Error(msg string) *Response {
	return &Response{
		Error: msg,
	}
}
