package server

type Response struct {
	Reason string `json:"reason"`
}

func Error(reason string) Response {
	return Response{
		Reason: reason,
	}
}
