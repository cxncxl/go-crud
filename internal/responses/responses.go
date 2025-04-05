package responses

import "encoding/json"

type Response struct {
	Status string        `json:"status"`
	Error  string        `json:"error,omitempty"`
	Data   *ResponseData `json:"data,omitempty"`
}

func (self Response) Json() []byte {
    res, err := json.Marshal(self);
    if err != nil {
        panic(err);
    }

    return res;
}

func (self Response) JsonString() string {
    return string(self.Json());
}

type ResponseData struct {
	Kind string `json:"kind"`
	Data any    `json:"data"`
}

func NewErrorResponse(err string) Response {
	return Response{
		Status: "error",
		Error:  err,
	}
}

func NewDataResponse(kind string, data any) Response {
	return Response{
		Status: "ok",
		Data: &ResponseData{
			Kind: kind,
			Data: data,
		},
	}
}
